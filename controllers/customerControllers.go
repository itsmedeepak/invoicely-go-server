package controllers

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"tmp-invoicely.co/database"
	"tmp-invoicely.co/models"
	"tmp-invoicely.co/utils"
)

var CustomerCollection = database.OpenCollection(database.ConnectDB(), "CUSTOMER")
var Validater = validator.New()

func GetCustomers() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Println("user_id")
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		userId, exists := c.Get("user_id")

		if !exists {
			utils.ApiResponse(c, http.StatusUnauthorized, false, "User not authenticated", nil)
			return
		}

		userIdStr, ok := userId.(string)
		log.Println("userIdStr:  ", userIdStr)
		if !ok {
			utils.ApiResponse(c, http.StatusInternalServerError, false, "Failed to parse user ID", nil)
			return
		}

		var customers []models.Customer

		cursor, err := CustomerCollection.Find(ctx, bson.M{"user_id": userIdStr})

		if err != nil {
			utils.ApiResponse(c, http.StatusInternalServerError, false, "Failed to fetch customers", nil)
			return
		}

		if err = cursor.All(ctx, &customers); err != nil {
			utils.ApiResponse(c, http.StatusInternalServerError, false, "Failed to parse customers", nil)
			return
		}

		utils.ApiResponse(c, http.StatusOK, true, "Customers fetched successfully", customers)
	}
}

func GetCustomer() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		customerId := c.Param("customerId")
		log.Println("customerId: ", customerId)

		if customerId == "" {
			utils.ApiResponse(c, http.StatusBadRequest, false, "Invalid customer ID", nil)
			return
		}

		userId, exists := c.Get("user_id")

		if !exists {
			utils.ApiResponse(c, http.StatusUnauthorized, false, "User not authenticated", nil)
			return
		}

		userIdStr, ok := userId.(string)
		log.Println("userIdStr:  ", userIdStr)
		if !ok {
			utils.ApiResponse(c, http.StatusInternalServerError, false, "Failed to parse user ID", nil)
			return
		}

		var customer models.Customer
		err := CustomerCollection.FindOne(ctx, bson.M{"_id": customerId, "user_id": userIdStr}).Decode(&customer)
		if err != nil {
			log.Println("Error:", err.Error())
			utils.ApiResponse(c, http.StatusOK, false, "Customer not found", nil)
			return
		}

		utils.ApiResponse(c, http.StatusOK, true, "Customer fetched successfully", customer)
	}
}

func CreateCustomer() gin.HandlerFunc {
	return func(c *gin.Context) {
		var customer models.Customer
		if err := c.BindJSON(&customer); err != nil {
			utils.ApiResponse(c, http.StatusBadRequest, false, "Invalid request body", nil)
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		userId, exists := c.Get("user_id")
		if !exists {
			utils.ApiResponse(c, http.StatusUnauthorized, false, "User not authenticated", nil)
			return
		}

		userIdStr, ok := userId.(string)
		if !ok {
			utils.ApiResponse(c, http.StatusInternalServerError, false, "Failed to parse user ID", nil)
			return
		}

		var existingCustomer models.Customer
		err := CustomerCollection.FindOne(ctx, bson.M{"email": customer.Email, "user_id": userIdStr}).Decode(&existingCustomer)
		if err == nil {
			utils.ApiResponse(c, http.StatusConflict, false, "Customer with this email already exists", nil)
			return
		}

		customerId := primitive.NewObjectID().Hex()
		customer.CustomerId = customerId

		customer.UserId = userIdStr
		customer.CreatedAt = time.Now()
		customer.UpdatedAt = customer.CreatedAt

		result, err := CustomerCollection.InsertOne(ctx, customer)
		if err != nil {
			utils.ApiResponse(c, http.StatusInternalServerError, false, "Failed to create customer", nil)
			return
		}

		utils.ApiResponse(c, http.StatusCreated, true, "Customer created successfully", result.InsertedID)
	}
}

func UpdateCustomer() gin.HandlerFunc {
	return func(c *gin.Context) {
		customerId := c.Param("customerId")
		log.Println("customerId: ", customerId)

		var customer models.Customer
		if err := c.ShouldBindJSON(&customer); err != nil {
			utils.ApiResponse(c, http.StatusBadRequest, false, "Invalid request body", nil)
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		customer.UpdatedAt = time.Now()
		userId, exists := c.Get("user_id")
		if !exists {
			utils.ApiResponse(c, http.StatusUnauthorized, false, "User not authenticated", nil)
			return
		}

		userIdStr, ok := userId.(string)
		if !ok {
			utils.ApiResponse(c, http.StatusInternalServerError, false, "Failed to parse user ID", nil)
			return
		}
		customer.UserId = userIdStr
		customer.CustomerId = customerId

		update := bson.M{"$set": customer}
		result, err := CustomerCollection.UpdateOne(ctx, bson.M{"_id": customerId, "user_id": userIdStr}, update)

		log.Printf("Result : ", result)

		if err != nil {
			utils.ApiResponse(c, http.StatusInternalServerError, false, "Failed to update customer", nil)
			return
		}

		if result.ModifiedCount == 0 {
			utils.ApiResponse(c, http.StatusNotFound, false, "Customer not found", nil)
			return
		}

		utils.ApiResponse(c, http.StatusOK, true, "Customer updated successfully", result.ModifiedCount)
	}
}

func DeleteCustomer() gin.HandlerFunc {
	return func(c *gin.Context) {
		customerId := c.Param("customerId")
		if customerId == "" {
			utils.ApiResponse(c, http.StatusBadRequest, false, "Customer ID is required", nil)
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		userId, exists := c.Get("user_id")
		if !exists {
			utils.ApiResponse(c, http.StatusUnauthorized, false, "User not authenticated", nil)
			return
		}

		userIdStr, ok := userId.(string)

		if !ok {
			utils.ApiResponse(c, http.StatusInternalServerError, false, "Failed to parse user ID", nil)
			return
		}

		result, err := CustomerCollection.DeleteOne(ctx, bson.M{"_id": customerId, "user_id": userIdStr})
		if err != nil {
			utils.ApiResponse(c, http.StatusInternalServerError, false, "Failed to delete customer", nil)
			return
		}

		if result.DeletedCount == 0 {
			utils.ApiResponse(c, http.StatusNotFound, false, "Customer not found", nil)
			return
		}

		utils.ApiResponse(c, http.StatusOK, true, "Customer deleted successfully", result.DeletedCount)
	}
}
