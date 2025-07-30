package controllers

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"restaurant-management/helpers"
	"restaurant-management/models"
	"time"

	"github.com/gin-gonic/gin"
)

func GetInvoices() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		id := c.Param("restaurant_id")
		if id == "" {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Please provide Restaurant ID"})
			return
		}

		query := `SELECT * FROM invoices WHERE restaurant_id = $1`
		rows, err := Db.QueryContext(ctx, query, id)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch invoices from database", "details": err.Error()})
			return
		}

		var invoices []models.Invoice
		for rows.Next() {
			var invoice models.Invoice
			if err := rows.Scan(&invoice.ID, &invoice.OrderID, &invoice.RestaurantID, &invoice.Amount, &invoice.Tax, &invoice.Total, &invoice.Status, &invoice.PaymentMethod, &invoice.CreatedAt, &invoice.UpdatedAt); err != nil {
				c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan invoice from database", "details": err.Error()})
				return
			}
			invoices = append(invoices, invoice)
		}

		c.IndentedJSON(http.StatusOK, gin.H{"message": "Invoices fetched successfully", "invoices": invoices})
	}
}

func GetInvoice() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Logic to get an invoice by ID
		c.JSON(200, gin.H{"message": "Get invoice by ID"})
	}
}

func CreateInvoice() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Logic to create a new invoice
		c.JSON(201, gin.H{"message": "Create invoice"})
	}
}

func UpdateInvoice() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Logic to update an existing invoice
		c.JSON(200, gin.H{"message": "Update invoice"})
	}
}
func DeleteInvoice() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Logic to delete an invoice
		c.JSON(200, gin.H{"message": "Delete invoice"})
	}
}

func GetInvoiceByOrderID() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Logic to get an invoice by order ID
		c.JSON(200, gin.H{"message": "Get invoice by order ID"})
	}
}

func CreateInvoiceFromOrder(order models.Order) error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	const (
		queryInsert = `
			INSERT INTO invoices (order_id, amount, tax, total, status, payment_method, restaurant_id)
			VALUES ($1, $2, $3, $4, $5, $6, $7)
			ON CONFLICT (order_id)
			DO UPDATE SET 
				amount = EXCLUDED.amount,
				tax = EXCLUDED.tax,
				total = EXCLUDED.total,
				status = EXCLUDED.status,
				payment_method = EXCLUDED.payment_method,
				restaurant_id = EXCLUDED.restaurant_id,
				updated_at = CURRENT_TIMESTAMP
		`
		queryDelete = `
			DELETE FROM invoices WHERE order_id = $1
		`
	)

	var tax = 0.00
	var total = order.TotalPrice + tax
	var status = "pending"
	var paymentMethod = "cash"

	var result sql.Result
	var err error

	switch order.Status {
	case "preparing", "ready", "served", "paid":
		if order.Status == "paid" {
			status = "paid"
		} else {
			status = "pending"
		}
		fmt.Println("Status: " + status)
		result, err = Db.ExecContext(ctx, queryInsert, order.ID, order.TotalPrice, tax, total, status, paymentMethod, order.RestaurantID)
	case "pending", "cancelled":
		result, err = Db.ExecContext(ctx, queryDelete, order.ID)
	}

	if err != nil {
		return err
	}

	affectedRows, err := result.RowsAffected()
	if err != nil || affectedRows == 0 {
		return err
	}

	return nil
}

func DownloadInvoice() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		invoiceID := c.Param("invoice_id")
		if invoiceID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invoice ID is required"})
			return
		}

		var invoice models.Invoice
		query := `SELECT * FROM invoices WHERE id = $1`
		err := Db.QueryRowContext(ctx, query, invoiceID).
			Scan(&invoice.ID, &invoice.OrderID, &invoice.RestaurantID, &invoice.Amount, &invoice.Tax, &invoice.Total,
				&invoice.Status, &invoice.PaymentMethod, &invoice.CreatedAt, &invoice.UpdatedAt)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch invoice", "details": err.Error()})
			return
		}

		var order models.Order
		err = Db.QueryRowContext(ctx, "SELECT * FROM orders WHERE id = $1", invoice.OrderID).
			Scan(&order.ID, &order.TableID, &order.RestaurantID, &order.OrderDate, &order.TotalPrice,
				&order.Status, &order.Notes, &order.CreatedAt, &order.UpdatedAt)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch order", "details": err.Error()})
			return
		}

		rows, err := Db.Query(`
			SELECT oi.id, oi.order_id, oi.food_id, f.name, oi.quantity, oi.unit_price, oi.subtotal, oi.created_at, oi.updated_at
			FROM orderitems oi
			JOIN foods f ON f.id = oi.food_id
			WHERE oi.order_id = $1
		`, order.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch order items", "details": err.Error()})
			return
		}
		defer rows.Close()

		for rows.Next() {
			var item models.OrderItem
			if err := rows.Scan(&item.ID, &item.OrderID, &item.FoodID, &item.FoodName, &item.Quantity, &item.UnitPrice, &item.SubTotal, &item.CreatedAt, &item.UpdatedAt); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan order item", "details": err.Error()})
				return
			}
			order.OrderItems = append(order.OrderItems, item)
		}

		var restaurant models.Restaurant
		err = Db.QueryRowContext(ctx, "SELECT * FROM restaurants WHERE id = $1", invoice.RestaurantID).
			Scan(&restaurant.ID, &restaurant.Name, &restaurant.OwnerID, &restaurant.Logo, &restaurant.Address, &restaurant.Description, &restaurant.CreatedAt, &restaurant.UpdatedAt)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch restaurant", "details": err.Error()})
			return
		}

		document, err := helpers.GeneratePdfFromData(invoice, order, restaurant)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate PDF", "details": err.Error()})
			return
		}

		bytes := document.GetBytes()

		c.Header("Content-Disposition", "attachment; filename=invoice-"+invoiceID+".pdf")
		c.Data(http.StatusOK, "application/pdf", bytes)
	}
}
