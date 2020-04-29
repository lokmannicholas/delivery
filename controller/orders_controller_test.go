package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http/httptest"
	"testing"

	"github.com/lokmannicholas/delivery/pkg/models"

	"github.com/lokmannicholas/delivery/pkg"

	"github.com/gin-gonic/gin"

	"github.com/lokmannicholas/delivery/pkg/managers/mocks"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/suite"
)

type OrdersControllerTestSuite struct {
	suite.Suite
}

func (suite *OrdersControllerTestSuite) SetupSuite() {
	gin.SetMode(gin.TestMode)
}

func TestOrdersControllerTestSuite(t *testing.T) {
	suite.Run(t, new(OrdersControllerTestSuite))
}

func (suite *OrdersControllerTestSuite) TestGetOrderController() {
	t := suite.T()
	t.Run("OrderControllerImp calls GetOrderController", func(t *testing.T) {
		ctl := &OrderControllerImp{
			OrdersManager: &mocks.OrdersManager{},
		}
		assert.NotNil(t, ctl, "no error")
	})
}
func (suite *OrdersControllerTestSuite) TestOrderControllerImp_PlaceOrder() {
	t := suite.T()
	type Request struct {
		Origin      []string `json:"origin"`
		Destination []string `json:"destination"`
	}
	start := []string{pkg.FloatToString(-90 + rand.Float64()*(90-(-90))), pkg.FloatToString(-180 + rand.Float64()*(180-(-180)))}
	end := []string{pkg.FloatToString(-90 + rand.Float64()*(90-(-90))), pkg.FloatToString(-180 + rand.Float64()*(180-(-180)))}

	req := &Request{
		Origin:      start,
		Destination: end,
	}
	distance := rand.Int()
	expectedOrd := &models.Orders{
		ID:       rand.Int63n(10000-1) + 1,
		Distance: distance,
		Status:   "UNASSIGNED",
	}
	t.Run("OrderControllerImp calls PlaceOrder", func(t *testing.T) {
		data, _ := json.Marshal(req)
		//suite.context.Request = &http.Request{
		//	Body:          ioutil.NopCloser(bytes.NewReader(data)),
		//	ContentLength: int64(len(data)),
		//}
		var orderManagerMock = &mocks.OrdersManager{}
		ctl := &OrderControllerImp{
			OrdersManager: orderManagerMock,
		}
		orderManagerMock.On("PlaceOrder", req.Origin, req.Destination).
			Return(expectedOrd, nil)
		r := gin.Default()
		r.POST("/orders", func(c *gin.Context) {
			ctl.PlaceOrder(c)
		})
		req := httptest.NewRequest("POST", "/orders", bytes.NewReader(data))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, w.Code, 200)
		output := &models.Orders{}
		json.Unmarshal(w.Body.Bytes(), output)
		assert.Equal(t, expectedOrd, output)
	})
}

func (suite *OrdersControllerTestSuite) TestOrderControllerImp_TakeOrder() {
	t := suite.T()
	type Request struct {
		Status string `json:"status"`
	}
	req := &Request{
		Status: "TAKEN",
	}
	distance := rand.Int()
	id := rand.Int63n(10000-1) + 1
	expectedOrd := &models.Orders{
		ID:       id,
		Distance: distance,
		Status:   "ASSIGNED",
	}
	t.Run("OrderControllerImp calls TakeOrder", func(t *testing.T) {
		data, _ := json.Marshal(req)
		//suite.context.Request = &http.Request{
		//	Body:          ioutil.NopCloser(bytes.NewReader(data)),
		//	ContentLength: int64(len(data)),
		//}
		var orderManagerMock = &mocks.OrdersManager{}
		ctl := &OrderControllerImp{
			OrdersManager: orderManagerMock,
		}
		orderManagerMock.On("TakeOrder", id).
			Return(expectedOrd, nil)
		r := gin.Default()
		r.PATCH("/orders/:id", func(c *gin.Context) {
			ctl.TakeOrder(c)
		})
		req := httptest.NewRequest("PATCH", fmt.Sprintf("/orders/%d", id), bytes.NewReader(data))

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, w.Code, 200)
		output := &Request{}
		json.Unmarshal(w.Body.Bytes(), output)
		assert.Equal(t, "SUCCESS", output.Status)
	})
}

func (suite *OrdersControllerTestSuite) TestOrderControllerImp_OrderList() {
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
	t.Run("OrderControllerImp calls OrderList", func(t *testing.T) {

		var orderManagerMock = &mocks.OrdersManager{}
		ctl := &OrderControllerImp{
			OrdersManager: orderManagerMock,
		}
		orderManagerMock.On("GetOrders", page, limit).
			Return(expectedOrds, nil)
		r := gin.Default()
		r.GET("/orders", func(c *gin.Context) {
			ctl.OrderList(c)
		})
		req := httptest.NewRequest("GET", fmt.Sprintf("/orders?page=%d&limit=%d", int(page), int(limit)), nil)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, w.Code, 200)
		output := []*models.Orders{}
		json.Unmarshal(w.Body.Bytes(), &output)
		assert.Equal(t, expectedOrds, output)
	})
}
