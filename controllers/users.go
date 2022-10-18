package controllers

import (
	"backend/models"
	"backend/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

func NewUsers(group *gin.RouterGroup) {
	c := usersController{service: services.NewUsers()}
	group.GET("", c.get)
	group.POST("", c.post)

	// group.GET("/:id", c.getOne)

	// group.PUT("", c.put)
	group.PATCH("", c.patch)
	// group.DELETE("/:id", c.delete)
}

type usersController struct {
	service services.Users
}

func (u *usersController) get(c *gin.Context) {
	user, err := u.service.ReadAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (u *usersController) post(c *gin.Context) {
	user := &models.User{}

	err := c.ShouldBindJSON(user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	result, err := u.service.CreateUser(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, result)
}

func (u *usersController) patch(c *gin.Context) {
	user := &models.User{}

	err := c.ShouldBindJSON(user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	result, err := u.service.UpdateUser(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, result)
}
