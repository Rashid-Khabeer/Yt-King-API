package models

import (
	base2 "backend/models/base"
)

type Campaign struct {
	Id            *int    `bson:"_id" json:"id,omitempty"`
	UserId        *int    `bson:"user_id" json:"user_id"`
	CampaignType  *string `bson:"campaign_type" json:"campaign_type"`
	ChannelUrl    *string `bson:"channel_url" json:"channel_url,omitempty"`
	IsStateBusy   *bool   `bson:"is_state_busy" json:"is_state_busy,omitempty"`
	IsCompleted   *bool   `bson:"is_completed" json:"is_completed,omitempty"`
	PlayerId      *string `bson:"player_id" json:"player_id,omitempty"`
	RequiredCount *int    `bson:"required_count" json:"required_count,omitempty"`
	RequiredTime  *int    `bson:"required_time" json:"required_time,omitempty"`
	VideoUrl      *string `bson:"vidoe_url" json:"vidoe_url,omitempty"`
	base2.Timestamped
}
