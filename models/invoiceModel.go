package models

import (
	"time"
)

// ProductInInvoice represents a product entry in the invoice
type ProductInInvoice struct {
	ProductID  string  `json:"product_id" bson:"product_id"`
	Name       string  `json:"name,omitempty" bson:"name,omitempty"`   // Optional for DB, but useful for display
	Price      float64 `json:"price,omitempty" bson:"price,omitempty"` // Optional if needed
	Quantity   int     `json:"quantity" bson:"quantity"`
	Discount   float64 `json:"discount,omitempty" bson:"discount,omitempty"`       // Optional
	Currency   string  `json:"currency,omitempty" bson:"currency,omitempty"`       // Optional
	FinalPrice float64 `json:"final_price,omitempty" bson:"final_price,omitempty"` // Optional
	TotalPrice float64 `json:"total_price,omitempty" bson:"total_price,omitempty"` // Optional
}

// Invoice represents the structure of the invoice document in MongoDB
type Invoice struct {
	UserID             string             `json:"user_id" bson:"user_id"`
	InvoiceID          string             `json:"invoice_id" bson:"invoice_id"`
	InvoiceNo          string             `json:"invoice_no" bson:"invoice_no"`
	InvoiceUrl         string             `json:"invoice_url" bson:"invoice_url"`
	InvoiceLogo        string             `json:"invoice_logo" bson:"invoice_logo"`
	IssuedDate         string             `json:"issued_date" bson:"issued_date"`
	DueDate            string             `json:"due_date" bson:"due_date"`
	Customer           Customer           `json:"customer" bson:"customer"`
	CustomerID         string             `json:"customer_id" bson:"customer_id"`
	Products           []ProductInInvoice `json:"products" bson:"products"`
	PaymentMethod      string             `json:"payment_method" bson:"payment_method"`
	PaymentStatus      string             `json:"payment_status" bson:"payment_status"`
	InvoiceGeneratedBy string             `json:"invoice_generated_by" bson:"invoice_generated_by"`
	TotalAmount        float64            `json:"total_amount" bson:"total_amount"`
	Currency           string             `json:"currency" bson:"currency"`
	Url                string             `json:"url" bson:"url"`
	CreatedAt          time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt          time.Time          `json:"updated_at" bson:"updated_at"`
}
