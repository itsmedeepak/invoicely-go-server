package models

import "time"

type Billing struct {
	UserID       string    `json:"user_id" bson:"user_id"`
	FullName     string    `json:"full_name" bson:"full_name" validate:"required"`
	Email        string    `json:"email" bson:"email" validate:"required"`
	AddressLine1 string    `json:"address_line1" bson:"address_line1" validate:"required"`
	AddressLine2 string    `json:"address_line2,omitempty" bson:"address_line2,omitempty"` // optional
	City         string    `json:"city" bson:"city" validate:"required"`
	State        string    `json:"state" bson:"state" validate:"required"`
	ZIP          string    `json:"zip" bson:"zip" validate:"required"`
	Country      string    `json:"country" bson:"country" validate:"required"`
	CreatedAt    time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" bson:"updated_at"`
}
