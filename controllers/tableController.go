package controllers

import (
	"context"
	"net/http"
	"restaurant-management/models"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func GetTables() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		rows, err := Db.QueryContext(ctx, "SELECT * FROM tables")
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch tables from database", "details": err.Error()})
			return
		}
		defer rows.Close()

		var tables []models.Table

		for rows.Next() {
			var table models.Table
			if err := rows.Scan(&table.ID, &table.Name, &table.Capacity, &table.RestaurantID, &table.Location, &table.Status, &table.CreatedAt, &table.UpdatedAt); err != nil {
				c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan table data", "details": err.Error()})
				return
			}
			tables = append(tables, table)
		}

		c.IndentedJSON(http.StatusOK, gin.H{"message": "Tables fetched successfully", "tables": tables})
	}
}

func GetTable() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		id := c.Param("table_id")
		if id == "" {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Table ID is required"})
			return
		}

		var table models.Table
		if err := Db.QueryRowContext(ctx, `
			SELECT id, name, capacity, restaurant_id, location, status, created_at, updated_at FROM tables WHERE id = $1
		`, id).Scan(&table.ID, &table.Name, &table.Capacity, &table.RestaurantID, &table.Location, &table.Status, &table.CreatedAt, &table.UpdatedAt); err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch table from database", "details": err.Error()})
			return
		}

		c.IndentedJSON(http.StatusOK, gin.H{"message": "Table fetched successfully", "table": table})
	}
}

func GetTablesByRestaurantId() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		id := c.Query("restaurant_id")
		if id == "" {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Restaurant ID is required."})
			return
		}

		rows, err := Db.QueryContext(ctx, "SELECT * FROM tables WHERE restaurant_id = $1", id)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch tables from Database", "details": err.Error()})
			return
		}
		defer rows.Close()

		var tables []models.Table

		for rows.Next() {
			var table models.Table
			if err := rows.Scan(&table.ID, &table.Name, &table.Capacity, &table.RestaurantID, &table.Location, &table.Status, &table.CreatedAt, &table.UpdatedAt); err != nil {
				c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Failed to Scan the table from Database", "details": err.Error()})
				return
			}
			tables = append(tables, table)
		}

		c.IndentedJSON(http.StatusOK, gin.H{"message": "Tables fetched successfully", "tables": tables})
	}
}

func CreateTable() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		var table models.Table
		err := c.BindJSON(&table)
		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Please provide correct data to create table", "details": err.Error()})
			return
		}

		// Validate the table struct
		if err := validate.Struct(table); err != nil {
			var validationErrors []string
			for _, err := range err.(validator.ValidationErrors) {
				validationErrors = append(validationErrors, err.Field()+" failed on the '"+err.Tag()+"' tag")
			}
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Validation failed", "details": validationErrors})
			return
		}

		if err := Db.QueryRowContext(ctx, `
			INSERT INTO tables (name, capacity, restaurant_id, location, status)
			VALUES ($1, $2, $3, $4, $5)
			RETURNING id, name, capacity, restaurant_id, location, status, created_at, updated_at
		`, table.Name, table.Capacity, table.RestaurantID, table.Location, table.Status).Scan(
			&table.ID, &table.Name, &table.Capacity, &table.RestaurantID, &table.Location, &table.Status, &table.CreatedAt, &table.UpdatedAt); err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Failed to create table in database", "details": err.Error()})
			return
		}

		c.IndentedJSON(http.StatusCreated, gin.H{"message": "Table created successfully", "table": table})
	}
}

func UpdateTable() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		id := c.Param("table_id")
		if id == "" {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Table ID is required"})
			return
		}

		var table models.Table

		if err := c.BindJSON(&table); err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Provide correct data for the table", "details": err.Error()})
			return
		}

		// Validate the table struct
		if err := validate.Struct(table); err != nil {
			var validationErrors []string
			for _, err := range err.(validator.ValidationErrors) {
				validationErrors = append(validationErrors, err.Field()+" failed on the '"+err.Tag()+"' tag")
			}
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Validation failed", "details": validationErrors})
			return
		}

		if err := Db.QueryRowContext(ctx, `
			UPDATE tables
			SET name = $1, capacity = $2, restaurant_id = $3, location = $4, status = $5, updated_at = CURRENT_TIMESTAMP
			WHERE id = $6 RETURNING id, name, capacity, restaurant_id, location, status, created_at, updated_at
		`, table.Name, table.Capacity, table.RestaurantID, table.Location, table.Status, id).Scan(&table.ID, &table.Name, &table.Capacity, &table.RestaurantID, &table.Location, &table.Status, &table.CreatedAt, &table.UpdatedAt); err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Failed to update the table in database", "details": err.Error()})
			return
		}

		c.IndentedJSON(http.StatusOK, gin.H{"message": "Table updated successfully", "table": table})
	}
}

func DeleteTable() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		id := c.Param("table_id")
		if id == "" {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Table ID is required"})
			return
		}

		result, err := Db.ExecContext(ctx, "DELETE FROM tables WHERE id = $1", id)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete table from database", "details": err.Error()})
			return
		}
		rowsAffected, err := result.RowsAffected()
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve affected rows", "details": err.Error()})
			return
		}
		if rowsAffected == 0 {
			c.IndentedJSON(http.StatusNotFound, gin.H{"error": "No table found with given ID"})
			return
		}
		c.IndentedJSON(http.StatusOK, gin.H{"message": "Table deleted successfully", "table_id": id})
	}
}
