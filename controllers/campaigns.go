package controllers

import (
	"backend/models"
	"backend/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func NewCampaign(group *gin.RouterGroup) {
	c := campaignController{service: services.NewCampaign()}
	group.GET("", c.get)
	group.POST("", c.post)
	group.PUT("", c.put)
	group.DELETE("/:id", c.delete)
	group.GET("/type/:type/user/:user/count/:count", c.fetchCampaigns)
	group.GET("/participated/user/:user/type/:type", c.fetchOwnActionCampaigns)
	group.GET("/user/:user", c.fetchOwnCampaigns)
	group.GET("/user/:user/count", c.fetchOwnCampaignsCount)
}

type campaignController struct {
	service services.Campaigns
}

func (u *campaignController) post(c *gin.Context) {
	campaign := &models.Campaign{}

	err := c.ShouldBindJSON(campaign)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	result, err := u.service.CreateCampaign(campaign)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, result)
}

func (u *campaignController) put(c *gin.Context) {
	campaign := &models.Campaign{}

	err := c.ShouldBindJSON(campaign)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	result, err := u.service.UpdateCampaign(campaign)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, result)
}

func (u *campaignController) delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{"error": err.Error()})
		return
	}
	response, err := u.service.DeleteCampaign(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, response)
}

func (u *campaignController) get(c *gin.Context) {
	campaigns, err := u.service.ReadAllCampaigns()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, campaigns)
}

func (u *campaignController) fetchCampaigns(c *gin.Context) {
	user, err := strconv.Atoi(c.Param("user"))
	if err != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{"error": err.Error()})
		return
	}
	count, err := strconv.Atoi(c.Param("count"))
	if err != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{"error": err.Error()})
		return
	}
	campaigns, err := u.service.FetchCampaigns(c.Param("type"), user, count)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, campaigns)
}

func (u *campaignController) fetchOwnActionCampaigns(c *gin.Context) {
	user, err := strconv.Atoi(c.Param("user"))
	if err != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{"error": err.Error()})
		return
	}
	campaigns, err := u.service.FetchOwnActionCampaigns(c.Param("type"), user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, campaigns)
}

func (u *campaignController) fetchOwnCampaigns(c *gin.Context) {
	user, err := strconv.Atoi(c.Param("user"))
	if err != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{"error": err.Error()})
		return
	}
	campaigns, err := u.service.FetchOwnCampaigns(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, campaigns)
}

func (u *campaignController) fetchOwnCampaignsCount(c *gin.Context) {
	user, err := strconv.Atoi(c.Param("user"))
	if err != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{"error": err.Error()})
		return
	}
	count, err := u.service.FetchOwnCampaignsCount(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, count)
}
