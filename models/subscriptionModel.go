package models

import "time"

type Subscription struct {
	UserID            string    `json:"user_id" bson:"user_id"`
	Plan              string    `json:"plan" bson:"plan"`
	ValidTill         time.Time `json:"valid_till" bson:"valid_till"`
	Status            string    `json:"status" bson:"status"`
	CreditsUsed       int       `json:"credits_used" bson:"credits_used"`
	CreditsRemaining  int       `json:"credits_remaining" bson:"credits_remaining"`
	AverageDailyUsage float64   `json:"average_daily_usage" bson:"average_daily_usage"`
	LastRefreshed     time.Time `json:"last_refreshed" bson:"last_refreshed"`
	CreatedAt         time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt         time.Time `json:"updated_at" bson:"updated_at"`
}
