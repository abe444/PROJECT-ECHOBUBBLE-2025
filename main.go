package main

import (
	"echobubble/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	routes.InitRoutes(router)

	if err := router.Run(":8080"); err != nil {
		panic(err)
	}
}
