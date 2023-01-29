package services

import (
	"backend/helpers"
	"backend/models"
	"database/sql"
	"time"
)

type Feedback interface {
	CreateFeedback(feedback *models.Feedback) (int, error)
	ReadAllFeedback() ([]*models.Feedback, error)
}

func NewFeedback() Feedback {
	return &feedbakService{db: helpers.GetDB()}
}

func (n *feedbakService) CreateFeedback(feedbak *models.Feedback) (int, error) {
	sql := "INSERT INTO feedbacks(user_id, type, text, createdAt, updatedAt) VALUES(?,?,?,?,?)"
	insert, err := n.db.Prepare(sql)
	if err != nil {
		return 0, err
	}
	feedbak.CreatedAt = time.Now()
	feedbak.UpdatedAt = time.Now()
	response, err := insert.Exec(feedbak.UserId, feedbak.Type, feedbak.Text, feedbak.CreatedAt, feedbak.UpdatedAt)
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

func (n *feedbakService) ReadAllFeedback() ([]*models.Feedback, error) {
	result, err := n.db.Query("select * from feedbacks")
	if err != nil {
		return make([]*models.Feedback, 0), err
	}
	defer result.Close()
	var feedbacks []*models.Feedback
	for result.Next() {
		row := models.Feedback{}
		result.Scan(&row.Id, &row.UserId, &row.Type, &row.Text, &row.CreatedAt, &row.UpdatedAt)
		feedbacks = append(feedbacks, &row)
	}
	if feedbacks == nil {
		return make([]*models.Feedback, 0), nil
	}
	return feedbacks, nil

}

type feedbakService struct {
	db *sql.DB
}
