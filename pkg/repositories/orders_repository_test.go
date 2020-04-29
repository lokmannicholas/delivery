package repositories

import (
	"errors"
	"math/rand"
	"testing"

	"github.com/lokmannicholas/delivery/pkg/models"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type OrdersRepositoryTestSuite struct {
	suite.Suite
	//DBMock sqlmock.Sqlmock
	//TestDB *sql.DB
	Fields []string
}

func (suite *OrdersRepositoryTestSuite) SetupSuite() {
	// one m DB connection for all tests

	suite.Fields = []string{
		"id",
		"distance",
		"status",
	}
	//db = db
	//mockDB = mock
}
func TestOrdersRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(OrdersRepositoryTestSuite))
}

func (suite *OrdersRepositoryTestSuite) TestGetOrdersRepository() {
	t := suite.T()
	t.Run("GetOrdersRepository", func(t *testing.T) {
		mgr := &OrdersRepositoryImpl{}
		assert.NotNil(t, mgr, "no error")
	})
}
func (suite *OrdersRepositoryTestSuite) TestOrdersRepositoryImpl_Insert() {
	t := suite.T()
	distance := rand.Intn(100000)
	insertID := rand.Int63n(100-1) + 1
	expectedOrd := &models.Orders{
		ID:       insertID,
		Distance: distance,
		Status:   "UNASSIGNED",
	}
	t.Run("Should insert order with distance", func(t *testing.T) {
		r := &OrdersRepositoryImpl{}
		result := sqlmock.NewResult(insertID, 1)
		db, mockDB, _ := sqlmock.New()
		mockDB.ExpectBegin()
		mockDB.ExpectExec("INSERT INTO orders \\(distance,status\\) VALUES \\(\\?,\"UNASSIGNED\"\\)").
			WithArgs(
				distance,
			).WillReturnResult(result)
		id, err := result.LastInsertId()
		assert.Nil(t, err)
		rows := sqlmock.NewRows(suite.Fields).
			AddRow(expectedOrd.ID, expectedOrd.Distance, expectedOrd.Status)
		mockDB.ExpectQuery("select \\* from orders where id = \\?  AND status = \"UNASSIGNED\" Limit 1 ").
			WithArgs(id).WillReturnRows(rows)
		mockDB.ExpectCommit()
		ord, err := r.Insert(db, distance)
		assert.Nil(t, err)
		assert.EqualValues(t, expectedOrd, ord)
	})
	t.Run("Should return error when insert failed", func(t *testing.T) {

		r := &OrdersRepositoryImpl{}
		result := sqlmock.NewResult(0, 0)
		db, mockDB, _ := sqlmock.New()
		mockDB.ExpectBegin()
		mockDB.ExpectExec("INSERT INTO orders \\(distance,status\\) VALUES \\(\\?,\"UNASSIGNED\"\\)").
			WithArgs(
				distance,
			).WillReturnResult(result)

		id, err := result.LastInsertId()
		assert.EqualValues(t, id, 0)
		rows := sqlmock.NewRows(suite.Fields).
			AddRow(expectedOrd.ID, expectedOrd.Distance, expectedOrd.Status)
		mockDB.ExpectQuery("select \\* from orders where id = \\?  AND status = \"UNASSIGNED\" Limit 1 ").
			WithArgs(id).WillReturnRows(rows)
		mockDB.ExpectCommit()
		ord, err := r.Insert(db, distance)
		assert.Nil(t, ord)
		assert.NotNil(t, err)
	})
}

func (suite *OrdersRepositoryTestSuite) TestOrdersRepositoryImpl_UpdateById() {
	t := suite.T()
	id := rand.Int63n(10000-1) + 1
	expectedOrd := &models.Orders{
		ID:       id,
		Distance: rand.Int(),
		Status:   "ASSIGNED",
	}
	t.Run("Should Update order By Id", func(t *testing.T) {

		r := &OrdersRepositoryImpl{}
		result := sqlmock.NewResult(id, 1)
		db, mockDB, _ := sqlmock.New()
		mockDB.ExpectBegin()
		rows := sqlmock.NewRows(suite.Fields).
			AddRow(expectedOrd.ID, expectedOrd.Distance, expectedOrd.Status)
		mockDB.ExpectQuery("select \\* from orders where id = \\?  AND status = \"UNASSIGNED\" Limit 1  FOR UPDATE").
			WithArgs(id).WillReturnRows(rows)
		mockDB.ExpectExec("UPDATE orders Set status = \"ASSIGNED\" where id = \\?").
			WithArgs(
				id,
			).WillReturnResult(result)
		mockDB.ExpectCommit()
		ord, err := r.UpdateById(db, id)
		assert.Nil(t, err)
		assert.EqualValues(t, expectedOrd, ord)
	})
	t.Run("Should return error when order not found", func(t *testing.T) {

		r := &OrdersRepositoryImpl{}
		result := sqlmock.NewResult(id, 1)
		rows := sqlmock.NewRows(suite.Fields)
		db, mockDB, _ := sqlmock.New()
		mockDB.ExpectBegin()
		mockDB.ExpectQuery("select \\* from orders where id = \\?  AND status = \"UNASSIGNED\" Limit 1  FOR UPDATE").
			WithArgs(id).WillReturnRows(rows)
		mockDB.ExpectExec("UPDATE orders Set status = \"ASSIGNED\" where id = \\?").
			WithArgs(
				id,
			).WillReturnResult(result)
		mockDB.ExpectCommit()

		mErr := errors.New("sql: no rows in result set")
		ord, err := r.UpdateById(db, id)
		assert.Nil(t, ord)
		assert.EqualValues(t, mErr, err)
	})
	t.Run("Should return error when update failed", func(t *testing.T) {

		r := &OrdersRepositoryImpl{}
		result := sqlmock.NewResult(0, 0)
		rows := sqlmock.NewRows(suite.Fields).
			AddRow(expectedOrd.ID, expectedOrd.Distance, expectedOrd.Status)
		db, mockDB, _ := sqlmock.New()
		mockDB.ExpectBegin()
		mockDB.ExpectQuery("select \\* from orders where id = \\?  AND status = \"UNASSIGNED\" Limit 1  FOR UPDATE").
			WithArgs(id).WillReturnRows(rows)
		mockDB.ExpectExec("UPDATE orders Set status = \"ASSIGNED\" where id = \\?").
			WithArgs(
				id,
			).WillReturnResult(result)
		mockDB.ExpectCommit()

		ord, err := r.UpdateById(db, id)
		assert.Nil(t, ord)
		assert.NotNil(t, err)
	})
	t.Run("Should return error when id = 0", func(t *testing.T) {
		id := int64(0)
		r := &OrdersRepositoryImpl{}
		result := sqlmock.NewResult(0, 0)
		rows := sqlmock.NewRows(suite.Fields)
		db, mockDB, _ := sqlmock.New()
		mockDB.ExpectBegin()
		mockDB.ExpectQuery("select \\* from orders where id = \\?  AND status = \"UNASSIGNED\" Limit 1  FOR UPDATE").
			WithArgs(id).WillReturnRows(rows)
		mockDB.ExpectExec("UPDATE orders Set status = \"ASSIGNED\" where id = \\?").
			WithArgs(
				id,
			).WillReturnResult(result)
		mockDB.ExpectCommit()

		mErr := errors.New("incorrect id")
		ord, err := r.UpdateById(db, id)
		assert.Nil(t, ord)
		assert.EqualValues(t, mErr, err)
	})
}

func (suite *OrdersRepositoryTestSuite) TestOrdersRepositoryImpl_Find() {
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
	t.Run("Should Select orders", func(t *testing.T) {

		page := uint64(rand.Int63n(100-1) + 1)
		limit := uint64(rand.Int63n(100-1) + 1)
		r := &OrdersRepositoryImpl{}
		rows := sqlmock.NewRows(suite.Fields).
			AddRow(expectedOrds[0].ID, expectedOrds[0].Distance, expectedOrds[0].Status).
			AddRow(expectedOrds[1].ID, expectedOrds[1].Distance, expectedOrds[1].Status)
		db, mockDB, _ := sqlmock.New()
		mockDB.ExpectQuery("select \\* from orders limit \\?  OFFSET \\?").
			WithArgs(limit, (page-1)*limit).WillReturnRows(rows)
		ords, err := r.Find(db, page, limit)
		assert.Nil(t, err)
		assert.EqualValues(t, expectedOrds, ords)
	})
	t.Run("Should Select orders no results", func(t *testing.T) {

		page := uint64(rand.Int63n(100-1) + 1)
		limit := uint64(rand.Int63n(100-1) + 1)
		r := &OrdersRepositoryImpl{}
		rows := sqlmock.NewRows(suite.Fields)
		db, mockDB, _ := sqlmock.New()
		mockDB.ExpectQuery("select \\* from orders limit \\?  OFFSET \\?").
			WithArgs(limit, (page-1)*limit).WillReturnRows(rows)
		ords, err := r.Find(db, page, limit)
		assert.Nil(t, err)
		assert.EqualValues(t, []*models.Orders(nil), ords)
	})
	t.Run("Should return error when page = 0", func(t *testing.T) {
		page := uint64(0)
		limit := uint64(rand.Int63n(100-1) + 1)
		r := &OrdersRepositoryImpl{}
		rows := sqlmock.NewRows(suite.Fields)
		db, mockDB, _ := sqlmock.New()
		mockDB.ExpectQuery("select \\* from orders limit \\?  OFFSET \\?").
			WithArgs(limit, (page-1)*limit).WillReturnRows(rows)
		ords, err := r.Find(db, page, limit)
		mErr := errors.New("page should not less than 1")
		assert.Nil(t, ords)
		assert.EqualValues(t, mErr, err)
	})
	t.Run("Should return error when limit = 0", func(t *testing.T) {
		page := uint64(rand.Int63n(100-1) + 1)
		limit := uint64(0)
		r := &OrdersRepositoryImpl{}
		rows := sqlmock.NewRows(suite.Fields)
		db, mockDB, _ := sqlmock.New()
		mockDB.ExpectQuery("select \\* from orders limit \\?  OFFSET \\?").
			WithArgs(limit, (page-1)*limit).WillReturnRows(rows)
		ords, err := r.Find(db, page, limit)
		mErr := errors.New("limit should not less than or equal 0")
		assert.Nil(t, ords)
		assert.EqualValues(t, mErr, err)
	})
}
