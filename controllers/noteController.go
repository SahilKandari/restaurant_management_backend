package controllers

import (
	"context"
	"net/http"
	"restaurant-management/models"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func GetNotes() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		id := c.Param("restaurant_id")
		if id == "" {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Restaurant ID is required"})
			return
		}
		var notes []models.Note
		rows, err := Db.QueryContext(ctx, "SELECT * FROM notes WHERE restaurant_id = $1", id)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch notes from database", "details": err.Error()})
			return
		}
		defer rows.Close()
		for rows.Next() {
			var note models.Note
			if err := rows.Scan(&note.ID, &note.Title, &note.Content, &note.Priority, &note.RestaurantID, &note.CreatedAt, &note.UpdatedAt); err != nil {
				c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan notes", "details": err.Error()})
				return
			}
			notes = append(notes, note)
		}
		c.IndentedJSON(http.StatusOK, gin.H{"message": "Notes fetched successfully", "notes": notes})
	}
}

func GetNote() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func CreateNote() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		var note models.Note

		if err := c.BindJSON(&note); err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Provide correct data for creating Note"})
			return
		}

		// Validate the note struct
		if err := validate.Struct(note); err != nil {
			var validationErrors []string
			for _, err := range err.(validator.ValidationErrors) {
				validationErrors = append(validationErrors, err.Field()+" failed on the '"+err.Tag()+"' tag")
			}
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Validation failed", "details": validationErrors})
			return
		}

		query := `
			INSERT INTO notes (title, content, priority, restaurant_id)
			VALUES ($1, $2, $3, $4)
			RETURNING id, title, content, priority, restaurant_id, created_at, updated_at
		`

		if err := Db.QueryRowContext(ctx, query, note.Title, note.Content, note.Priority, note.RestaurantID).Scan(&note.ID, &note.Title, &note.Content, &note.Priority, &note.RestaurantID, &note.CreatedAt, &note.UpdatedAt); err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Failed to create Note in Database", "details": err.Error()})
			return
		}

		c.IndentedJSON(http.StatusOK, gin.H{"message": "Note created successfully", "note": note})
	}
}

func UpdateNote() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		id := c.Param("note_id")
		if id == "" {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Note ID is required"})
			return
		}

		var note models.Note
		if err := c.BindJSON(&note); err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Provide correct data for updating Note"})
			return
		}

		// Validate the note struct
		if err := validate.Struct(note); err != nil {
			var validationErrors []string
			for _, err := range err.(validator.ValidationErrors) {
				validationErrors = append(validationErrors, err.Field()+" failed on the '"+err.Tag()+"' tag")
			}
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Validation failed", "details": validationErrors})
			return
		}

		query := `
			UPDATE notes
			SET title = $1, content = $2, priority = $3, updated_at = CURRENT_TIMESTAMP
			WHERE id = $4
			RETURNING id, title, content, priority, restaurant_id, created_at, updated_at
		`

		if err := Db.QueryRowContext(ctx, query, note.Title, note.Content, note.Priority, id).Scan(&note.ID, &note.Title, &note.Content, &note.Priority, &note.RestaurantID, &note.CreatedAt, &note.UpdatedAt); err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Failed to update Note in Database", "details": err.Error()})
			return
		}

		c.IndentedJSON(http.StatusOK, gin.H{"message": "Note updated successfully", "note": note})
	}
}

func DeleteNote() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		id := c.Param("note_id")
		if id == "" {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Note ID is required"})
			return
		}

		query := "DELETE FROM notes WHERE id = $1"
		if _, err := Db.ExecContext(ctx, query, id); err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete Note from Database", "details": err.Error()})
			return
		}

		c.IndentedJSON(http.StatusOK, gin.H{"message": "Note deleted successfully"})
	}
}
