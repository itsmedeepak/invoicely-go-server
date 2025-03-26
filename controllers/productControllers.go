package controllers

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"tmp-invoicely.co/database"
	"tmp-invoicely.co/models"
	"tmp-invoicely.co/utils"
)

var ProductCollection = database.OpenCollection(database.ConnectDB(), "PRODUCTS")

func GetProducts() gin.HandlerFunc {
	return func(c *gin.Context) {
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

		var products []models.Product
		cursor, err := ProductCollection.Find(ctx, bson.M{"user_id": userIdStr})
		if err != nil {
			utils.ApiResponse(c, http.StatusInternalServerError, false, "Failed to fetch products", nil)
			return
		}
		defer cursor.Close(ctx)

		if err = cursor.All(ctx, &products); err != nil {
			utils.ApiResponse(c, http.StatusInternalServerError, false, "Failed to parse products", nil)
			return
		}

		utils.ApiResponse(c, http.StatusOK, true, "Products fetched successfully", products)
	}
}

func GetProduct() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		productId := c.Param("productId")
		if productId == "" {
			utils.ApiResponse(c, http.StatusBadRequest, false, "Product ID is required", nil)
			return
		}

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

		var product models.Product
		err := ProductCollection.FindOne(ctx, bson.M{"_id": productId, "user_id": userIdStr}).Decode(&product)
		if err != nil {
			utils.ApiResponse(c, http.StatusOK, true, "Product not found", nil)
			return

		}

		utils.ApiResponse(c, http.StatusOK, true, "Product fetched successfully", product)
	}
}

func CreateProduct() gin.HandlerFunc {
	return func(c *gin.Context) {
		var product models.Product
		if err := c.ShouldBindJSON(&product); err != nil {
			utils.ApiResponse(c, http.StatusBadRequest, false, "Invalid request body", nil)
			return
		}

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

		product.ProductId = primitive.NewObjectID().Hex()
		product.UserId = userIdStr
		product.CreatedAt = time.Now()
		product.UpdatedAt = time.Now()

		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		result, err := ProductCollection.InsertOne(ctx, product)
		if err != nil {
			utils.ApiResponse(c, http.StatusInternalServerError, false, "Failed to create product", nil)
			return
		}

		utils.ApiResponse(c, http.StatusCreated, true, "Product created successfully", result.InsertedID)
	}
}

func UpdateProduct() gin.HandlerFunc {
	return func(c *gin.Context) {
		productId := c.Param("productId")
		var product models.Product

		if err := c.ShouldBindJSON(&product); err != nil {
			utils.ApiResponse(c, http.StatusBadRequest, false, "Invalid request body", nil)
			return
		}

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

		product.UpdatedAt = time.Now()
		product.ProductId = productId
		product.UserId = userIdStr
		update := bson.M{"$set": product}

		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		result, err := ProductCollection.UpdateOne(ctx, bson.M{"_id": productId, "user_id": userIdStr}, update)

		log.Println(result, err, productId, userIdStr)
		if err != nil {
			utils.ApiResponse(c, http.StatusInternalServerError, false, "Failed to update product", nil)
			return
		}

		if result.ModifiedCount == 0 {
			utils.ApiResponse(c, http.StatusNotFound, false, "Product not found", nil)
			return
		}

		utils.ApiResponse(c, http.StatusOK, true, "Product updated successfully", result.ModifiedCount)
	}
}

func DeleteProduct() gin.HandlerFunc {
	return func(c *gin.Context) {
		productId := c.Param("productId")
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

		result, err := ProductCollection.DeleteOne(ctx, bson.M{"_id": productId, "user_id": userIdStr})
		if err != nil {
			utils.ApiResponse(c, http.StatusInternalServerError, false, "Failed to delete product", nil)
			return
		}

		if result.DeletedCount == 0 {
			utils.ApiResponse(c, http.StatusNotFound, false, "Product not found", nil)
			return
		}

		utils.ApiResponse(c, http.StatusOK, true, "Product deleted successfully", nil)
	}
}
