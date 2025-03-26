package controllers

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/bson"
	"tmp-invoicely.co/database"
	"tmp-invoicely.co/models"
	"tmp-invoicely.co/utils"
)

var Validate = validator.New()
var UserCollection = database.OpenCollection(database.ConnectDB(), "USER")

func SignUp() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var user models.User

		if err := c.BindJSON(&user); err != nil {
			utils.ApiResponse(c, http.StatusBadRequest, false, "Invalid input data", nil)
			return
		}

		if err := Validate.Struct(user); err != nil {
			utils.ApiResponse(c, http.StatusBadRequest, false, "Validation failed", err.Error())
			return
		}

		var existingUser models.User
		err := UserCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&existingUser)
		log.Println(existingUser)
		if err == nil {
			utils.ApiResponse(c, http.StatusConflict, false, "User already exists", nil)
			return
		}

		hashedPassword, hashErr := utils.HashPassword(user.Password)
		if hashErr != nil {
			utils.ApiResponse(c, http.StatusInternalServerError, false, "Error hashing password", nil)
			return
		}
		user.Password = string(hashedPassword)

		user.CreatedAt = time.Now()
		user.UpdatedAt = user.CreatedAt
		user.UserId = primitive.NewObjectID().Hex()

		_, insertErr := UserCollection.InsertOne(ctx, user)
		if insertErr != nil {
			utils.ApiResponse(c, http.StatusInternalServerError, false, "Error creating user", nil)
			return
		}

		trialSubscription := models.Subscription{
			UserID:            user.UserId,
			Plan:              "Trial",
			ValidTill:         time.Now().AddDate(0, 1, 0),
			Status:            "Active",
			CreditsUsed:       0,
			CreditsRemaining:  10,
			AverageDailyUsage: 0,
			LastRefreshed:     time.Now(),
			CreatedAt:         time.Now(),
			UpdatedAt:         time.Now(),
		}

		_, subErr := SubscriptionCollection.InsertOne(ctx, trialSubscription)
		if subErr != nil {
			utils.ApiResponse(c, http.StatusInternalServerError, false, "User created but failed to initialize subscription", nil)
			return
		}

		responseUser := map[string]interface{}{
			"_id":           user.UserId,
			"first_name":    user.FirstName,
			"last_name":     user.LastName,
			"email":         user.Email,
			"phone":         user.Phone,
			"auth_token":    user.AuthToken,
			"refresh_token": user.RefreshToken,
			"created_at":    user.CreatedAt,
			"updated_at":    user.UpdatedAt,
		}

		utils.ApiResponse(c, http.StatusCreated, true, "User created successfully", responseUser)
	}
}

func LogIn() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		// Struct for login request
		var loginRequest struct {
			Email    string `json:"email" binding:"required,email"`
			Password string `json:"password" binding:"required"`
		}

		// Bind JSON request
		if err := c.BindJSON(&loginRequest); err != nil {
			utils.ApiResponse(c, http.StatusBadRequest, false, "Invalid input data", nil)
			return
		}

		var user models.User
		err := UserCollection.FindOne(ctx, bson.M{"email": loginRequest.Email}).Decode(&user)
		if err != nil {
			utils.ApiResponse(c, http.StatusUnauthorized, false, "Invalid email or password", nil)
			return
		}

		if !utils.ComparePassword(loginRequest.Password, user.Password) {
			utils.ApiResponse(c, http.StatusUnauthorized, false, "Invalid email or password", nil)
			return
		}

		// Generate JWT token
		authToken, refreshToken, tokenErr := utils.GenerateToken(user.UserId, user.Email)
		if tokenErr != nil {
			utils.ApiResponse(c, http.StatusInternalServerError, false, "Error generating token", nil)
			return
		}

		// Update user tokens in DB
		updateData := bson.M{
			"$set": bson.M{
				"auth_token":    authToken,
				"refresh_token": refreshToken,
				"updated_at":    time.Now(),
			},
		}

		_, updateErr := UserCollection.UpdateOne(ctx, bson.M{"email": user.Email}, updateData)
		if updateErr != nil {
			utils.ApiResponse(c, http.StatusInternalServerError, false, "Failed to update tokens", nil)
			return
		}

		// Send response excluding password
		responseUser := gin.H{
			"_id":           user.UserId,
			"first_name":    user.FirstName,
			"last_name":     user.LastName,
			"email":         user.Email,
			"phone":         user.Phone,
			"auth_token":    authToken,
			"refresh_token": refreshToken,
			"created_at":    user.CreatedAt,
			"updated_at":    user.UpdatedAt,
		}

		utils.ApiResponse(c, http.StatusOK, true, "Login successful", responseUser)
	}
}

func ResetPassword() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func SendOtpEmail() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func SendHelloController() gin.HandlerFunc {
	return func(c *gin.Context) {
		email := "deep.bes.us@gmail.com"
		result, err := utils.SendHelloEmail(email)
		if err != nil {
			utils.ApiResponse(c, 500, false, "Error", err)
		}
		utils.ApiResponse(c, 200, true, "Message Send Success", result)
	}
}
