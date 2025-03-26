package routes

import (
	"github.com/gin-gonic/gin"
	"tmp-invoicely.co/controllers"
)

func BillingRoute(incomingRoute *gin.Engine) {
	incomingRoute.GET("/billing", controllers.GetBilling())
	incomingRoute.POST("/billing", controllers.UpdateBilling())
}
