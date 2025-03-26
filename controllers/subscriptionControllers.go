package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
	"tmp-invoicely.co/models"
	"tmp-invoicely.co/utils"
)

func GetSubscription() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		userID, exists := c.Get("user_id")
		if !exists {
			utils.ApiResponse(c, http.StatusUnauthorized, false, "Unauthorized", nil)
			return
		}

		var sub models.Subscription
		err := SubscriptionCollection.FindOne(ctx, bson.M{"user_id": userID}).Decode(&sub)
		if err != nil {
			utils.ApiResponse(c, http.StatusNotFound, false, "No subscription found", nil)
			return
		}

		utils.ApiResponse(c, http.StatusOK, true, "Subscription fetched", sub)
	}
}
func UpdateSubscription() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		userID, exists := c.Get("user_id")
		if !exists {
			utils.ApiResponse(c, http.StatusUnauthorized, false, "Unauthorized", nil)
			return
		}

		var subInput models.Subscription
		if err := c.BindJSON(&subInput); err != nil {
			utils.ApiResponse(c, http.StatusBadRequest, false, "Invalid input data", nil)
			return
		}

		// Check if subscription already exists
		var existing models.Subscription
		err := SubscriptionCollection.FindOne(ctx, bson.M{"user_id": userID}).Decode(&existing)

		if err != nil { // Create new trial if not found
			sub := models.Subscription{
				UserID:            userID.(string),
				Plan:              "trial",
				ValidTill:         time.Now().AddDate(0, 1, 0),
				Status:            "active",
				CreditsUsed:       0,
				CreditsRemaining:  10,
				AverageDailyUsage: 0,
				LastRefreshed:     time.Now(),
				CreatedAt:         time.Now(),
				UpdatedAt:         time.Now(),
			}

			_, insertErr := SubscriptionCollection.InsertOne(ctx, sub)
			if insertErr != nil {
				utils.ApiResponse(c, http.StatusInternalServerError, false, "Failed to create subscription", nil)
				return
			}

			utils.ApiResponse(c, http.StatusCreated, true, "Trial subscription created", sub)
			return
		}

		// Update existing subscription
		subInput.UserID = userID.(string)
		subInput.UpdatedAt = time.Now()

		update := bson.M{"$set": subInput}
		_, updateErr := SubscriptionCollection.UpdateOne(ctx, bson.M{"user_id": userID}, update)
		if updateErr != nil {
			utils.ApiResponse(c, http.StatusInternalServerError, false, "Failed to update subscription", nil)
			return
		}

		utils.ApiResponse(c, http.StatusOK, true, "Subscription updated", subInput)
	}
}
