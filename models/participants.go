package models

import (
	base2 "backend/models/base"
)

type Participants struct {
	Id         *int `bson:"_id" json:"id,omitempty"`
	UserId     *int `bson:"user_id" json:"user_id,omitempty"`
	CampaignId *int `bson:"campaign_type" json:"campaign_id,omitempty"`
	base2.Timestamped
}

type PopulatedParticipants struct {
	Id         *int  `bson:"_id" json:"id,omitempty"`
	UserId     *User `bson:"user_id" json:"user_id,omitempty"`
	CampaignId *int  `bson:"campaign_type" json:"campaign_id,omitempty"`
	base2.Timestamped
}
