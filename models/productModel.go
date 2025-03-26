package models

import (
	"time"
)

type Product struct {
	ProductId    string    `json:"_id" bson:"_id"`
	UserId       string    `json:"user_id" bson:"user_id"`
	Name         string    `json:"name" bson:"name" validate:"required,min=2,max=20"`
	Currency     string    `json:"currency" bson:"currency" validate:"required,min=2,max=20"`
	Category     string    `json:"category" bson:"category" validate:"required"`
	Price        float64   `json:"price" bson:"price" validate:"required,gt=0"`
	ProductImage string    `json:"product_image" bson:"product_image"`
	Discount     float64   `json:"discount" bson:"discount" validate:"min=0,max=100"`
	CreatedAt    time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" bson:"updated_at"`
}
