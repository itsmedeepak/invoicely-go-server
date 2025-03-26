package routes

import (
	"github.com/gin-gonic/gin"
	"tmp-invoicely.co/controllers"
)

func UserRoutes(incomingRoute *gin.Engine) {

	incomingRoute.GET("/profile", controllers.GetProfile())
	incomingRoute.PUT("/profile", controllers.EditProfile())
	incomingRoute.PUT("/change-password", controllers.ChangePassword())
}
