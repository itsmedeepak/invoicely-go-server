package models

import (
	"time"
)

type Customer struct {
	CustomerId    string    `json:"_id" bson:"_id"`
	UserId        string    `json:"user_id" bson:"user_id"`
	Email         string    `json:"email" bson:"email"`
	Phone         string    `json:"phone" bson:"phone"`
	StreetAddress string    `json:"street_address" bson:"street_address"`
	City          string    `json:"city" bson:"city"`
	State         string    `json:"state" bson:"state"`
	Country       string    `json:"country" bson:"country"`
	FirstName     string    `json:"first_name" bson:"first_name"`
	LastName      string    `json:"last_name" bson:"last_name"`
	CreatedAt     time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" bson:"updated_at"`
}
