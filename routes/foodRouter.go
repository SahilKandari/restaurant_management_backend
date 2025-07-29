package routes

import (
	"restaurant-management/controllers"

	"github.com/gin-gonic/gin"
)

func FoodRoutes(incomingRoutes gin.IRoutes) {
	incomingRoutes.GET("/foods", controllers.GetFoods())
	incomingRoutes.GET("/foods/:food_id", controllers.GetFood())
	incomingRoutes.POST("/foods", controllers.CreateFood())
	incomingRoutes.PATCH("/foods/:food_id", controllers.UpdateFood())
	incomingRoutes.DELETE("/foods/:food_id", controllers.DeleteFood())
}
