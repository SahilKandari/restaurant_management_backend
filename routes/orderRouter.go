package routes

import (
	"restaurant-management/controllers"

	"github.com/gin-gonic/gin"
)

func OrderRoutes(incomingRoutes gin.IRoutes) {
	incomingRoutes.GET("/orders", controllers.GetOrders())
	incomingRoutes.GET("/orders/:order_id", controllers.GetOrder())
	incomingRoutes.GET("/order-id", controllers.CreateOrderId())
	incomingRoutes.POST("/orders", controllers.CreateOrder())
	incomingRoutes.PATCH("/orders/:order_id", controllers.UpdateOrder())
	incomingRoutes.PATCH("/orders-status/:order_id", controllers.UpdateOrderStatus())
	incomingRoutes.DELETE("/orders/:order_id", controllers.DeleteOrder())
}
