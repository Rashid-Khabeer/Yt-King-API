package services

import (
	"backend/helpers"
	"backend/models"
	"database/sql"
	"strconv"
	"time"
)

type Participants interface {
	AddParticipant(participant *models.Participants) (int, error)
}

func NewPartripants() Participants {
	notificationService := NewNotificationService()
	return &participantsService{db: helpers.GetDB(), notificationService: *notificationService}
}

func (pS *participantsService) AddParticipant(participant *models.Participants) (int, error) {
	sql := "INSERT INTO participants(campaign_id, user_id, created_at, updated_at) VALUES(?,?,?,?)"
	insert, err := pS.db.Prepare(sql)
	if err != nil {
		return 0, err
	}
	participant.CreatedAt = time.Now()
	participant.UpdatedAt = time.Now()
	response, err := insert.Exec(participant.CampaignId, participant.UserId, participant.CreatedAt, participant.UpdatedAt)
	if err != nil {
		return 0, err
	}
	defer insert.Close()
	query := "SELECT COUNT(*) FROM `participants` WHERE campaign_id = ?"
	result, err := pS.db.Query(query, participant.CampaignId)
	if err != nil {
		return 0, err
	}
	defer result.Close()
	if result.Next() {
		var totalParticipants int
		result.Scan(&totalParticipants)
		query1 := "select required_count, user_id from campaigns where id =?"
		result1, err := pS.db.Query(query1, participant.CampaignId)
		if err != nil {
			return 0, err
		}
		defer result1.Close()
		if result1.Next() {
			var requiredCount int
			var campaignOwner int
			result1.Scan(&requiredCount, &campaignOwner)
			if totalParticipants >= requiredCount {
				/// udpate campaign
				updateSql := "UPDATE campaigns SET is_completed = ?, updated_at = ? WHERE id = ?"
				update, err := pS.db.Prepare(updateSql)
				if err != nil {
					return 0, err
				}
				isCompleted := true
				updatedAt := time.Now()
				response1, err := update.Exec(isCompleted, updatedAt, participant.CampaignId)
				defer update.Close()
				if err != nil {
					return 0, err
				}
				n, err := response1.RowsAffected()
				if err != nil {
					return 0, err
				}
				if n < 1 {
					return 0, err
				}
				topic := strconv.Itoa(campaignOwner)
				pS.notificationService.SendNotificationToTopic("Congratulations", "Your campaign has been completed", topic)
			}
		}
	}
	id, err := response.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

type participantsService struct {
	db                  *sql.DB
	notificationService notificationService
}
