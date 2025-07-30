package controllers

import (
	"context"
	"net/http"
	"restaurant-management/models"
	"time"

	"github.com/gin-gonic/gin"
)

func GetOrderItems() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		var orderItems []models.OrderItem

		rows, err := Db.QueryContext(ctx, "SELECT * FROM orderitems")
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch order items from database", "details": err.Error()})
			return
		}

		for rows.Next() {
			var orderItem models.OrderItem
			if err := rows.Scan(&orderItem.ID, &orderItem.OrderID, &orderItem.FoodID, &orderItem.Quantity, &orderItem.UnitPrice, &orderItem.SubTotal, &orderItem.CreatedAt, &orderItem.UpdatedAt); err != nil {
				c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan order items", "details": err.Error()})
				return
			}
			orderItems = append(orderItems, orderItem)
		}

		c.IndentedJSON(http.StatusOK, gin.H{"message": "Order items fetched successfully", "order_items": orderItems})
	}
}

func GetOrderItem() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		id := c.Param("order_item_id")
		if id == "" {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Order item ID is required"})
			return
		}

		var orderItem models.OrderItem
		if err := Db.QueryRowContext(ctx, "SELECT * FROM orderitems WHERE id = $1", id).Scan(&orderItem.ID, &orderItem.OrderID, &orderItem.FoodID, &orderItem.Quantity, &orderItem.UnitPrice, &orderItem.SubTotal, &orderItem.CreatedAt, &orderItem.UpdatedAt); err != nil {
			c.IndentedJSON(http.StatusNotFound, gin.H{"error": "Order item not found", "details": err.Error()})
			return
		}

		c.IndentedJSON(http.StatusOK, gin.H{"message": "Order item fetched successfully", "order_item": orderItem})
	}
}

func CreateOrderItem() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		var orderItem models.OrderItem
		var foodPrice float64

		if err := c.BindJSON(&orderItem); err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Provide correct data for creating order item", "details": err.Error()})
			return
		}

		if err := Db.QueryRowContext(ctx, "SELECT name, price FROM foods WHERE id = $1", orderItem.FoodID).Scan(&orderItem.FoodName, &foodPrice); err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch food in database with the given ID", "food_id": orderItem.FoodID, "details": err.Error()})
			return
		}

		if orderItem.Quantity == 0 {
			orderItem.Quantity = 1
		}

		orderItem.UnitPrice = foodPrice
		orderItem.SubTotal = foodPrice * float64(orderItem.Quantity)

		if err := Db.QueryRowContext(ctx, "INSERT INTO orderitems (order_id, food_id, quantity, unit_price, subtotal) VALUES ($1, $2, $3, $4, $5) RETURNING id, order_id, food_id, quantity, unit_price, subtotal, created_at, updated_at", orderItem.OrderID, orderItem.FoodID, orderItem.Quantity, orderItem.UnitPrice, orderItem.SubTotal).Scan(&orderItem.ID, &orderItem.OrderID, &orderItem.FoodID, &orderItem.Quantity, &orderItem.UnitPrice, &orderItem.SubTotal, &orderItem.CreatedAt, &orderItem.UpdatedAt); err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Failed to create order item in database", "details": err.Error()})
			return
		}

		c.IndentedJSON(http.StatusCreated, gin.H{"message": "Order item created successfully", "order_item": orderItem})
	}
}

func UpdateOrderItem() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		id := c.Param("order_item_id")
		if id == "" {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Order item ID is required"})
			return
		}

		var orderItem models.OrderItem
		var updateOrderItem models.UpdateOrderItem

		if err := c.BindJSON(&updateOrderItem); err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Provide the correct Quantity to update order items", "details": err.Error()})
			return
		}

		if updateOrderItem.Quantity == 0 {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Food quantity for Order Items can't be zero"})
			return
		}

		if err := Db.QueryRowContext(ctx, "SELECT * FROM orderitems WHERE id = $1", id).Scan(&orderItem.ID, &orderItem.OrderID, &orderItem.FoodID, &orderItem.Quantity, &orderItem.UnitPrice, &orderItem.SubTotal, &orderItem.CreatedAt, &orderItem.UpdatedAt); err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Failed fetch order item in database with the given ID", "details": err.Error()})
			return
		}

		orderItem.SubTotal = orderItem.UnitPrice * float64(updateOrderItem.Quantity)

		err := Db.QueryRowContext(ctx,
			"UPDATE orderitems SET quantity = $1, subtotal = $2, updated_at = CURRENT_TIMESTAMP WHERE id = $3 RETURNING id, order_id, food_id, quantity, unit_price, subtotal, created_at, updated_at",
			updateOrderItem.Quantity, orderItem.SubTotal, id,
		).Scan(&orderItem.ID, &orderItem.OrderID, &orderItem.FoodID, &orderItem.Quantity, &orderItem.UnitPrice, &orderItem.SubTotal, &orderItem.CreatedAt, &orderItem.UpdatedAt)
		if err != nil {
			c.IndentedJSON(http.StatusNotFound, gin.H{"error": "No order item found with the given ID or failed to update", "details": err.Error(), "order_item_id": id})
			return
		}

		c.IndentedJSON(http.StatusOK, gin.H{"message": "Order item updated successfully", "order_item": orderItem})
	}
}

func DeleteOrderItem() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		id := c.Param("order_item_id")
		if id == "" {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Order item ID is required"})
			return
		}

		result, err := Db.ExecContext(ctx, "DELETE FROM orderitems WHERE id = $1", id)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete order item from database", "details": err.Error()})
			return
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve affected rows", "details": err.Error()})
			return
		}
		if rowsAffected == 0 {
			c.IndentedJSON(http.StatusNotFound, gin.H{"error": "No order item found with the given ID", "order_item_id": id})
			return
		}

		c.IndentedJSON(http.StatusOK, gin.H{"message": "Order item deleted successfully", "order_item_id": id})
	}
}

func GetOrderItemsByOrder() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		id := c.Param("order_id")
		if id == "" {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Order ID is required"})
			return
		}

		var orderItems []models.OrderItem

		rows, err := Db.QueryContext(ctx, "SELECT * FROM orderitems WHERE order_id = $1", id)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch order items", "details": err.Error()})
			return
		}

		for rows.Next() {
			var orderItem models.OrderItem
			if err := rows.Scan(&orderItem.ID, &orderItem.OrderID, &orderItem.FoodID, &orderItem.Quantity, &orderItem.UnitPrice, &orderItem.SubTotal, &orderItem.CreatedAt, &orderItem.UpdatedAt); err != nil {
				c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan order items", "details": err.Error()})
				return
			}
			orderItems = append(orderItems, orderItem)
		}

		c.IndentedJSON(http.StatusOK, gin.H{"message": "Order items by order_id fetched successfully", "order_id": id, "order_items": orderItems})
	}
}
