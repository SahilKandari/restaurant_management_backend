package routes

import (
	"restaurant-management/controllers"

	"github.com/gin-gonic/gin"
)

func RestaurantRoutes(incomingRoutes gin.IRoutes) {
	incomingRoutes.GET("/restaurant", controllers.GetRestaurants())
	incomingRoutes.GET("/restaurant/:restaurant_id", controllers.GetRestaurant())
	incomingRoutes.GET("/restaurant-owner/:owner_id", controllers.GetRestaurantsByOwner())
	incomingRoutes.POST("/restaurant", controllers.CreateRestaurant())
	incomingRoutes.PATCH("/restaurant/:restaurant_id", controllers.UpdateRestaurant())
	incomingRoutes.DELETE("/restaurant/:restaurant_id", controllers.DeleteRestaurant())
}
