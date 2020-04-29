package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lokmannicholas/delivery/pkg/managers"
)

type OrderController interface {
	PlaceOrder(c *gin.Context)
	TakeOrder(c *gin.Context)
	OrderList(c *gin.Context)
}
type OrderControllerImp struct {
	OrdersManager managers.OrdersManager
}

func GetOrderController() OrderController {
	return &OrderControllerImp{
		OrdersManager: managers.GetOrdersManager(),
	}
}

func (ctl *OrderControllerImp) PlaceOrder(c *gin.Context) {

	type Request struct {
		Origin      []string `json:"origin"`
		Destination []string `json:"destination"`
	}
	request := new(Request)
	defer c.Request.Body.Close()
	if c.Request.ContentLength == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "request body is missing"})
		return
	}
	if err := c.BindJSON(request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ord, err := ctl.OrdersManager.PlaceOrder(request.Origin, request.Destination)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, ord)
}

func (ctl *OrderControllerImp) TakeOrder(c *gin.Context) {
	type Request struct {
		Status string `json:"status"`
	}
	request := new(Request)
	defer c.Request.Body.Close()
	if c.Request.ContentLength == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "request body is missing"})
		return
	}
	if err := c.BindJSON(request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if request.Status == "TAKEN" {
		s := c.Param("id")
		id, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "incorrect id"})
			return
		}
		ord, err := ctl.OrdersManager.TakeOrder(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if ord == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "order cannot not be taken "})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "SUCCESS"})
		return
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "request body not correct"})
	}
}
func (ctl *OrderControllerImp) OrderList(c *gin.Context) {
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "page should be number"})
		return
	}
	if page < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "page should be larger or equal to 1"})
		return
	}
	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "limit should be number"})
		return
	}
	ords, err := ctl.OrdersManager.GetOrders(uint64(page), uint64(limit))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if len(ords) == 0 {
		c.JSON(http.StatusOK, []interface{}{})
		return
	}

	c.JSON(http.StatusOK, ords)
}
