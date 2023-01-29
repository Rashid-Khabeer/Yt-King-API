package models

import (
	base2 "backend/models/base"
)

type FeedbakType int

const (
	Bug FeedbakType = iota
	Feature
	Suggestion
)

type Feedback struct {
	Id     *int         `bson:"_id" json:"id,omitempty"`
	UserId *int         `bson:"user_id" json:"user_id"`
	Type   *FeedbakType `bson:"type" json:"type"`
	Text   *string      `bson:"text" json:"text"`
	base2.Timestamped
}
