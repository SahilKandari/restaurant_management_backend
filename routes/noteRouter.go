package routes

import (
	"restaurant-management/controllers"

	"github.com/gin-gonic/gin"
)

func NoteRoutes(incomingRoutes gin.IRoutes) {
	incomingRoutes.GET("/notes/:restaurant_id", controllers.GetNotes())
	// incomingRoutes.GET("/notes/:note_id", controllers.GetNote())
	incomingRoutes.POST("/notes", controllers.CreateNote())
	incomingRoutes.PATCH("/notes/:note_id", controllers.UpdateNote())
	incomingRoutes.DELETE("/notes/:note_id", controllers.DeleteNote())
}
