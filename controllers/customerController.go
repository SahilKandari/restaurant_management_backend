package controllers

import (
	"context"
	"database/sql"
	"net/http"
	"time"

	"restaurant-management/models"

	"github.com/gin-gonic/gin"
)

func CustomerGetRestaurant() gin.HandlerFunc {
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

func CustomerGetTable() gin.HandlerFunc {
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

func CustomerGetMenus() gin.HandlerFunc {
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

func CustomerCreateOrderId() gin.HandlerFunc {
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

func CustomerUpdateOrder() gin.HandlerFunc {
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

func CustomerGetFoodsByRestaurantID() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
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

func CustomerCreateOrderItem() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
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

func CustomerUpdateOrderItem() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
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

func CustomerDeleteOrderItem() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
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
