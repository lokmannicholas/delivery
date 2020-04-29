package datacollection

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"runtime/debug"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/lokmannicholas/delivery/pkg/config"
)

var mysql *MySQLHelperImp

type MySQLHelper interface {
	DB() *sql.DB
	Tx(ctx context.Context, f func(db *sql.Tx) error) (err error)
	Close()
}
type MySQLHelperImp struct {
	db *sql.DB
}

//TODO: remove this function to prevent potential connection deadlock
func GetMySQLHelper() *MySQLHelperImp {
	if mysql == nil {
		mysql = mySQLConnect()
	} else {
		if err := mysql.db.Ping(); err != nil {
			mysql = mySQLConnect()
		}
	}
	return mysql
}

func mySQLConnect() *MySQLHelperImp {
	mysqlConf := config.Get().Mysql
	//user:password@/dbname?charset=utf8&parseTime=True&loc=Local
	uri := fmt.Sprintf("%s:%s@%s(%s)/%s%s",
		mysqlConf.User, mysqlConf.Password, "tcp", mysqlConf.Addr, mysqlConf.DB, "?charset=utf8mb4&parseTime=True&loc=Local")
	db, err := sql.Open("mysql", uri)
	if err != nil {
		panic(err)
	}
	db.SetMaxIdleConns(20)
	db.SetConnMaxLifetime(2 * time.Minute)

	return &MySQLHelperImp{
		db: db,
	}
}

//Tx Begin Transaction
//TODO: remove this function to prevent potential connection deadlock
func (my *MySQLHelperImp) Tx(ctx context.Context, f func(db *sql.Tx) error) (err error) {
	tx, err := my.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: 0,
		ReadOnly:  false,
	})

	//TODO: review this part coz wrapping panic to error causing lost of trace info
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
	err = f(tx)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return err
	} else {
		return tx.Commit()
	}
}
func (my *MySQLHelperImp) DB() *sql.DB {
	return my.db
}
func (my *MySQLHelperImp) Close() {
	_ = my.db.Close()
}
