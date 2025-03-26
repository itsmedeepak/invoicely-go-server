package routes

import (
	"github.com/gin-gonic/gin"
	"tmp-invoicely.co/controllers"
)

func CustomerRoute(incomingRoute *gin.Engine) {
	incomingRoute.GET("/customer", controllers.GetCustomers())
	incomingRoute.GET("/customer/:customerId", controllers.GetCustomer())
	incomingRoute.POST("/customer", controllers.CreateCustomer())
	incomingRoute.PUT("/customer/:customerId", controllers.UpdateCustomer())
	incomingRoute.DELETE("/customer/:customerId", controllers.DeleteCustomer())
}
