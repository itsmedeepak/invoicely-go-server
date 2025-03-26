package models

import "time"

type InvoiceConfiguration struct {
	UserID    string    `json:"user_id" bson:"user_id"`
	Name      string    `json:"name,omitempty" bson:"name,omitempty"`
	Address   string    `json:"address,omitempty" bson:"address,omitempty"`
	City      string    `json:"city,omitempty" bson:"city,omitempty"`
	Country   string    `json:"country,omitempty" bson:"country,omitempty"`
	Phone1    string    `json:"phone1,omitempty" bson:"phone1,omitempty"`
	Phone2    string    `json:"phone2,omitempty" bson:"phone2,omitempty"`
	Email     string    `json:"email,omitempty" bson:"email,omitempty"`
	LogoURL   string    `json:"logo_url,omitempty" bson:"logo_url,omitempty"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}
