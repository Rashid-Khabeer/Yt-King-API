package models

import (
	base2 "backend/models/base"
)

type User struct {
	Id            *int    `bson:"_id" json:"id,omitempty"`
	Name          *string `bson:"name" json:"name"`
	Email         *string `bson:"email" json:"email"`
	Image         *string `bson:"image" json:"image"`
	TotalCoins    *int    `bson:"total_coins" json:"total_coins,omitempty"`
	PremiumType   *string `bson:"premium_type" json:"premium_type,omitempty"`
	HasPremium    *bool   `bson:"has_premium" json:"has_premium,omitempty"`
	LastDate      *string `bson:"last_date" json:"last_date,omitempty"`
	RememberToken *string `bson:"remember_token" json:"remember_token,omitempty"`
	Password      *string `bson:"password" json:"password,omitempty"`
	base2.Timestamped
}
