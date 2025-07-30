package controllers

import (
	"context"
	"database/sql"
	"net/http"
	"restaurant-management/models"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func GetFoods() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		id, ok := c.GetQuery("restaurant_id")
		var rows *sql.Rows
		var err error

		if ok && id != "" {
			rows, err = Db.QueryContext(ctx, "SELECT * FROM foods WHERE restaurant_id = $1", id)
		} else {
			rows, err = Db.QueryContext(ctx, "SELECT * FROM foods")
		}

		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch foods from database", "details": err.Error()})
			return
		}

		defer rows.Close()
		var foods []models.Food

		for rows.Next() {
			var food models.Food
			if err := rows.Scan(&food.ID, &food.Name, &food.Price, &food.Description, &food.ImageURL, &food.MenuID, &food.RestaurantID, &food.Ingredients, &food.PrepTime, &food.Calories, &food.SpicyLevel, &food.Vegetarian, &food.Available, &food.CreatedAt, &food.UpdatedAt); err != nil {
				c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan food data", "details": err.Error()})
				return
			}
			foods = append(foods, food)
		}
		c.IndentedJSON(http.StatusOK, gin.H{"message": "Foods fetched successfully", "foods": foods})
	}
}

func GetFood() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("food_id")
		if id == "" {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Food ID is required"})
			return
		}
		var food models.Food
		if err := Db.QueryRow("SELECT * FROM foods WHERE id = $1", id).Scan(&food.ID, &food.Name, &food.Price, &food.Description, &food.ImageURL, &food.MenuID, &food.RestaurantID, &food.Ingredients, &food.PrepTime, &food.Calories, &food.SpicyLevel, &food.Vegetarian, &food.Available, &food.CreatedAt, &food.UpdatedAt); err != nil {
			c.IndentedJSON(http.StatusNotFound, gin.H{"error": "Food not found", "details": err.Error()})
			return
		}

		c.IndentedJSON(http.StatusOK, gin.H{"message": "Food fetched successfully", "food": food})
	}
}

func CreateFood() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		var food models.Food

		if err := c.BindJSON(&food); err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Food data provided is not correct", "details": err.Error()})
			return
		}

		// Validate the food struct
		if err := validate.Struct(food); err != nil {
			var validationErrors []string
			for _, err := range err.(validator.ValidationErrors) {
				validationErrors = append(validationErrors, err.Field()+" failed on the '"+err.Tag()+"' tag")
			}
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Validation failed", "details": validationErrors})
			return
		}

		err := Db.QueryRowContext(ctx, `
			INSERT INTO foods 
			(name, price, description, image_url, menu_id, restaurant_id, ingredients, prep_time, calories, spicy_level, vegetarian, available)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
			RETURNING id, name, price, description, image_url, menu_id, restaurant_id, ingredients, prep_time, calories, spicy_level, vegetarian, available, created_at, updated_at`,
			food.Name, toFixed(food.Price, 2), food.Description, food.ImageURL, food.MenuID,
			food.RestaurantID, food.Ingredients, food.PrepTime, food.Calories, food.SpicyLevel,
			food.Vegetarian, food.Available,
		).Scan(
			&food.ID, &food.Name, &food.Price, &food.Description, &food.ImageURL, &food.MenuID,
			&food.RestaurantID, &food.Ingredients, &food.PrepTime, &food.Calories, &food.SpicyLevel,
			&food.Vegetarian, &food.Available, &food.CreatedAt, &food.UpdatedAt,
		)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Failed to create food item in database", "details": err.Error()})
			return
		}

		c.IndentedJSON(http.StatusCreated, gin.H{"message": "Food item created successfully", "food": food})
	}
}

func UpdateFood() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		id := c.Param("food_id")
		if id == "" {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Food ID is required"})
			return
		}

		var food models.Food
		if err := c.BindJSON(&food); err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Food data provided is not correct", "details": err.Error()})
			return
		}

		// Validate the food struct
		if err := validate.Struct(food); err != nil {
			var validationErrors []string
			for _, err := range err.(validator.ValidationErrors) {
				validationErrors = append(validationErrors, err.Field()+" failed on the '"+err.Tag()+"' tag")
			}
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Validation failed", "details": validationErrors})
			return
		}

		// Update the food item in the database
		if err := Db.QueryRowContext(ctx, `
			UPDATE foods SET 
			name = $1, price = $2, description = $3, image_url = $4,
			menu_id = $5, restaurant_id = $6, ingredients = $7, prep_time = $8,
			calories = $9, spicy_level = $10, vegetarian = $11, available = $12,
			updated_at = CURRENT_TIMESTAMP 
			WHERE id = $13 RETURNING id, name, price, description, image_url, menu_id, restaurant_id, ingredients, prep_time, calories, spicy_level, vegetarian, available, created_at, updated_at`,
			food.Name, toFixed(food.Price, 2), food.Description, food.ImageURL,
			food.MenuID, food.RestaurantID, food.Ingredients, food.PrepTime,
			food.Calories, food.SpicyLevel, food.Vegetarian, food.Available, id,
		).Scan(
			&food.ID, &food.Name, &food.Price, &food.Description, &food.ImageURL, &food.MenuID,
			&food.RestaurantID, &food.Ingredients, &food.PrepTime, &food.Calories, &food.SpicyLevel,
			&food.Vegetarian, &food.Available, &food.CreatedAt, &food.UpdatedAt,
		); err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Failed to update food item in database", "details": err.Error()})
			return
		}

		c.IndentedJSON(http.StatusOK, gin.H{"message": "Food item updated successfully", "food": food})
	}
}

func DeleteFood() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		id := c.Param("food_id")
		if id == "" {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Food ID is required"})
			return
		}
		result, err := Db.ExecContext(ctx, "DELETE FROM foods WHERE id = $1", id)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete food item from database", "details": err.Error()})
			return
		}
		rowsAffected, err := result.RowsAffected()
		if err != nil || rowsAffected == 0 {
			c.IndentedJSON(http.StatusNotFound, gin.H{"error": "Food not found in database", "details": err.Error()})
			return
		}
		c.IndentedJSON(http.StatusNoContent, gin.H{"message": "Food item deleted successfully"})
	}
}

func round(num float64) int {
	if num < 0 {
		return int(num - 0.5)
	}
	return int(num + 0.5)
}

func toFixed(num float64, precision int) float64 {
	factor := float64(1)
	for i := 0; i < precision; i++ {
		factor *= 10
	}
	return float64(round(num*factor)) / factor
}
