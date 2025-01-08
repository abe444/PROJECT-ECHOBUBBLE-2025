package controllers

import (
	"fmt"
	"net/http"

	"github.com/PuerkitoBio/goquery"
	"github.com/gin-gonic/gin"
)

func ScrapeHandler(c *gin.Context) {
	url := c.Query("url")
	timestamp := c.Query("timestamp")

	if url == "" || timestamp == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Both 'url' and 'timestamp' query parameters are required",
		})
	}

	waybackURL := fmt.Sprintf("https://web.archive.org/web/%s/%s", timestamp, url)

	resp, err := http.Get(waybackURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch data from Wayback Machine",
		})
		return
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to parse the Wayback Machine page.",
		})
		return
	}

	title := doc.Find("title").Text()

	c.JSON(http.StatusOK, gin.H{
		"requested_url": url,
		"timestamp":     timestamp,
		"archive_url":   waybackURL,
		"title":         title,
		"message":       "Scraping successful",
	})
}
