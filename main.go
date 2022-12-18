package main

import (
	"time"

	"os"

	"backend/controllers"
	"backend/helpers"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	helpers.InitiateMySql()
	defer helpers.GetDB().Close()
	router := gin.Default()
	router.StaticFile("/app-ads.txt", "./assets/app-ads.txt")
	router.StaticFile("/privacy-policy", "./assets/privacy-policy.html")
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"*"},
		AllowHeaders:     []string{"*"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	router.GET("/", func(c *gin.Context) {
		c.String(200, "Welcome to Parenthises")
	})

	controllers.NewUsers(router.Group("/users"))
	controllers.NewCampaign(router.Group("/campaigns"))
	controllers.NewPartripants(router.Group("/participants"))
	controllers.NewTransactions(router.Group("/iap"))

	port, flag := os.LookupEnv("PORT")
	if flag {
		err := router.Run(":" + port)
		if err != nil {
			panic(err)
		}
	}

	err := router.Run(":5000")
	if err != nil {
		panic(err)
	}
}
