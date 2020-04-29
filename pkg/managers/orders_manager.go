package managers

import (
	"errors"
	"strconv"

	"github.com/lokmannicholas/delivery/pkg/datacollection"

	"github.com/lokmannicholas/delivery/pkg/repositories"

	"github.com/lokmannicholas/delivery/pkg/models"
	"github.com/lokmannicholas/delivery/pkg/services"
)

type OrdersManager interface {
	PlaceOrder(startCor, EndCor []string) (*models.Orders, error)
	TakeOrder(id int64) (*models.Orders, error)
	GetOrders(page, limit uint64) ([]*models.Orders, error)
}

type OrdersManagerImp struct {
	DistanceFounder services.DistanceFounder
	OrdersRepo      repositories.OrdersRepository
	MySqlHelper     datacollection.MySQLHelper
}

func GetOrdersManager() OrdersManager {
	return &OrdersManagerImp{
		DistanceFounder: services.GetDistanceFounder(),
		OrdersRepo:      repositories.GetOrdersRepository(),
		MySqlHelper:     datacollection.GetMySQLHelper(),
	}
}

func (m *OrdersManagerImp) PlaceOrder(startCor, endCor []string) (*models.Orders, error) {
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
	dis := m.DistanceFounder.CountDistance(startCor, endCor)
	if dis == 0 {
		return nil, errors.New("distance calculation error")
	}
	db := m.MySqlHelper.DB()
	return m.OrdersRepo.Insert(db, dis)
}

func (m *OrdersManagerImp) TakeOrder(id int64) (*models.Orders, error) {
	db := m.MySqlHelper.DB()
	return m.OrdersRepo.UpdateById(db, id)
}

func (m *OrdersManagerImp) GetOrders(page, limit uint64) ([]*models.Orders, error) {
	db := m.MySqlHelper.DB()
	return m.OrdersRepo.Find(db, page, limit)
}
