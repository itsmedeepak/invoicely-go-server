package routes

import (
	"github.com/gin-gonic/gin"
	"tmp-invoicely.co/controllers"
)

func InvoiceConfigRoute(incomingRoute *gin.Engine) {
	incomingRoute.GET("/invoice-config", controllers.GetInvoiceConfiguration())
	incomingRoute.POST("/invoice-config", controllers.UpdateInvoiceConfiguration())

}
