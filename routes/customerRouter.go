package routes

import (
	"restaurant-management/controllers"

	"github.com/gin-gonic/gin"
)

func PublicCustomerRoutes(incomingRoutes *gin.Engine) {
	// Restaurant routes for customers
	incomingRoutes.GET("/customer/restaurants/:restaurant_id", controllers.CustomerGetRestaurant())

	// Table routes for customers
	incomingRoutes.GET("/customer/tables/:table_id", controllers.CustomerGetTable())

	// Menu routes for customers
	incomingRoutes.GET("/customer/menus", controllers.CustomerGetMenus())

	// Food routes for customers
	incomingRoutes.GET("/customer/foods", controllers.CustomerGetFoodsByRestaurantID())

	// Order routes for customers
	incomingRoutes.POST("/customer/orders", controllers.CustomerCreateOrderId())
	incomingRoutes.PATCH("/customer/orders/:order_id", controllers.CustomerUpdateOrder())

	// Order item routes for customers
	incomingRoutes.POST("/customer/order-items", controllers.CustomerCreateOrderItem())
	incomingRoutes.PATCH("/customer/order-items/:order_item_id", controllers.CustomerUpdateOrderItem())
	incomingRoutes.DELETE("/customer/order-items/:order_item_id", controllers.CustomerDeleteOrderItem())
}
