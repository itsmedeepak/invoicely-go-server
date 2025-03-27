package routes

import (
	"log"

	"github.com/gin-gonic/gin"
	"tmp-invoicely.co/controllers"
)

func AuthRoute(incomingRoute *gin.Engine) {
	log.Println("signUp")
	incomingRoute.POST("/auth/sign-up", controllers.SignUp())
	incomingRoute.POST("/auth/sign-in", controllers.LogIn())
	incomingRoute.POST("/auth/forgot-password", controllers.ResetPassword())
	incomingRoute.POST("/auth/send-otp", controllers.SendOtpEmail())
	// incomingRoute.POST("/auth/test", controllers.SendHelloController())
}
