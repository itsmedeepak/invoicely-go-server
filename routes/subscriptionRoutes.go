package routes

import (
	"github.com/gin-gonic/gin"
	"tmp-invoicely.co/controllers"
)

func SubscriptionRoutes(incomingRoute *gin.Engine) {

	incomingRoute.GET("/subscription", controllers.GetSubscription())
	incomingRoute.POST("/subscription", controllers.UpdateSubscription())
}
