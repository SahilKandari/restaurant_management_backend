package routes

import (
	"restaurant-management/controllers"

	"github.com/gin-gonic/gin"
)

func PublicUserRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.GET("/users", controllers.GetUsers())
	incomingRoutes.GET("/users/:user_id", controllers.GetUser())
	incomingRoutes.POST("/users/signup", controllers.SignUp())
	incomingRoutes.POST("/users/login", controllers.Login())
	incomingRoutes.POST("/users/reset-password-otp", controllers.SendPasswordResetEmail())
	incomingRoutes.POST("/users/reset-password", controllers.PasswordReset())
}

func ProtectedUserRoutes(incomingRoutes gin.IRoutes) {
	incomingRoutes.PATCH("/users/:user_id", controllers.UpdateUser())
	incomingRoutes.DELETE("/users/:user_id", controllers.DeleteUser())

}
