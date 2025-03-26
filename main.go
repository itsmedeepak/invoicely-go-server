package main

import (
	"log"

	"github.com/gin-gonic/gin"

	middleware "tmp-invoicely.co/middlewares"
	"tmp-invoicely.co/routes"
	"tmp-invoicely.co/utils"
)

func main() {
	port := utils.GetEnv("PORT")
	if port == "" {
		port = "4040"
	}

	router := gin.Default() // Includes Logger and Recovery middleware

	router.Use(middleware.Authenticate())

	// Routes
	routes.AuthRoute(router)
	routes.CustomerRoute(router)
	routes.InvoiceRoute(router)
	routes.ProductRoute(router)
	routes.BillingRoute(router)
	routes.InvoiceConfigRoute(router)
	routes.UserRoutes(router)
	routes.SubscriptionRoutes(router)

	router.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{"error": "Not Found"})
	})

	log.Printf("Server running at http://localhost:%s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
