package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"tmp-invoicely.co/models"
	"tmp-invoicely.co/utils"
)

func UpdateInvoiceConfiguration() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		userID, exists := c.Get("user_id")
		if !exists {
			utils.ApiResponse(c, http.StatusUnauthorized, false, "Unauthorized", nil)
			return
		}

		var config models.InvoiceConfiguration
		if err := c.BindJSON(&config); err != nil {
			utils.ApiResponse(c, http.StatusBadRequest, false, "Invalid input data", nil)
			return
		}

		config.UserID = userID.(string)
		config.UpdatedAt = time.Now()

		filter := bson.M{"user_id": config.UserID}
		update := bson.M{"$set": config}

		opts := options.UpdateOne().SetUpsert(true)

		_, err := InvoiceConfigCollection.UpdateOne(ctx, filter, update, opts)
		if err != nil {
			utils.ApiResponse(c, http.StatusInternalServerError, false, "Failed to update invoice configuration", nil)
			return
		}

		utils.ApiResponse(c, http.StatusOK, true, "Invoice configuration saved", config)
	}
}

func GetInvoiceConfiguration() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		userID, exists := c.Get("user_id")
		if !exists {
			utils.ApiResponse(c, http.StatusUnauthorized, false, "Unauthorized", nil)
			return
		}

		var config models.InvoiceConfiguration
		err := InvoiceConfigCollection.FindOne(ctx, bson.M{"user_id": userID}).Decode(&config)
		if err != nil {
			utils.ApiResponse(c, http.StatusNotFound, false, "Invoice configuration not found", nil)
			return
		}

		utils.ApiResponse(c, http.StatusOK, true, "Invoice configuration fetched", config)
	}
}
