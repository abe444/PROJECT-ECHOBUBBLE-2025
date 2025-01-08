package routes

import (
	"web-scraper/controllers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Wayback scraper API is running!",
		})
	})

	router.GET("/api/scrape", controllers.ScrapeHandler)
}
