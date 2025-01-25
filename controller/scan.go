package controller

import (
	"echobubble/services"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type ScanError struct {
	Service string
	Err     error
}

func (se ScanError) Error() string {
	return fmt.Sprintf("Error in %s scan: %v", se.Service, se.Err)
}

type ScanResult struct {
	Data  map[string]interface{}
	Error error
}

func Scanner(c *gin.Context) {
	url := c.Param("domain")
	if url == "" {
		c.HTML(http.StatusBadRequest, "index.html", gin.H{
			"error": "URL is required",
		})
		return
	}

	whoisChan := make(chan ScanResult)
	subdomainChan := make(chan ScanResult)
	nslookupChan := make(chan ScanResult)
	serverIPChan := make(chan ScanResult)

	go func() {
		result, err := services.Whois(url)
		whoisChan <- ScanResult{result, err}
	}()

	go func() {
		result, err := services.FindSubs(url)
		subdomainChan <- ScanResult{result, err}
	}()

	go func() {
		result, err := services.GetServerIP(url)
		serverIPChan <- ScanResult{result, err}
	}()

	go func() {
		result, err := services.NSLookup(url)
		nslookupChan <- ScanResult{result, err}
	}()

	results := make(map[string]interface{})
	var errors []string
	timeout := time.After(60 * time.Second)

	for i := 0; i < 4; i++ {
		select {
		case result := <-whoisChan:
			if result.Error != nil {
				errors = append(errors, fmt.Sprintf("WHOIS error: %v", result.Error))
			} else {
				results["whois"] = result.Data["whois"]
			}

		case result := <-subdomainChan:
			if result.Error != nil {
				errors = append(errors, fmt.Sprintf("Subdomain error: %v", result.Error))
			} else {
				results["subdomains"] = result.Data["subdomains"]
			}

		case result := <-nslookupChan:
			if result.Error != nil {
				errors = append(errors, fmt.Sprintf("DNS lookup error: %v", result.Error))
			} else {
				results["nslookup"] = result.Data["nslookup"]
			}

		case result := <-serverIPChan:
			if result.Error != nil {
				errors = append(errors, fmt.Sprintf("Server IP error: %v", result.Error))
			} else {
				results["serverip"] = result.Data["serverip"]
			}

		case <-timeout:
			errors = append(errors, "Some scans timed out")
			goto RenderResults
		}
	}

RenderResults:
	c.HTML(http.StatusOK, "index.html", gin.H{
		"result": results,
		"url":    url,
		"error":  strings.Join(errors, "; "),
	})
}
