package managers

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"database/sql"

	"github.com/lokmannicholas/delivery/pkg/datacollection"
	"github.com/lokmannicholas/delivery/pkg/models"
	"github.com/lokmannicholas/delivery/pkg/services"
)

type OrdersManagerImp interface {
	PlaceOrder(ctx context.Context, startCor, EndCor []string) (*models.Orders, error)
	TakeOrder(ctx context.Context, id int64) (*models.Orders, error)
	GetOrders(ctx context.Context, page, limit uint64) ([]*models.Orders, error)
}

type OrdersManager struct {
}

func GetOrdersManager() OrdersManagerImp {
	return &OrdersManager{}
}

func (m *OrdersManager) PlaceOrder(ctx context.Context, startCor, endCor []string) (*models.Orders, error) {
	//validate latitude and longitude value of coordinates
	if len(startCor) < 2 || len(endCor) < 2 {
		return nil, errors.New("incorrect coordinate")
	}
	if f, err := strconv.ParseFloat(startCor[0], 64); err != nil || f > 90 || f < -90 {
		if err != nil {
			return nil, err
		}
		return nil, errors.New("incorrect coordinate")
	}
	if f, err := strconv.ParseFloat(startCor[1], 64); err != nil || f > 180 || f < -180 {
		if err != nil {
			return nil, err
		}
		return nil, errors.New("incorrect coordinate")
	}
	if f, err := strconv.ParseFloat(endCor[0], 64); err != nil || f > 90 || f < -90 {
		if err != nil {
			return nil, err
		}
		return nil, errors.New("incorrect coordinate")
	}
	if f, err := strconv.ParseFloat(endCor[1], 64); err != nil || f > 180 || f < -180 {
		if err != nil {
			return nil, err
		}
		return nil, errors.New("incorrect coordinate")
	}
	//cal distance
	dis := services.CountDistance(startCor, endCor)
	if dis == 0 {
		return nil, errors.New("distance calculation error")
	}
	var ord *models.Orders
	if err := datacollection.GetMySQLHelper().Tx(ctx, func(db *sql.Tx) error {
		result, err := db.Exec(` INSERT INTO orders (distance,status) VALUES (?,"UNASSIGNED")`, dis)
		if err != nil {
			return err
		}
		id, err := result.LastInsertId()
		if err != nil {
			return err
		}
		ord = new(models.Orders)
		ord.ID = id
		ord.Distance = dis
		ord.Status = "UNASSIGNED"
		return nil
	}); err != nil {
		return nil, err
	}
	return ord, nil
}

func (m *OrdersManager) TakeOrder(ctx context.Context, id int64) (*models.Orders, error) {

	var ord *models.Orders
	if err := datacollection.GetMySQLHelper().Tx(ctx, func(db *sql.Tx) error {
		row := db.QueryRow(`select * from orders where id = ? AND status = "UNASSIGNED" Limit 1 FOR UPDATE`, id)
		var oid int64
		var dis int
		var status string
		if err := row.Scan(&oid, &dis, &status); err != nil {
			return err
		}
		if result, err := db.Exec(`UPDATE  orders Set status = "ASSIGNED" where id = ?`, oid); err != nil {
			return err
		} else {
			if raf, err := result.RowsAffected(); err != nil {
				return err
			} else if raf != 0 {
				ord = &models.Orders{
					ID:       oid,
					Distance: dis,
					Status:   "ASSIGNED",
				}
			}
		}
		return nil
	}); err != nil {

		fmt.Println(err)
		return nil, err
	}
	return ord, nil
}

func (m *OrdersManager) GetOrders(ctx context.Context, page, limit uint64) ([]*models.Orders, error) {
	var ords []*models.Orders
	if err := datacollection.GetMySQLHelper().Tx(ctx, func(db *sql.Tx) error {
		var rows *sql.Rows
		var err error
		if limit > 0 {
			rows, err = db.Query(`select * from orders limit ?  OFFSET ?`, limit, (page-1)*limit)
		}
		if err != nil {
			return err
		}
		for rows.Next() {
			var id int64
			var dis int
			var status string
			err = rows.Scan(&id, &dis, &status)
			if err != nil {
				panic(err.Error())
			}
			ords = append(ords, &models.Orders{
				ID:       id,
				Distance: dis,
				Status:   status,
			})
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return ords, nil
}
