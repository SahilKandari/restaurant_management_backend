package controllers

import (
	"context"
	"net/http"
	"restaurant-management/models"
	"time"

	"github.com/gin-gonic/gin"
)

func GetOrders() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		id := c.Query("restaurant_id")
		if id == "" {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Please provide the Restaurant ID"})
			return
		}

		var orders []models.Order
		var orderItems []models.OrderItem

		// Order by order_date DESC to get latest orders first
		rows, err := Db.QueryContext(ctx, "SELECT * FROM orders WHERE restaurant_id = $1 ORDER BY order_date DESC", id)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch orders from database", "details": err.Error()})
			return
		}
		defer rows.Close()

		for rows.Next() {
			var order models.Order
			if err := rows.Scan(&order.ID, &order.TableID, &order.RestaurantID, &order.OrderDate, &order.TotalPrice, &order.Status, &order.Notes, &order.CreatedAt, &order.UpdatedAt); err != nil {
				c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan order data", "details": err.Error()})
				return
			}
			orders = append(orders, order)
		}
		rowsOrderItems, err := Db.QueryContext(ctx, "SELECT * FROM orderitems WHERE order_id IN (SELECT id FROM orders WHERE restaurant_id = $1)", id)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch order items from database", "details": err.Error()})
			return
		}
		defer rowsOrderItems.Close()
		for rowsOrderItems.Next() {
			var orderItem models.OrderItem
			if err := rowsOrderItems.Scan(&orderItem.ID, &orderItem.OrderID, &orderItem.FoodID, &orderItem.Quantity, &orderItem.UnitPrice, &orderItem.SubTotal, &orderItem.CreatedAt, &orderItem.UpdatedAt); err != nil {
				c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan order item data", "details": err.Error()})
				return
			}
			orderItems = append(orderItems, orderItem)
		}

		for i, order := range orders {
			for _, orderItem := range orderItems {
				if order.ID == orderItem.OrderID {
					orders[i].OrderItems = append(orders[i].OrderItems, orderItem)
				}
			}
		}

		c.IndentedJSON(http.StatusOK, gin.H{"message": "Orders fetched successfully", "orders": orders})
	}
}

func GetOrder() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		id := c.Param("order_id")
		if id == "" {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Order ID is required"})
			return
		}

		var order models.Order
		var orderItems []models.OrderItem
		// var totalPrice float64

		if err := Db.QueryRowContext(ctx, "SELECT * FROM orders WHERE id = $1", id).
			Scan(&order.ID, &order.TableID, &order.RestaurantID, &order.OrderDate, &order.TotalPrice, &order.Status, &order.Notes, &order.CreatedAt, &order.UpdatedAt); err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch order from database", "details": err.Error()})
			return
		}

		rows, err := Db.QueryContext(ctx, "SELECT * FROM orderitems WHERE order_id = $1", id)
		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Failed to fetch order items related to order: " + id, "details": err.Error()})
			return
		}
		defer rows.Close()

		for rows.Next() {
			var orderItem models.OrderItem
			if err := rows.Scan(&orderItem.ID, &orderItem.OrderID, &orderItem.FoodID, &orderItem.Quantity, &orderItem.UnitPrice, &orderItem.SubTotal, &orderItem.CreatedAt, &orderItem.UpdatedAt); err != nil {
				c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan the order items where order is " + id, "details": err.Error()})
				return
			}
			// totalPrice += orderItem.SubTotal
			orderItems = append(orderItems, orderItem)
		}
		order.OrderItems = orderItems

		c.IndentedJSON(http.StatusOK, gin.H{"message": "Order fetched successfully", "order": order})
	}
}

func CreateOrder() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		var order models.Order

		if err := c.BindJSON(&order); err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Please provide correct data to create order", "details": err.Error()})
			return
		}

		query := `
			INSERT INTO orders (table_id, restaurant_id, order_date, status, notes)
			VALUES ($1, $2, $3, $4, $5)
			RETURNING id, table_id, restaurant_id, order_date, total_price, status, notes, created_at, updated_at
		`
		if err := Db.QueryRowContext(ctx, query, order.TableID, order.RestaurantID, order.OrderDate, order.Status, order.Notes).
			Scan(&order.ID, &order.TableID, &order.RestaurantID, &order.OrderDate, &order.TotalPrice, &order.Status, &order.Notes, &order.CreatedAt, &order.UpdatedAt); err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Failed to create order in database", "details": err.Error()})
			return
		}

		c.IndentedJSON(http.StatusCreated, gin.H{"message": "Order created successfully", "order": order})
	}
}

func CreateOrderId() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		var id int
		if err := Db.QueryRowContext(ctx, "INSERT INTO orders DEFAULT VALUES RETURNING id").Scan(&id); err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Failed to created Order ID", "details": err.Error()})
			return
		}

		c.IndentedJSON(http.StatusOK, gin.H{"message": "Order ID created successfully", "order_id": id})
	}
}

func UpdateOrder() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		id := c.Param("order_id")
		if id == "" {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Order ID is required"})
			return
		}

		var order models.Order
		var orderItems []models.OrderItem
		var totalPrice float64

		if err := c.BindJSON(&order); err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Please provide correct data to update order", "details": err.Error()})
			return
		}

		rows, err := Db.QueryContext(ctx, "SELECT * FROM orderitems WHERE order_id = $1", id)
		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Failed to fetch order items related to order: " + id, "details": err.Error()})
			return
		}
		defer rows.Close()

		for rows.Next() {
			var orderItem models.OrderItem
			if err := rows.Scan(&orderItem.ID, &orderItem.OrderID, &orderItem.FoodID, &orderItem.Quantity, &orderItem.UnitPrice, &orderItem.SubTotal, &orderItem.CreatedAt, &orderItem.UpdatedAt); err != nil {
				c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan the order items where order is " + id, "details": err.Error()})
				return
			}
			totalPrice += orderItem.SubTotal
			orderItems = append(orderItems, orderItem)
		}

		if totalPrice <= 0 {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Failed to add total price of the order items where order is " + id})
			return
		}

		order.OrderItems = orderItems
		order.TotalPrice = totalPrice

		query := `
			UPDATE orders
			SET table_id = $1, restaurant_id = $2, order_date = $3, status = $4, total_price = $5, notes = $6, updated_at = CURRENT_TIMESTAMP
			WHERE id = $7
			RETURNING id, table_id, restaurant_id, order_date, total_price, status, notes, created_at, updated_at
		`
		row := Db.QueryRowContext(ctx, query, order.TableID, order.RestaurantID, order.OrderDate, order.Status, order.TotalPrice, order.Notes, id)
		if err := row.Scan(&order.ID, &order.TableID, &order.RestaurantID, &order.OrderDate, &order.TotalPrice, &order.Status, &order.Notes, &order.CreatedAt, &order.UpdatedAt); err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Failed to update the order in database", "details": err.Error()})
			return
		}

		c.IndentedJSON(http.StatusOK, gin.H{"message": "Order updated successfully", "order": order})
	}
}

func UpdateOrderStatus() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		id := c.Param("order_id")
		if id == "" {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Order ID is required"})
			return
		}

		var orderStatus models.OrderStatus
		var order models.Order
		if err := c.BindJSON(&orderStatus); err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Please provide the correct status to update order status", "details": err.Error()})
			return
		}

		query := `
			UPDATE orders
			SET status = $1, updated_at = CURRENT_TIMESTAMP
			WHERE id = $2
			RETURNING id, table_id, restaurant_id, order_date, total_price, status, notes, created_at, updated_at			
		`
		if err := Db.QueryRowContext(ctx, query, orderStatus.Status, id).Scan(&order.ID, &order.TableID, &order.RestaurantID, &order.OrderDate, &order.TotalPrice, &order.Status, &order.Notes, &order.CreatedAt, &order.UpdatedAt); err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Failed to update status in database", "details": err.Error()})
			return
		}

		if err := CreateInvoiceFromOrder(order); err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Failed to create invoices for the order in database", "details": err.Error()})
			return
		}

		c.IndentedJSON(http.StatusOK, gin.H{"message": "Order updated successfully", "order_id": id})
	}
}

func DeleteOrder() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		id := c.Param("order_id")
		if id == "" {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Order ID is required"})
			return
		}

		result, err := Db.ExecContext(ctx, "DELETE FROM orders WHERE id = $1", id)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete order from database", "details": err.Error()})
			return
		}
		rowsAffected, err := result.RowsAffected()
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve affected rows for orders", "details": err.Error()})
			return
		}
		if rowsAffected == 0 {
			c.IndentedJSON(http.StatusNotFound, gin.H{"error": "No order found with given ID"})
			return
		}

		c.IndentedJSON(http.StatusNoContent, gin.H{"message": "Order and it's associated order items deleted successfully", "order_id": id})
	}
}
