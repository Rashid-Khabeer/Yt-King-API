package services

import (
	"backend/helpers"
	"backend/models"
	"database/sql"
)

type Campaigns interface {
	ReadAllCampaigns() ([]*models.Campaign, error)
}

func NewCampaign() Campaigns {
	return &campaignService{db: helpers.GetDB()}
}

func (n *campaignService) ReadAllCampaigns() ([]*models.Campaign, error) {
	result, err := n.db.Query("select * from campaigns")
	if err != nil {
		return make([]*models.Campaign, 0), err
	}
	defer result.Close()
	var campaigns []*models.Campaign
	for result.Next() {
		row := models.Campaign{}
		result.Scan(&row.Id, &row.UserId, &row.CampaignType, &row.ChannelUrl, &row.IsStateBusy, &row.IsCompleted, &row.PlayerId, &row.RequiredCount, &row.RequiredTime, &row.VideoUrl, &row.CreatedAt, &row.UpdatedAt)
		campaigns = append(campaigns, &row)
	}
	return campaigns, nil
}

type campaignService struct {
	db *sql.DB
}
