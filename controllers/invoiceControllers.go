package controllers

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"tmp-invoicely.co/database"
	"tmp-invoicely.co/models"
	"tmp-invoicely.co/utils"
)

var InvoiceCollection = database.OpenCollection(database.ConnectDB(), "INVOICES")

func GetInvoices() gin.HandlerFunc {
	return func(c *gin.Context) {

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

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		var invoices []models.Invoice
		cursor, err := InvoiceCollection.Find(ctx, bson.M{"user_id": userIdStr})
		if err != nil {
			log.Printf("Error fetching invoices: %v", err)
			utils.ApiResponse(c, http.StatusInternalServerError, false, "Failed to fetch invoices", nil)
			return
		}

		if err = cursor.All(ctx, &invoices); err != nil {
			log.Printf("Error parsing invoices: %v", err)
			utils.ApiResponse(c, http.StatusInternalServerError, false, "Failed to parse invoices", nil)
			return
		}

		log.Printf("Fetched Invoices: %+v", invoices)

		if len(invoices) == 0 {
			utils.ApiResponse(c, http.StatusOK, true, "No invoices found", []models.Invoice{})
			return
		}

		utils.ApiResponse(c, http.StatusOK, true, "Invoices fetched successfully", invoices)
	}
}

func GetInvoice() gin.HandlerFunc {
	return func(c *gin.Context) {
		invoiceId := c.Param("invoiceId")
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

		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var invoice models.Invoice
		err := InvoiceCollection.FindOne(ctx, bson.M{"_id": invoiceId, "user_id": userIdStr}).Decode(&invoice)

		if err == mongo.ErrNoDocuments {
			utils.ApiResponse(c, http.StatusNotFound, false, "Invoice not found", nil)
			return
		} else if err != nil {
			utils.ApiResponse(c, http.StatusInternalServerError, false, "Failed to fetch invoice", nil)
			return
		}

		utils.ApiResponse(c, http.StatusOK, true, "Invoice fetched successfully", invoice)
	}
}

func CreateInvoices() gin.HandlerFunc {
	return func(c *gin.Context) {
		var invoice models.Invoice

		// Bind JSON to the Invoice struct
		if err := c.ShouldBindJSON(&invoice); err != nil {
			utils.ApiResponse(c, http.StatusBadRequest, false, "Invalid request body", nil)
			return
		}

		// Validate total amount
		if invoice.TotalAmount <= 0 {
			utils.ApiResponse(c, http.StatusBadRequest, false, "Total amount must be greater than zero", nil)
			return
		}

		// Extract user ID from context
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

		// Context with timeout for MongoDB operations
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		// Fetch user's subscription
		var subscription models.Subscription
		err := SubscriptionCollection.FindOne(ctx, bson.M{"user_id": userIdStr}).Decode(&subscription)
		if err != nil {
			utils.ApiResponse(c, http.StatusNotFound, false, "Subscription not found", nil)
			return
		}

		// Check if the user has enough credits
		if subscription.CreditsRemaining <= 0 {
			utils.ApiResponse(c, http.StatusPaymentRequired, false, "Insufficient credits to create an invoice", nil)
			return
		}

		// Set invoice details
		invoice.InvoiceID = primitive.NewObjectID().Hex()
		invoice.UserID = userIdStr
		invoice.CreatedAt = time.Now()
		invoice.UpdatedAt = time.Now()

		// Fetch customer details
		var customer models.Customer
		err = CustomerCollection.FindOne(ctx, bson.M{"_id": invoice.CustomerID}).Decode(&customer)
		if err != nil {
			utils.ApiResponse(c, http.StatusInternalServerError, false, "Failed to fetch customer details", nil)
			return
		}
		invoice.Customer = customer

		// Insert invoice into MongoDB
		result, err := InvoiceCollection.InsertOne(ctx, invoice)
		if err != nil {
			utils.ApiResponse(c, http.StatusInternalServerError, false, "Failed to create invoice", nil)
			return
		}

		// **Update subscription: decrement credits**
		update := bson.M{
			"$inc": bson.M{
				"credits_remaining": -1,
				"credits_used":      1,
			},
			"$set": bson.M{
				"updated_at": time.Now(),
			},
		}

		_, err = SubscriptionCollection.UpdateOne(ctx, bson.M{"user_id": userIdStr}, update)
		if err != nil {
			utils.ApiResponse(c, http.StatusInternalServerError, false, "Failed to update subscription", nil)
			return
		}

		// Fetch invoice configuration
		var config models.InvoiceConfiguration
		err = InvoiceConfigCollection.FindOne(ctx, bson.M{"user_id": userIdStr}).Decode(&config)
		if err != nil {
			utils.ApiResponse(c, http.StatusNotFound, false, "Invoice configuration not found", nil)
			return
		}

		// Send invoice email
		res, err := utils.SendInvoiceEmail(invoice, config)
		if err != nil {
			utils.ApiResponse(c, http.StatusInternalServerError, false, "Failed to send invoice email", nil)
			return
		}

		log.Println(res)
		utils.ApiResponse(c, http.StatusCreated, true, "Invoice created successfully", result.InsertedID)
	}
}

func UpdateInvoices() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the invoiceId from the URL parameter
		invoiceId := c.Param("invoiceId")

		// Ensure the user is authenticated
		userId, exists := c.Get("user_id")
		if !exists {
			utils.ApiResponse(c, http.StatusUnauthorized, false, "User not authenticated", nil)
			return
		}

		userIdStr, ok := userId.(string)
		if !ok {
			utils.ApiResponse(c, http.StatusUnauthorized, false, "User ID format is invalid", nil)
			return
		}

		// Bind the request body to the invoice struct
		var invoice models.Invoice
		if err := c.ShouldBindJSON(&invoice); err != nil {
			utils.ApiResponse(c, http.StatusBadRequest, false, "Invalid request body", nil)
			return
		}

		// Set the UpdatedAt field
		invoice.UpdatedAt = time.Now()

		// Create an update operation
		updateFields := bson.M{
			"$set": bson.M{
				"products":             invoice.Products, // Only update the necessary fields
				"payment_method":       invoice.PaymentMethod,
				"payment_status":       invoice.PaymentStatus,
				"invoice_generated_by": invoice.InvoiceGeneratedBy,
				"total_amount":         invoice.TotalAmount,
				"currency":             invoice.Currency,
				"updated_at":           invoice.UpdatedAt,
			},
		}

		// MongoDB context with timeout
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		// Update the invoice document in MongoDB by matching invoiceId and userId
		result, err := InvoiceCollection.UpdateOne(ctx, bson.M{
			"invoice_id": invoiceId,
			"user_id":    userIdStr,
		}, updateFields)

		if err != nil {
			log.Printf("Error updating invoice: %v", err) // Log error with context
			utils.ApiResponse(c, http.StatusInternalServerError, false, "Failed to update invoice", nil)
			return
		}

		// Check if any document was modified
		if result.ModifiedCount == 0 {
			utils.ApiResponse(c, http.StatusNotFound, false, "Invoice not found or no changes made", nil)
			return
		}

		// Return success response
		utils.ApiResponse(c, http.StatusOK, true, "Invoice updated successfully", result.ModifiedCount)
	}
}

func DeleteInvoices() gin.HandlerFunc {
	return func(c *gin.Context) {
		invoiceId := c.Param("invoiceId")
		userId, exists := c.Get("user_id")

		log.Println("Delete Invoice")
		if !exists {
			utils.ApiResponse(c, http.StatusUnauthorized, false, "User not authenticated", nil)
			return
		}

		userIdStr, _ := userId.(string)

		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		log.Println(invoiceId, userId)

		result, _ := InvoiceCollection.DeleteOne(ctx, bson.M{"invoice_id": invoiceId, "user_id": userIdStr})
		log.Println(result)
		if result.DeletedCount == 0 {
			utils.ApiResponse(c, http.StatusNotFound, false, "Invoice not found", nil)
			return
		}

		utils.ApiResponse(c, http.StatusOK, true, "Invoice deleted successfully", nil)
	}
}

func SendInvoiceEmail() gin.HandlerFunc {
	return func(c *gin.Context) {
		invoiceId := c.Param("invoiceId")
		userId, _ := c.Get("user_id")

		utils.ApiResponse(c, http.StatusOK, true, "Email sent successfully for invoice "+invoiceId, userId)
	}
}

func GetReports() gin.HandlerFunc {
	return func(c *gin.Context) {}
}

func Dashboard() gin.HandlerFunc {
	return func(c *gin.Context) {}
}
