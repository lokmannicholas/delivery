package managers

import (
	"context"
	"database/sql"

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
	start := []string{"22.4606381", "113.9690943"}
	end := []string{"22.364754", "114.1201489"}
	distance := 100
	expectedOrd := &models.Orders{
		ID:       1,
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
		//suite.DBMock.ExpectBegin()
		//suite.DBMock.ExpectExec(`INSERT INTO orders (distance,status) VALUES (?,"UNASSIGNED")`).
		//	WithArgs(distance).
		//	WillReturnResult(sqlmock.NewResult(1, 1))
		//suite.DBMock.ExpectCommit()
		ord, err := mgr.PlaceOrder(suite.Context, start, end)
		assert.Nil(t, err, "no errors")
		assert.EqualValues(t, expectedOrd, ord)
	})
}

//func (suite *OrdersManagerTestSuite) TestOrdersManagerImp_TakeOrder() {
//
//}
//func (suite *OrdersManagerTestSuite) TestOrdersManagerImp_GetOrders() {
//
//}
func TestOrdersManagerTestSuite(t *testing.T) {
	suite.Run(t, new(OrdersManagerTestSuite))
}
