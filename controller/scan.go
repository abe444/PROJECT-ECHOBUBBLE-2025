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

type JobStatus struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
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
	urlScanChan := make(chan ScanResult)

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

	go func() {
		result, err := services.ScanURL(url)
		urlScanChan <- ScanResult{result, err}
	}()

	results := make(map[string]interface{})
	var errors []string
	timeout := time.After(60 * time.Second)

	jobStatus := map[string]JobStatus{
		"whois":      {Status: "Pending"},
		"subdomains": {Status: "Pending"},
		"nslookup":   {Status: "Pending"},
		"serverip":   {Status: "Pending"},
		"urlscan":    {Status: "Pending"},
	}

	for i := 0; i < 5; i++ {
		select {
		case result := <-whoisChan:
			if result.Error != nil {
				errors = append(errors, fmt.Sprintf("WHOIS error: %v", result.Error))
				jobStatus["whois"] = JobStatus{Status: "Failed", Error: result.Error.Error()}
			} else {
				results["whois"] = result.Data["whois"]
				jobStatus["whois"] = JobStatus{Status: "Completed"}
			}

		case result := <-subdomainChan:
			if result.Error != nil {
				errors = append(errors, fmt.Sprintf("Subdomain error: %v", result.Error))
				jobStatus["subdomains"] = JobStatus{Status: "Failed", Error: result.Error.Error()}
			} else {
				results["subdomains"] = result.Data["subdomains"]
				jobStatus["subdomains"] = JobStatus{Status: "Completed"}
			}

		case result := <-nslookupChan:
			if result.Error != nil {
				errors = append(errors, fmt.Sprintf("DNS lookup error: %v", result.Error))
				jobStatus["nslookup"] = JobStatus{Status: "Failed", Error: result.Error.Error()}
			} else {
				results["nslookup"] = result.Data["nslookup"]
				jobStatus["nslookup"] = JobStatus{Status: "Completed"}
			}

		case result := <-serverIPChan:
			if result.Error != nil {
				errors = append(errors, fmt.Sprintf("Server IP error: %v", result.Error))
				jobStatus["serverip"] = JobStatus{Status: "Failed", Error: result.Error.Error()}
			} else {
				results["serverip"] = result.Data["serverip"]
				jobStatus["serverip"] = JobStatus{Status: "Completed"}
			}

		case result := <-urlScanChan:
			if result.Error != nil {
				errors = append(errors, fmt.Sprintf("URL scan error: %v", result.Error))
				jobStatus["urlscan"] = JobStatus{Status: "Failed", Error: result.Error.Error()}
			} else {
				results["urlscan"] = result.Data
				jobStatus["urlscan"] = JobStatus{Status: "Completed"}
			}

		case <-timeout:
			errors = append(errors, "Some scans timed out")
			goto RenderResults
		}
	}

RenderResults:
	c.HTML(http.StatusOK, "index.html", gin.H{
		"result":    results,
		"url":       url,
		"error":     strings.Join(errors, "; "),
		"jobStatus": jobStatus,
	})
}
