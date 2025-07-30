package main

import (
	"fmt"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"restaurant-management/config"
	"restaurant-management/controllers"
	"restaurant-management/database"
	"restaurant-management/middlewares"
	"restaurant-management/routes"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port if not set
	}

	// Initialize database connection after loading .env
	database.InitializeDB()
	defer database.CloseDBConnection() // Ensure the database connection is closed when the application exits

	// Set up database tables
	database.SetupTables()

	// Initialize controllers
	controllers.InitControllers()

	// Create a new Gin router
	router := gin.New()
	router.Use(gin.Logger())

	// Get allowed origins from environment variable
	allowedOrigins := os.Getenv("CORS_ALLOWED_ORIGINS")
	fmt.Printf("Allowed Origins: %s\n", allowedOrigins)

	if allowedOrigins == "" {
		allowedOrigins = "http://localhost:3000" // Default for development
	}
	// Set up CORS middleware
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{allowedOrigins},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	authGroup := router.Group("/")
	authGroup.Use(middlewares.Authentication())

	// Initialize Redis client
	config.InitRedis()

	// SMTP connection
	if err := config.SMTPConnect(); err != nil {
		panic("Failed to connect to SMTP server: " + err.Error())
	}

	// Public Routes
	routes.PublicUserRoutes(router)
	routes.PublicCustomerRoutes(router)

	// Private Routes
	routes.ProtectedUserRoutes(authGroup)
	routes.RestaurantRoutes(authGroup)
	routes.FoodRoutes(authGroup)
	routes.MenuRoutes(authGroup)
	routes.TableRoutes(authGroup)
	routes.OrderRoutes(authGroup)
	routes.OrderItemRoutes(authGroup)
	routes.InvoiceRoutes(authGroup)
	routes.NoteRoutes(authGroup)

	router.Run(":" + port) // Start the server on the specified port
}
