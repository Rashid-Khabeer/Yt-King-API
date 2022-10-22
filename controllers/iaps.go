package controllers

import (
	"backend/services"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
)

func NewTransactions(group *gin.RouterGroup) {
	c := iapController{service: services.NewTransactions()}
	group.GET("", c.get)
	group.POST("/revenuecat-purchase-update", c.revenueCatPurchaseUpdate)
}

type iapController struct {
	service services.Iap
}

func (u *iapController) get(c *gin.Context) {
	tr, err := u.service.ReadAllTransactions()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, tr)
}

func (u *iapController) revenueCatPurchaseUpdate(c *gin.Context) {
	all, _ := ioutil.ReadAll(c.Request.Body)
	u.service.HandleRevenueCatWebhooks(all)
	// c.JSON(http.StatusOK, "")
}
