package models

import (
	base2 "backend/models/base"
)

type Transactions struct {
	Id              *int                    `bson:"_id" json:"id,omitempty"`
	Email           *string                 `bson:"email" json:"email,omitempty"`
	Date            *string                 `bson:"date" json:"date,omitempty"`
	Price           *float32                `bson:"price" json:"price,omitempty"`
	TransactionId   *string                 `bson:"transaction_id" json:"transaction_id,omitempty"`
	WebhookResponse *map[string]interface{} `bson:"webhook_response" json:"webhook_response,omitempty"`
	Type            *string                 `bson:"type" json:"type,omitempty"`
	ProductId       *string                 `bson:"product_id" json:"product_id,omitempty"`
	base2.Timestamped
}
