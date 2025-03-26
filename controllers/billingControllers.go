package controllers

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"tmp-invoicely.co/models"
	"tmp-invoicely.co/utils"
)

func GetBilling() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Println("BILLING")
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		userID, exists := c.Get("user_id")
		if !exists {
			utils.ApiResponse(c, http.StatusUnauthorized, false, "Unauthorized", nil)
			return
		}

		var billing models.Billing
		err := BillingCollection.FindOne(ctx, bson.M{"user_id": userID}).Decode(&billing)
		if err != nil {
			utils.ApiResponse(c, http.StatusOK, false, "Billing details not found", nil)
			return
		}

		utils.ApiResponse(c, http.StatusOK, true, "Billing details fetched", billing)
	}
}

func UpdateBilling() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		userID, exists := c.Get("user_id")
		if !exists {
			utils.ApiResponse(c, http.StatusUnauthorized, false, "Unauthorized", nil)
			return
		}

		var billing models.Billing
		if err := c.BindJSON(&billing); err != nil {
			utils.ApiResponse(c, http.StatusBadRequest, false, "Invalid input", nil)
			return
		}

		// Optional: Validate billing fields
		if err := Validate.Struct(billing); err != nil {
			utils.ApiResponse(c, http.StatusBadRequest, false, "Validation failed", err.Error())
			return
		}

		billing.UserID = userID.(string)
		billing.UpdatedAt = time.Now()

		filter := bson.M{"user_id": billing.UserID}
		update := bson.M{"$set": billing}
		opts := options.UpdateOne().SetUpsert(true)

		_, err := BillingCollection.UpdateOne(ctx, filter, update, opts)
		if err != nil {
			utils.ApiResponse(c, http.StatusInternalServerError, false, "Failed to update billing details", nil)
			return
		}

		utils.ApiResponse(c, http.StatusOK, true, "Billing details updated", billing)
	}
}
