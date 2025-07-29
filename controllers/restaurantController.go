package controllers

import (
	"context"
	"database/sql"
	"net/http"
	"restaurant-management/models"
	"time"

	"github.com/gin-gonic/gin"
)

func GetRestaurants() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		var restaurants []models.Restaurant

		rows, err := Db.QueryContext(ctx, "SELECT * FROM restaurants")
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch restaurants from database", "details": err.Error()})
			return
		}
		defer rows.Close()

		for rows.Next() {
			var restaurant models.Restaurant
			if err := rows.Scan(&restaurant.ID, &restaurant.Name, &restaurant.OwnerID, &restaurant.Logo, &restaurant.Address, &restaurant.Description, &restaurant.CreatedAt, &restaurant.UpdatedAt); err != nil {
				c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan restaurants", "details": err.Error()})
				return
			}
			restaurants = append(restaurants, restaurant)
		}

		c.IndentedJSON(http.StatusOK, gin.H{"message": "Restaurants fetched successfully", "restaurants": restaurants})
	}
}

func GetRestaurant() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		id := c.Param("restaurant_id")
		if id == "" {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Restaurant ID is required"})
			return
		}

		var restaurant models.Restaurant
		if err := Db.QueryRowContext(ctx, "SELECT * FROM restaurants WHERE id = $1", id).Scan(&restaurant.ID, &restaurant.Name, &restaurant.OwnerID, &restaurant.Logo, &restaurant.Address, &restaurant.Description, &restaurant.CreatedAt, &restaurant.UpdatedAt); err != nil {
			if err == sql.ErrNoRows {
				c.IndentedJSON(http.StatusNotFound, gin.H{"error": "Restaurant not found"})
			} else {
				c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch restaurant", "details": err.Error()})
			}
			return
		}

		c.IndentedJSON(http.StatusOK, gin.H{"message": "Restaurant fetched successfully", "restaurant": restaurant})
	}
}

func CreateRestaurant() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		var restaurant models.Restaurant
		if err := c.ShouldBindJSON(&restaurant); err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
			return
		}

		if restaurant.Name == "" || restaurant.OwnerID == 0 || restaurant.Address == "" {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Name, OwnerID, and Address are required"})
			return
		}

		err := Db.QueryRowContext(ctx, "INSERT INTO restaurants (name, owner_id, logo, address, description) VALUES ($1, $2, $3, $4, $5) RETURNING id",
			restaurant.Name, restaurant.OwnerID, restaurant.Logo, restaurant.Address, restaurant.Description).Scan(&restaurant.ID)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Failed to create restaurant", "details": err.Error()})
			return
		}

		c.IndentedJSON(http.StatusCreated, gin.H{"message": "Restaurant created successfully", "restaurant": restaurant})
	}
}

func UpdateRestaurant() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		id := c.Param("restaurant_id")
		if id == "" {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Restaurant ID is required"})
			return
		}

		var restaurant models.Restaurant
		if err := c.ShouldBindJSON(&restaurant); err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
			return
		}

		if restaurant.Name == "" || restaurant.OwnerID == 0 || restaurant.Address == "" {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Name, OwnerID, and Address are required"})
			return
		}

		if err := Db.QueryRowContext(ctx, "UPDATE restaurants SET name = $1, owner_id = $2, logo = $3, address = $4, description = $5 WHERE id = $6 RETURNING id, name, owner_id, logo, address, description, created_at, updated_at",
			restaurant.Name, restaurant.OwnerID, restaurant.Logo, restaurant.Address, restaurant.Description, id).Scan(&restaurant.ID, &restaurant.Name, &restaurant.OwnerID, &restaurant.Logo, &restaurant.Address, &restaurant.Description, &restaurant.CreatedAt, &restaurant.UpdatedAt); err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Failed to update restaurant", "details": err.Error()})
			return
		}

		c.IndentedJSON(http.StatusOK, gin.H{"message": "Restaurant updated successfully", "restaurant": restaurant})
	}
}

func DeleteRestaurant() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		id := c.Param("restaurant_id")
		if id == "" {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Restaurant ID is required"})
			return
		}

		result, err := Db.ExecContext(ctx, "DELETE FROM restaurants WHERE id = $1", id)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete restaurant", "details": err.Error()})
			return
		}

		rowsAffected, _ := result.RowsAffected()
		if rowsAffected == 0 {
			c.IndentedJSON(http.StatusNotFound, gin.H{"error": "Restaurant not found"})
			return
		}

		c.IndentedJSON(http.StatusOK, gin.H{"message": "Restaurant deleted successfully"})
	}
}

func GetRestaurantsByOwner() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		ownerID := c.Param("owner_id")
		if ownerID == "" {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Owner ID is required"})
			return
		}

		var restaurants []models.Restaurant

		rows, err := Db.QueryContext(ctx, "SELECT * FROM restaurants WHERE owner_id = $1 ORDER BY id ASC", ownerID)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch restaurants from database", "details": err.Error()})
			return
		}
		defer rows.Close()

		for rows.Next() {
			var restaurant models.Restaurant
			if err := rows.Scan(&restaurant.ID, &restaurant.Name, &restaurant.OwnerID, &restaurant.Logo, &restaurant.Address, &restaurant.Description, &restaurant.CreatedAt, &restaurant.UpdatedAt); err != nil {
				c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan restaurants", "details": err.Error()})
				return
			}
			restaurants = append(restaurants, restaurant)
		}

		c.IndentedJSON(http.StatusOK, gin.H{"message": "Restaurants fetched successfully", "restaurants": restaurants})
	}
}
