package routes

import (
	"restaurant-management/controllers"

	"github.com/gin-gonic/gin"
)

func InvoiceRoutes(incomingRoutes gin.IRoutes) {
	incomingRoutes.GET("/invoices/:restaurant_id", controllers.GetInvoices())
	incomingRoutes.GET("/invoice-pdf/:invoice_id", controllers.DownloadInvoice())
	// incomingRoutes.GET("/invoices/:invoice_id", controllers.GetInvoice())
	// incomingRoutes.POST("/invoices", controllers.CreateInvoice())
	// incomingRoutes.PATCH("/invoices/:invoice_id", controllers.UpdateInvoice())
	// incomingRoutes.DELETE("/invoices/:invoice_id", controllers.DeleteInvoice())
}
