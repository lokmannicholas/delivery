package managers

import (
	"context"
	"database/sql"
	"errors"
	"math/rand"

	"github.com/lokmannicholas/delivery/pkg"

	"github.com/lokmannicholas/delivery/pkg/models"

	"github.com/stretchr/testify/assert"

	"testing"

	dbMock "github.com/lokmannicholas/delivery/pkg/datacollection/mocks"
	repoMock "github.com/lokmannicholas/delivery/pkg/repositories/mocks"
	servicesMock "github.com/lokmannicholas/delivery/pkg/services/mocks"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/suite"
)

type OrdersManagerTestSuite struct {
	suite.Suite
	DBMock  sqlmock.Sqlmock
	TestDB  *sql.DB
	Context context.Context
	Fields  []string
}

func (suite *OrdersManagerTestSuite) SetupSuite() {
	// one m DB connection for all tests
	db, mock, _ := sqlmock.New()
	suite.TestDB = db
	suite.DBMock = mock
	suite.Context = context.Background()
}

func TestOrdersManagerTestSuite(t *testing.T) {
	suite.Run(t, new(OrdersManagerTestSuite))
}

func (suite *OrdersManagerTestSuite) TestGetOrdersManager() {
	t := suite.T()
	t.Run("OrdersManagerImp calls DistanceFounder", func(t *testing.T) {
		mgr := &OrdersManagerImp{
			DistanceFounder: &servicesMock.DistanceFounder{},
		}
		assert.NotNil(t, mgr, "no error")
	})
}

func (suite *OrdersManagerTestSuite) TestOrdersManagerImp_PlaceOrder() {
	t := suite.T()
	start := []string{pkg.FloatToString(-90 + rand.Float64()*(90-(-90))), pkg.FloatToString(-180 + rand.Float64()*(180-(-180)))}
	end := []string{pkg.FloatToString(-90 + rand.Float64()*(90-(-90))), pkg.FloatToString(-180 + rand.Float64()*(180-(-180)))}
	distance := rand.Int()
	expectedOrd := &models.Orders{
		ID:       rand.Int63n(10000-1) + 1,
		Distance: distance,
		Status:   "UNASSIGNED",
	}
	t.Run("OrdersManagerImp calls PlaceOrder", func(t *testing.T) {
		var distFounder = &servicesMock.DistanceFounder{}
		var ordersRepo = &repoMock.OrdersRepository{}
		var mysqlMock = &dbMock.MySQLHelper{}
		mgr := &OrdersManagerImp{
			DistanceFounder: distFounder,
			OrdersRepo:      ordersRepo,
			MySqlHelper:     mysqlMock,
		}
		mysqlMock.On("DB").Return(suite.TestDB)
		distFounder.On("CountDistance", start, end).Return(distance, nil)
		ordersRepo.On("Insert", mysqlMock.DB(), distance).Return(expectedOrd, nil)
		ord, err := mgr.PlaceOrder(start, end)
		assert.Nil(t, err, "no errors")
		assert.EqualValues(t, expectedOrd, ord)
		mysqlMock.AssertCalled(t, "DB")
		distFounder.AssertCalled(t, "CountDistance", start, end)
		ordersRepo.AssertCalled(t, "Insert", suite.TestDB, distance)
	})
	t.Run("Should return error when get db failed", func(t *testing.T) {
		var distFounder = &servicesMock.DistanceFounder{}
		var ordersRepo = &repoMock.OrdersRepository{}
		var mysqlMock = &dbMock.MySQLHelper{}
		mgr := &OrdersManagerImp{
			DistanceFounder: distFounder,
			OrdersRepo:      ordersRepo,
			MySqlHelper:     mysqlMock,
		}
		mErr := errors.New("get db failed")
		mysqlMock.On("DB").Return(nil)
		distFounder.On("CountDistance", start, end).Return(distance, nil)
		ordersRepo.On("Insert", mysqlMock.DB(), distance).Return(nil, mErr)
		ord, err := mgr.PlaceOrder(start, end)
		assert.Nil(t, ord)
		assert.EqualValues(t, mErr, err)
		mysqlMock.AssertCalled(t, "DB")
	})
	t.Run("Should return error when CountDistance fail", func(t *testing.T) {
		var distFounder = &servicesMock.DistanceFounder{}
		var ordersRepo = &repoMock.OrdersRepository{}
		var mysqlMock = &dbMock.MySQLHelper{}
		mgr := &OrdersManagerImp{
			DistanceFounder: distFounder,
			OrdersRepo:      ordersRepo,
			MySqlHelper:     mysqlMock,
		}
		mErr := errors.New("distance calculation error")
		mysqlMock.On("DB").Return(suite.TestDB)
		distFounder.On("CountDistance", start, end).Return(0, mErr)
		ordersRepo.On("Insert", mysqlMock.DB(), distance).Return(nil, mErr)
		ord, err := mgr.PlaceOrder(start, end)
		assert.Nil(t, ord)
		assert.EqualValues(t, mErr, err)
		mysqlMock.AssertCalled(t, "DB")
		distFounder.AssertCalled(t, "CountDistance", start, end)
	})
	t.Run("Should return error when coordinate incorrect", func(t *testing.T) {
		start := []string{pkg.FloatToString(-180 + rand.Float64()*(-90-(-180))), pkg.FloatToString(-360 + rand.Float64()*(-180-(-360)))}
		end := []string{pkg.FloatToString(-180 + rand.Float64()*(-90-(-180))), pkg.FloatToString(-360 + rand.Float64()*(-180-(-360)))}
		var distFounder = &servicesMock.DistanceFounder{}
		var ordersRepo = &repoMock.OrdersRepository{}
		var mysqlMock = &dbMock.MySQLHelper{}
		mgr := &OrdersManagerImp{
			DistanceFounder: distFounder,
			OrdersRepo:      ordersRepo,
			MySqlHelper:     mysqlMock,
		}
		mErr := errors.New("incorrect coordinate")
		mysqlMock.On("DB").Return(suite.TestDB)
		distFounder.On("CountDistance", start, end).Return(0, mErr)
		ordersRepo.On("Insert", mysqlMock.DB(), distance).Return(nil, mErr)
		ord, err := mgr.PlaceOrder(start, end)
		assert.Nil(t, ord)
		assert.EqualValues(t, mErr, err)
		mysqlMock.AssertCalled(t, "DB")
	})
	t.Run("Should return error when PlaceOrder fail", func(t *testing.T) {
		var distFounder = &servicesMock.DistanceFounder{}
		var ordersRepo = &repoMock.OrdersRepository{}
		var mysqlMock = &dbMock.MySQLHelper{}
		mgr := &OrdersManagerImp{
			DistanceFounder: distFounder,
			OrdersRepo:      ordersRepo,
			MySqlHelper:     mysqlMock,
		}
		mysqlMock.On("DB").Return(suite.TestDB)
		distFounder.On("CountDistance", start, end).Return(distance, nil)
		mErr := errors.New("get orders error")
		ordersRepo.On("Insert", mysqlMock.DB(), distance).Return(nil, mErr)
		ord, err := mgr.PlaceOrder(start, end)
		assert.Nil(t, ord)
		assert.EqualValues(t, mErr, err)
		mysqlMock.AssertCalled(t, "DB")
		distFounder.AssertCalled(t, "CountDistance", start, end)
		ordersRepo.AssertCalled(t, "Insert", suite.TestDB, distance)
	})
}

func (suite *OrdersManagerTestSuite) TestOrdersManagerImp_TakeOrder() {
	t := suite.T()
	id := rand.Int63n(10000-1) + 1
	expectedOrd := &models.Orders{
		ID:       id,
		Distance: 100,
		Status:   "ASSIGNED",
	}
	t.Run("OrdersManagerImp calls PlaceOrder", func(t *testing.T) {
		var distFounder = &servicesMock.DistanceFounder{}
		var ordersRepo = &repoMock.OrdersRepository{}
		var mysqlMock = &dbMock.MySQLHelper{}
		mgr := &OrdersManagerImp{
			DistanceFounder: distFounder,
			OrdersRepo:      ordersRepo,
			MySqlHelper:     mysqlMock,
		}
		mysqlMock.On("DB").Return(suite.TestDB)
		ordersRepo.On("UpdateById", mysqlMock.DB(), id).Return(expectedOrd, nil)
		ord, err := mgr.TakeOrder(id)
		assert.Nil(t, err, "no errors")
		assert.EqualValues(t, expectedOrd, ord)
		mysqlMock.AssertCalled(t, "DB")
		ordersRepo.AssertCalled(t, "UpdateById", suite.TestDB, id)
	})

	t.Run("Should return error when get db failed", func(t *testing.T) {
		var distFounder = &servicesMock.DistanceFounder{}
		var ordersRepo = &repoMock.OrdersRepository{}
		var mysqlMock = &dbMock.MySQLHelper{}

		mgr := &OrdersManagerImp{
			DistanceFounder: distFounder,
			OrdersRepo:      ordersRepo,
			MySqlHelper:     mysqlMock,
		}
		mErr := errors.New("get db failed")
		mysqlMock.On("DB").Return(nil)
		ordersRepo.On("UpdateById", mysqlMock.DB(), id).Return(nil, mErr)
		ord, err := mgr.TakeOrder(id)
		assert.Nil(t, ord)
		assert.EqualValues(t, mErr, err)
		mysqlMock.AssertCalled(t, "DB")
	})
	t.Run("Should return error when PlaceOrder fail", func(t *testing.T) {
		var distFounder = &servicesMock.DistanceFounder{}
		var ordersRepo = &repoMock.OrdersRepository{}
		var mysqlMock = &dbMock.MySQLHelper{}

		mgr := &OrdersManagerImp{
			DistanceFounder: distFounder,
			OrdersRepo:      ordersRepo,
			MySqlHelper:     mysqlMock,
		}
		mErr := errors.New("get orders error")
		mysqlMock.On("DB").Return(suite.TestDB)
		ordersRepo.On("UpdateById", mysqlMock.DB(), id).Return(nil, mErr)
		ord, err := mgr.TakeOrder(id)
		assert.Nil(t, ord)
		assert.EqualValues(t, mErr, err)
		mysqlMock.AssertCalled(t, "DB")
		ordersRepo.AssertCalled(t, "UpdateById", suite.TestDB, id)

	})
}

func (suite *OrdersManagerTestSuite) TestOrdersManagerImp_GetOrders() {
	t := suite.T()
	expectedOrds := []*models.Orders{
		{
			ID:       rand.Int63n(10000-1) + 1,
			Distance: rand.Int(),
			Status:   "ASSIGNED",
		},
		{
			ID:       rand.Int63n(10000-1) + 1,
			Distance: rand.Int(),
			Status:   "UNASSIGNED",
		},
	}

	page := uint64(rand.Int63n(100-1) + 1)
	limit := uint64(rand.Int63n(100-1) + 1)
	t.Run("OrdersManagerImp calls GetOrders", func(t *testing.T) {
		var distFounder = &servicesMock.DistanceFounder{}
		var ordersRepo = &repoMock.OrdersRepository{}
		var mysqlMock = &dbMock.MySQLHelper{}

		mgr := &OrdersManagerImp{
			DistanceFounder: distFounder,
			OrdersRepo:      ordersRepo,
			MySqlHelper:     mysqlMock,
		}
		mysqlMock.On("DB").Return(suite.TestDB)
		ordersRepo.On("Find", mysqlMock.DB(), page, limit).Return(expectedOrds, nil)
		ords, err := mgr.GetOrders(page, limit)
		assert.Nil(t, err, "no errors")
		assert.EqualValues(t, expectedOrds, ords)
		mysqlMock.AssertCalled(t, "DB")
		ordersRepo.AssertCalled(t, "Find", suite.TestDB, page, limit)
	})
	t.Run("OrdersManagerImp calls GetOrders with no results", func(t *testing.T) {
		var distFounder = &servicesMock.DistanceFounder{}
		var ordersRepo = &repoMock.OrdersRepository{}
		var mysqlMock = &dbMock.MySQLHelper{}

		mgr := &OrdersManagerImp{
			DistanceFounder: distFounder,
			OrdersRepo:      ordersRepo,
			MySqlHelper:     mysqlMock,
		}
		mysqlMock.On("DB").Return(suite.TestDB)
		ordersRepo.On("Find", mysqlMock.DB(), page, limit).Return([]*models.Orders{}, nil)
		ords, err := mgr.GetOrders(page, limit)
		assert.Nil(t, err, "no errors")
		assert.EqualValues(t, []*models.Orders{}, ords)
		mysqlMock.AssertCalled(t, "DB")
		ordersRepo.AssertCalled(t, "Find", suite.TestDB, page, limit)
	})

	t.Run("Should return error when get db failed", func(t *testing.T) {
		var distFounder = &servicesMock.DistanceFounder{}
		var ordersRepo = &repoMock.OrdersRepository{}
		var mysqlMock = &dbMock.MySQLHelper{}

		mgr := &OrdersManagerImp{
			DistanceFounder: distFounder,
			OrdersRepo:      ordersRepo,
			MySqlHelper:     mysqlMock,
		}
		mErr := errors.New("get db failed")
		mysqlMock.On("DB").Return(nil)
		ordersRepo.On("Find", mysqlMock.DB(), page, limit).Return(nil, mErr)
		ords, err := mgr.GetOrders(page, limit)
		assert.Nil(t, ords)
		assert.EqualValues(t, mErr, err)
		mysqlMock.AssertCalled(t, "DB")
	})
	t.Run("Should return error when GetOrders fail", func(t *testing.T) {
		var distFounder = &servicesMock.DistanceFounder{}
		var ordersRepo = &repoMock.OrdersRepository{}
		var mysqlMock = &dbMock.MySQLHelper{}

		mgr := &OrdersManagerImp{
			DistanceFounder: distFounder,
			OrdersRepo:      ordersRepo,
			MySqlHelper:     mysqlMock,
		}
		mErr := errors.New("get orders error")
		mysqlMock.On("DB").Return(suite.TestDB)
		ordersRepo.On("Find", mysqlMock.DB(), page, limit).Return(nil, mErr)
		ords, err := mgr.GetOrders(page, limit)
		assert.Nil(t, ords)
		assert.EqualValues(t, mErr, err)
		mysqlMock.AssertCalled(t, "DB")
		ordersRepo.AssertCalled(t, "Find", suite.TestDB, page, limit)
	})
	t.Run("Should return error when page = 0 ", func(t *testing.T) {

		page := uint64(0)
		limit := uint64(rand.Int63n(100-1) + 1)
		var distFounder = &servicesMock.DistanceFounder{}
		var ordersRepo = &repoMock.OrdersRepository{}
		var mysqlMock = &dbMock.MySQLHelper{}

		mgr := &OrdersManagerImp{
			DistanceFounder: distFounder,
			OrdersRepo:      ordersRepo,
			MySqlHelper:     mysqlMock,
		}
		mErr := errors.New("page should not less than 1")
		mysqlMock.On("DB").Return(suite.TestDB)
		ordersRepo.On("Find", mysqlMock.DB(), page, limit).Return(nil, mErr)
		ords, err := mgr.GetOrders(page, limit)
		assert.Nil(t, ords)
		assert.EqualValues(t, mErr, err)
		mysqlMock.AssertCalled(t, "DB")
		ordersRepo.AssertCalled(t, "Find", suite.TestDB, page, limit)
	})
	t.Run("Should return error when limit = 0 ", func(t *testing.T) {

		page := uint64(rand.Int63n(100-1) + 1)
		limit := uint64(0)
		var distFounder = &servicesMock.DistanceFounder{}
		var ordersRepo = &repoMock.OrdersRepository{}
		var mysqlMock = &dbMock.MySQLHelper{}

		mgr := &OrdersManagerImp{
			DistanceFounder: distFounder,
			OrdersRepo:      ordersRepo,
			MySqlHelper:     mysqlMock,
		}
		mErr := errors.New("limit should not less than or equal 0")
		mysqlMock.On("DB").Return(suite.TestDB)
		ordersRepo.On("Find", mysqlMock.DB(), page, limit).Return(nil, mErr)
		ords, err := mgr.GetOrders(page, limit)
		assert.Nil(t, ords)
		assert.EqualValues(t, mErr, err)
		mysqlMock.AssertCalled(t, "DB")
		ordersRepo.AssertCalled(t, "Find", suite.TestDB, page, limit)
	})
}
