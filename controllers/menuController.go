package controllers

import (
	"context"
	"net/http"
	"restaurant-management/models"
	"time"

	"github.com/gin-gonic/gin"
)

func GetMenus() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		rows, err := Db.QueryContext(ctx, "SELECT * FROM menus")
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch menus data from database", "details": err.Error()})
			return
		}
		defer rows.Close()

		var menus []models.Menu

		for rows.Next() {
			var menu models.Menu
			if err := rows.Scan(&menu.ID, &menu.Name, &menu.RestaurantID, &menu.CreatedAt, &menu.UpdatedAt); err != nil {
				c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch menu data", "details": err.Error()})
				return
			}
			menus = append(menus, menu)
		}

		c.IndentedJSON(http.StatusOK, gin.H{"message": "Menus fetched successfully", "menus": menus})
	}
}

func GetMenu() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		id := c.Param("menu_id")
		if id == "" {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Menu ID is required"})
			return
		}

		var menu models.Menu

		if err := Db.QueryRowContext(ctx, "SELECT * FROM menus WHERE id = $1", id).Scan(&menu.ID, &menu.Name, &menu.RestaurantID, &menu.CreatedAt, &menu.UpdatedAt); err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch the menu", "details": err.Error()})
			return
		}

		c.IndentedJSON(http.StatusOK, gin.H{"message": "Menu fetched successfully", "menu": menu})
	}
}

func CreateMenu() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		var menu models.Menu
		if err := c.BindJSON(&menu); err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Menu data is required", "details": err.Error()})
			return
		}

		err := Db.QueryRowContext(ctx, "INSERT INTO menus (name, restaurant_id) VALUES ($1, $2) RETURNING id, name, restaurant_id, created_at, updated_at", menu.Name, menu.RestaurantID).Scan(&menu.ID, &menu.Name, &menu.RestaurantID, &menu.CreatedAt, &menu.UpdatedAt)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Failed to create menu in database", "details": err.Error()})
			return
		}

		c.IndentedJSON(http.StatusCreated, gin.H{"message": "Menu created successfully", "menu": menu})
	}
}

func UpdateMenu() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		id := c.Param("menu_id")
		if id == "" {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Menu ID is required"})
			return
		}

		var menu models.Menu

		if err := c.BindJSON(&menu); err != nil {
			c.IndentedJSON(http.StatusNotFound, gin.H{"error": "Please provide correct menu data", "details": err.Error()})
			return
		}

		result, err := Db.ExecContext(ctx, "UPDATE menus SET name = $1, restaurant_id = $2, updated_at = CURRENT_TIMESTAMP WHERE id = $3", menu.Name, menu.RestaurantID, id)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Failed to update menu in database", "details": err.Error()})
			return
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve rows affected", "details": err.Error()})
			return
		}
		if rowsAffected == 0 {
			c.IndentedJSON(http.StatusNotFound, gin.H{"error": "No item with the given ID"})
			return
		}

		c.IndentedJSON(http.StatusNoContent, gin.H{"message": "Menu updated successfully"})
	}
}

func DeleteMenu() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		id := c.Param("menu_id")
		if id == "" {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Menu ID is required"})
			return
		}

		result, err := Db.ExecContext(ctx, "DELETE FROM menus WHERE id = $1", id)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete the menu from database", "details": err.Error()})
			return
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve rows affected", "details": err.Error()})
			return
		}
		if rowsAffected == 0 {
			c.IndentedJSON(http.StatusNotFound, gin.H{"error": "Menu not found with the given ID"})
			return
		}

		c.IndentedJSON(http.StatusNoContent, gin.H{"message": "Menu deleted successfully"})
	}
}
