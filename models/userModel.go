package models

import (
	"time"
)

type User struct {
	UserId       string    `json:"_id" bson:"_id,omitempty"`
	FirstName    string    `json:"first_name" bson:"first_name" validate:"required,min=2,max=20"`
	LastName     string    `json:"last_name" bson:"last_name" validate:"required,min=2,max=20"`
	Email        string    `json:"email" bson:"email" validate:"required,email"`
	Phone        string    `json:"phone" bson:"phone" validate:"required"`
	Password     string    `json:"password" bson:"password" validate:"required,min=8"`
	AuthToken    string    `json:"auth_token" bson:"auth_token"`
	RefreshToken string    `json:"refresh_token" bson:"refresh_token"`
	CreatedAt    time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" bson:"updated_at"`
}
