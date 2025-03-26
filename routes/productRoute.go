package routes

import (
	"github.com/gin-gonic/gin"
	"tmp-invoicely.co/controllers"
)

func ProductRoute(incomingRoute *gin.Engine) {
	incomingRoute.GET("/product", controllers.GetProducts())
	incomingRoute.GET("/product/:productId", controllers.GetProduct())
	incomingRoute.POST("/product", controllers.CreateProduct())
	incomingRoute.PUT("/product/:productId", controllers.UpdateProduct())
	incomingRoute.DELETE("/product/:productId", controllers.DeleteProduct())
}
