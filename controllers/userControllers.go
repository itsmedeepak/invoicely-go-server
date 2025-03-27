package controllers

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
	"tmp-invoicely.co/database"
	"tmp-invoicely.co/models"
	"tmp-invoicely.co/utils"
)

var BillingCollection = database.OpenCollection(database.ConnectDB(), "BILLING")

var SubscriptionCollection = database.OpenCollection(database.ConnectDB(), "SUBSCRIPTION")

var InvoiceConfigCollection = database.OpenCollection(database.ConnectDB(), "INVOICECONFIG")

func EditProfile() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		userID, exists := c.Get("user_id")
		if !exists {
			utils.ApiResponse(c, http.StatusUnauthorized, false, "Unauthorized", nil)
			return
		}

		var input struct {
			FirstName   string `json:"firstName" validate:"required,min=2,max=20"`
			LastName    string `json:"lastName" validate:"required,min=2,max=20"`
			WalkThrough bool   `json:"walk_through" bson:"walk_through"`
			Phone       string `json:"phone" validate:"required"`
		}

		if err := c.BindJSON(&input); err != nil {
			utils.ApiResponse(c, http.StatusBadRequest, false, "Invalid input", nil)
			return
		}

		if err := Validate.Struct(input); err != nil {
			utils.ApiResponse(c, http.StatusBadRequest, false, "Validation failed", err.Error())
			return
		}

		update := bson.M{
			"$set": bson.M{
				"first_name":   input.FirstName,
				"last_name":    input.LastName,
				"phone":        input.Phone,
				"walk_through": input.WalkThrough,
				"updated_at":   time.Now(),
			},
		}

		log.Println("INPUT", input)

		result, err := UserCollection.UpdateOne(ctx, bson.M{"_id": userID}, update)
		if err != nil || result.MatchedCount == 0 {
			utils.ApiResponse(c, http.StatusInternalServerError, false, "Failed to update profile", nil)
			return
		}

		utils.ApiResponse(c, http.StatusOK, true, "Profile updated", nil)
	}
}

func GetProfile() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		userID, exists := c.Get("user_id")
		if !exists {
			utils.ApiResponse(c, http.StatusUnauthorized, false, "Unauthorized", nil)
			return
		}

		var user models.User
		err := UserCollection.FindOne(ctx, bson.M{"_id": userID}).Decode(&user)
		if err != nil {
			utils.ApiResponse(c, http.StatusNotFound, false, "User not found", nil)
			return
		}

		responseUser := map[string]interface{}{
			"_id":          user.UserId,
			"first_name":   user.FirstName,
			"last_name":    user.LastName,
			"email":        user.Email,
			"walk_through": user.WalkThrough,
			"phone":        user.Phone,
			"created_at":   user.CreatedAt,
			"updated_at":   user.UpdatedAt,
		}

		utils.ApiResponse(c, http.StatusOK, true, "Profile fetched", responseUser)
	}
}

func ChangePassword() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		userID, exists := c.Get("user_id")
		if !exists {
			utils.ApiResponse(c, http.StatusUnauthorized, false, "Unauthorized", nil)
			return
		}

		var input struct {
			OldPassword string `json:"current_password" binding:"required"`
			NewPassword string `json:"new_password" binding:"required,min=8"`
		}

		if err := c.BindJSON(&input); err != nil {
			utils.ApiResponse(c, http.StatusBadRequest, false, "Invalid input", nil)
			return
		}

		var user models.User
		err := UserCollection.FindOne(ctx, bson.M{"_id": userID}).Decode(&user)
		if err != nil {
			utils.ApiResponse(c, http.StatusNotFound, false, "User not found", nil)
			return
		}

		if !utils.ComparePassword(input.OldPassword, user.Password) {
			utils.ApiResponse(c, http.StatusUnauthorized, false, "Incorrect old password", nil)
			return
		}

		hashedPassword, err := utils.HashPassword(input.NewPassword)
		if err != nil {
			utils.ApiResponse(c, http.StatusInternalServerError, false, "Error hashing new password", nil)
			return
		}

		update := bson.M{
			"$set": bson.M{
				"password":   hashedPassword,
				"updated_at": time.Now(),
			},
		}

		_, err = UserCollection.UpdateOne(ctx, bson.M{"_id": userID}, update)
		if err != nil {
			utils.ApiResponse(c, http.StatusInternalServerError, false, "Failed to update password", nil)
			return
		}

		utils.ApiResponse(c, http.StatusOK, true, "Password changed successfully", nil)
	}
}
