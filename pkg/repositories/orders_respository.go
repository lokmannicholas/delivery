package repositories

import (
	"database/sql"
	"errors"
	"fmt"
	"runtime/debug"

	"github.com/lokmannicholas/delivery/pkg/models"
)

type OrdersRepository interface {
	Insert(db *sql.DB, distance int) (*models.Orders, error)
	UpdateById(db *sql.DB, id int64) (*models.Orders, error)
	Find(db *sql.DB, page, limit uint64) ([]*models.Orders, error)
}

type OrdersRepositoryImpl struct{}

func GetOrdersRepository() OrdersRepository {
	return &OrdersRepositoryImpl{}
}

func (repo *OrdersRepositoryImpl) Insert(db *sql.DB, distance int) (*models.Orders, error) {
	if db == nil {
		return nil, errors.New("get db failed")
	}
	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}
	defer func() {
		if r := recover(); r != nil {

			debug.PrintStack()
			if err := tx.Rollback(); err != nil {
				_ = fmt.Errorf("%+v", err)
			}
			switch v := r.(type) {
			case string:
				err = errors.New(v)
			case error:
				err = v
			default:
				err = errors.New("unknown panic")
			}
		}
	}()
	result, err := tx.Exec(` INSERT INTO orders (distance,status) VALUES (?,"UNASSIGNED")`, distance)
	if err != nil {
		return nil, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		if err := tx.Rollback(); err != nil {
			return nil, err
		}
		return nil, err
	}
	if id == 0 {
		if err := tx.Rollback(); err != nil {
			return nil, err
		}
		return nil, errors.New("insert failed")
	}
	row := tx.QueryRow(`select * from orders where id = ? AND status = "UNASSIGNED" Limit 1`, id)
	var oid int64
	var dis int
	var status string
	if err := row.Scan(&oid, &dis, &status); err != nil {
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return &models.Orders{
		ID:       id,
		Distance: distance,
		Status:   status,
	}, nil
}

func (repo *OrdersRepositoryImpl) UpdateById(db *sql.DB, id int64) (*models.Orders, error) {
	if db == nil {
		return nil, errors.New("get db failed")
	}
	if id == 0 {
		return nil, errors.New("incorrect id")
	}
	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}
	defer func() {
		if r := recover(); r != nil {

			debug.PrintStack()
			if err := tx.Rollback(); err != nil {
				_ = fmt.Errorf("%+v", err)
			}
			switch v := r.(type) {
			case string:
				err = errors.New(v)
			case error:
				err = v
			default:
				err = errors.New("unknown panic")
			}
		}
	}()
	row := tx.QueryRow(`select * from orders where id = ? AND status = "UNASSIGNED" Limit 1 FOR UPDATE`, id)
	var oid int64
	var dis int
	var status string
	if err := row.Scan(&oid, &dis, &status); err != nil {
		return nil, err
	}
	if result, err := tx.Exec(`UPDATE  orders Set status = "ASSIGNED" where id = ?`, oid); err != nil {
		return nil, err
	} else {
		if raf, err := result.RowsAffected(); err != nil {
			if err := tx.Rollback(); err != nil {
				return nil, err
			}
			return nil, err
		} else if raf != 0 {
			if err := tx.Commit(); err != nil {
				return nil, err
			}
			return &models.Orders{
				ID:       oid,
				Distance: dis,
				Status:   "ASSIGNED",
			}, nil
		}
	}
	return nil, errors.New(fmt.Sprintf("cannot update order id:%d", id))
}

func (repo *OrdersRepositoryImpl) Find(db *sql.DB, page, limit uint64) ([]*models.Orders, error) {
	if db == nil {
		return nil, errors.New("get db failed")
	}
	if page < 1 {
		return nil, errors.New("page should not less than 1")
	}
	if limit <= 0 {
		return nil, errors.New("limit should not less than or equal 0")
	}
	var ords []*models.Orders
	var rows *sql.Rows
	var err error
	if limit > 0 {
		rows, err = db.Query(`select * from orders limit ?  OFFSET ?`, limit, (page-1)*limit)
	}
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var id int64
		var dis int
		var status string
		err = rows.Scan(&id, &dis, &status)
		if err != nil {
			return nil, err
		}
		ords = append(ords, &models.Orders{
			ID:       id,
			Distance: dis,
			Status:   status,
		})
	}
	return ords, nil
}
