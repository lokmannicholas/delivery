package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/lokmannicholas/delivery/controller"
	"github.com/lokmannicholas/delivery/pkg/config"
	"github.com/lokmannicholas/delivery/pkg/datacollection"
)

func main() {
	defer datacollection.GetMySQLHelper().Close()
	r := gin.Default()
	r.POST("/orders", controller.GetOrderController().PlaceOrder)
	r.PATCH("/orders/:id", controller.GetOrderController().TakeOrder)
	r.GET("/orders", controller.GetOrderController().OrderList)
	r.Run(fmt.Sprintf(":%s", config.Get().Port)) // default listen and serve on 0.0.0.0:8080

}
