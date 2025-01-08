package main

import (
	"log"
	"os"

	"web-scraper/routes"

	"github.com/gin-gonic/gin"
)

func main() {

	ginMode := os.Getenv("GIN_MODE")
	if ginMode == "" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(ginMode)
	}

	router := gin.Default()

	routes.SetupRoutes(router)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server is running on http://localhost:%s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

}
