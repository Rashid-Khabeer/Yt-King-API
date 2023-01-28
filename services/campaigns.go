package services

import (
	"backend/helpers"
	"backend/models"
	"database/sql"
	"time"
)

type Campaigns interface {
	CreateCampaign(campaign *models.Campaign) (int, error)
	UpdateCampaign(campaign *models.Campaign) (*models.Campaign, error)
	DeleteCampaign(id int) (bool, error)
	ReadAllCampaigns() ([]*models.Campaign, error)
	FetchCampaigns(campaignType string, user int, count int) ([]*models.Campaign, error)
	FetchOwnActionCampaigns(campaignType string, user int) ([]*models.Campaign, error)
	FetchOwnCampaigns(user int) ([]*models.Campaign, error)
	FetchOwnCampaignsCount(user int) (int, error)
}

func NewCampaign() Campaigns {
	return &campaignService{db: helpers.GetDB()}
}

func (n *campaignService) CreateCampaign(campaign *models.Campaign) (int, error) {
	sql := "INSERT INTO campaigns(user_id, campaign_type, channel_url, is_state_busy,is_completed, player_id, required_count, required_time, video_url, created_at, updated_at) VALUES(?,?,?,?,?,?,?,?,?,?,?)"
	insert, err := n.db.Prepare(sql)
	if err != nil {
		return 0, err
	}
	campaign.CreatedAt = time.Now()
	campaign.UpdatedAt = time.Now()
	response, err := insert.Exec(campaign.UserId, campaign.CampaignType, campaign.ChannelUrl, campaign.IsStateBusy, campaign.IsCompleted, campaign.PlayerId, campaign.RequiredCount, campaign.RequiredTime, campaign.VideoUrl, campaign.CreatedAt, campaign.UpdatedAt)
	if err != nil {
		return 0, err
	}
	defer insert.Close()
	id, err := response.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

func (n *campaignService) UpdateCampaign(campaign *models.Campaign) (*models.Campaign, error) {
	sql := "UPDATE campaigns SET is_completed = ?, campaign_type = ?, required_count=?, required_time=?, video_url=?, updated_at=? WHERE id = ?"
	insert, err := n.db.Prepare(sql)
	if err != nil {
		return nil, err
	}
	campaign.UpdatedAt = time.Now()
	response, err := insert.Exec(campaign.IsCompleted, campaign.CampaignType, campaign.RequiredCount, campaign.RequiredTime, campaign.VideoUrl, campaign.UpdatedAt, campaign.Id)
	if err != nil {
		return nil, err
	}
	no, err := response.RowsAffected()
	if err != nil {
		return nil, err
	}
	if no < 1 {
		return nil, err
	}
	defer insert.Close()
	return campaign, nil
}

func (n *campaignService) DeleteCampaign(id int) (bool, error) {
	// sql := "delete from participants where campaign_id = ?"
	// deleteParticipants, err := n.db.Prepare(sql)
	// if err != nil {
	// 	return false, err
	// }
	// response, err := deleteParticipants.Exec(id)
	// if err != nil {
	// 	return false, err
	// }
	// response.RowsAffected()
	// defer deleteParticipants.Close()
	sql1 := "update campaigns SET is_deleted = true where id = ?"
	deleteCamp, err := n.db.Prepare(sql1)
	response1, err := deleteCamp.Exec(id)
	if err != nil {
		return false, err
	}
	no, err := response1.RowsAffected()
	if err != nil {
		return false, err
	}
	if no < 1 {
		return false, err
	}
	defer deleteCamp.Close()
	return true, nil
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
		result.Scan(&row.Id, &row.UserId, &row.CampaignType, &row.ChannelUrl, &row.IsStateBusy, &row.IsCompleted, &row.IsDeleted, &row.PlayerId, &row.RequiredCount, &row.RequiredTime, &row.VideoUrl, &row.CreatedAt, &row.UpdatedAt)
		campaigns = append(campaigns, &row)
	}
	if campaigns == nil {
		return make([]*models.Campaign, 0), nil
	}
	return campaigns, nil
}

func (n *campaignService) FetchOwnActionCampaigns(campaignType string, user int) ([]*models.Campaign, error) {
	query := "SELECT DISTINCT campaign_id FROM `participants` WHERE user_id = ?"
	result, err := n.db.Query(query, user)
	if err != nil {
		return make([]*models.Campaign, 0), err
	}
	defer result.Close()
	var campaigns []*models.Campaign
	for result.Next() {
		var campaignId int
		result.Scan(&campaignId)
		query1 := "select * from campaigns where id = ? AND campaign_type = ?"
		result1, err := n.db.Query(query1, campaignId, campaignType)
		if err != nil {
			return make([]*models.Campaign, 0), err
		}
		defer result1.Close()
		if result1.Next() {
			row := models.Campaign{}
			result1.Scan(&row.Id, &row.UserId, &row.CampaignType, &row.ChannelUrl, &row.IsStateBusy, &row.IsCompleted, &row.IsDeleted, &row.PlayerId, &row.RequiredCount, &row.RequiredTime, &row.VideoUrl, &row.CreatedAt, &row.UpdatedAt)
			campaigns = append(campaigns, &row)
		}
	}
	if campaigns == nil {
		return make([]*models.Campaign, 0), nil
	}
	return campaigns, nil
}

func (n *campaignService) FetchCampaigns(campaignType string, user int, count int) ([]*models.Campaign, error) {
	query := "select * from campaigns where is_deleted = false AND is_completed = false AND user_id != ? AND campaign_type = ? ORDER BY created_at DESC limit ?"
	result, err := n.db.Query(query, user, campaignType, count)
	if err != nil {
		return make([]*models.Campaign, 0), err
	}
	defer result.Close()
	var campaigns []*models.Campaign
	for result.Next() {
		row := models.Campaign{}
		result.Scan(&row.Id, &row.UserId, &row.CampaignType, &row.ChannelUrl, &row.IsStateBusy, &row.IsCompleted, &row.IsDeleted, &row.PlayerId, &row.RequiredCount, &row.RequiredTime, &row.VideoUrl, &row.CreatedAt, &row.UpdatedAt)
		par, _ := n.db.Query("select * from participants where campaign_id = ? AND user_id = ?", row.Id, user)
		defer par.Close()
		if par.Next() {
			continue
		}
		campaigns = append(campaigns, &row)
	}
	var prCampaigns []*models.Campaign
	var npCampaigns []*models.Campaign
	for _, camp := range campaigns {
		uResult, _ := n.db.Query("select * from users where id =?", camp.UserId)
		defer uResult.Close()
		for uResult.Next() {
			row := models.User{}
			uResult.Scan(&row.Id, &row.Name, &row.Email, &row.Image, &row.TotalCoins, &row.PremiumType, &row.HasPremium, &row.LastDate, &row.Password, &row.RememberToken, &row.CreatedAt, &row.UpdatedAt, &row.AppVersion, &row.IsBlocked, &row.BlockedDays)
			if *row.HasPremium == true {
				prCampaigns = append(prCampaigns, camp)
			} else {
				npCampaigns = append(npCampaigns, camp)
			}
		}
	}
	if prCampaigns == nil && npCampaigns == nil {
		return make([]*models.Campaign, 0), nil
	}
	return append(prCampaigns, npCampaigns...), nil
}

func (n *campaignService) FetchOwnCampaigns(user int) ([]*models.Campaign, error) {
	query := "select * from campaigns where is_deleted = false AND user_id = ? ORDER BY created_at DESC"
	result, err := n.db.Query(query, user)
	if err != nil {
		return make([]*models.Campaign, 0), err
	}
	defer result.Close()
	var campaigns []*models.Campaign
	for result.Next() {
		row := models.Campaign{}
		result.Scan(&row.Id, &row.UserId, &row.CampaignType, &row.ChannelUrl, &row.IsStateBusy, &row.IsCompleted, &row.IsCompleted, &row.PlayerId, &row.RequiredCount, &row.RequiredTime, &row.VideoUrl, &row.CreatedAt, &row.UpdatedAt)
		var participants []*models.PopulatedParticipants
		parResult, err := n.db.Query("select * from participants where campaign_id = ?", row.Id)
		if err != nil {
			return make([]*models.Campaign, 0), err
		}
		defer parResult.Close()
		for parResult.Next() {
			row1 := models.Participants{}
			participant := models.PopulatedParticipants{}
			parResult.Scan(&row1.Id, &row1.CampaignId, &row1.UserId, &row1.CreatedAt, &row1.UpdatedAt)
			userResult, err := n.db.Query("select * from users where id =?", row1.UserId)
			if err != nil {
				return make([]*models.Campaign, 0), err
			}
			if userResult.Next() {
				userRow := models.User{}
				userResult.Scan(&userRow.Id, &userRow.Name, &userRow.Email, &userRow.Image, &userRow.TotalCoins, &userRow.PremiumType, &userRow.HasPremium, &userRow.LastDate, &userRow.Password, &userRow.RememberToken, &userRow.CreatedAt, &userRow.UpdatedAt, &userRow.AppVersion, &userRow.IsBlocked, &userRow.BlockedDays)
				participant.Id = row1.Id
				participant.CampaignId = row1.CampaignId
				participant.UserId = &userRow
				participant.CreatedAt = row1.CreatedAt
				participant.UpdatedAt = row1.UpdatedAt
			}
			participants = append(participants, &participant)
			userResult.Close()
		}
		row.Participants = participants
		campaigns = append(campaigns, &row)
	}
	if campaigns == nil {
		return make([]*models.Campaign, 0), nil
	}
	return campaigns, nil
}

func (n *campaignService) FetchOwnCampaignsCount(user int) (int, error) {
	query := "select count(*) from campaigns where is_deleted = false AND user_id = ?"
	result, err := n.db.Query(query, user)
	if err != nil {
		return 0, err
	}
	defer result.Close()
	var count int
	if result.Next() {
		result.Scan(&count)
	}
	return count, nil
}

type campaignService struct {
	db *sql.DB
}
