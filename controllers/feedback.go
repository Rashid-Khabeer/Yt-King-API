package controllers

import (
	"backend/models"
	"backend/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

func NewFeedback(group *gin.RouterGroup) {
	c := feedbackController{service: services.NewFeedback()}
	group.GET("", c.get)
	group.POST("", c.post)
}

type feedbackController struct {
	service services.Feedback
}

func (u *feedbackController) post(c *gin.Context) {
	feedback := &models.Feedback{}

	err := c.ShouldBindJSON(feedback)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	result, err := u.service.CreateFeedback(feedback)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, result)
}

func (u *feedbackController) get(c *gin.Context) {
	feedbacks, err := u.service.ReadAllFeedback()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, feedbacks)
}
