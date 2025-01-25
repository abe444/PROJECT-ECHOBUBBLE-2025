package routes

import (
	"echobubble/controller"
	"echobubble/handler"

	"github.com/gin-gonic/gin"
)

func InitRoutes(routes *gin.Engine) {

	routes.GET("/status", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "OK"})
	})

	routes.LoadHTMLGlob("templates/*")
	routes.Static("/static", "./static")
	routes.StaticFile("/css", "./static/css/styles.css")

	routes.GET("/", handler.GreeterPage)

	routes.GET("/scan", handler.GreeterPage)
	routes.GET("/scan/:domain", controller.Scanner)

	routes.POST("/scan/:domain", controller.Scanner)
	routes.POST("/api/scan", controller.Scanner)

	routes.GET("/results", handler.ResultsPage)

}
