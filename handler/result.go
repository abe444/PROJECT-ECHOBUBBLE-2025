package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func ResultsPage(c *gin.Context) {
	// Retrieve results from the context or session
	scanResult, exists := c.Get("scanResult")
	if !exists {
		c.HTML(http.StatusBadRequest, "index.html", gin.H{
			"error": "No results available. Please perform a scan.",
		})
		return
	}

	c.HTML(http.StatusOK, "results.html", gin.H{
		"result": scanResult,
	})
}
