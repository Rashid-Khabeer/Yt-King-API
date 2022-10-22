package controllers

import (
	"backend/models"
	"backend/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

func NewPartripants(group *gin.RouterGroup) {
	c := participantsController{service: services.NewPartripants()}
	group.POST("", c.post)
}

func (u *participantsController) post(c *gin.Context) {
	participant := &models.Participants{}

	err := c.ShouldBindJSON(participant)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	result, err := u.service.AddParticipant(participant)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, result)
}

type participantsController struct {
	service services.Participants
}
