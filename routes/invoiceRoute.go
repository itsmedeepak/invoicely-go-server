package routes

import (
	"github.com/gin-gonic/gin"
	"tmp-invoicely.co/controllers"
)

func InvoiceRoute(incomingRoute *gin.Engine) {
	incomingRoute.GET("/invoice", controllers.GetInvoices())
	incomingRoute.GET("/invoice/:invoiceId", controllers.GetInvoice())
	incomingRoute.POST("/invoice", controllers.CreateInvoices())
	incomingRoute.GET("/invoice/send/:invoiceId", controllers.SendInvoiceEmail())
	incomingRoute.PUT("/invoice/:invoiceId", controllers.UpdateInvoices())
	incomingRoute.DELETE("/invoice/:invoiceId", controllers.DeleteInvoices())
}
