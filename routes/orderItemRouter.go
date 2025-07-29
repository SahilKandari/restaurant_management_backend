package routes

import (
	"restaurant-management/controllers"

	"github.com/gin-gonic/gin"
)

func OrderItemRoutes(incomingRoutes gin.IRoutes) {
	incomingRoutes.GET("/order-items", controllers.GetOrderItems())
	incomingRoutes.GET("/order-items/:order_item_id", controllers.GetOrderItem())
	incomingRoutes.GET("/order-items-order/:order_id", controllers.GetOrderItemsByOrder())
	incomingRoutes.POST("/order-items", controllers.CreateOrderItem())
	incomingRoutes.PATCH("/order-items/:order_item_id", controllers.UpdateOrderItem())
	incomingRoutes.DELETE("/order-items/:order_item_id", controllers.DeleteOrderItem())
}
