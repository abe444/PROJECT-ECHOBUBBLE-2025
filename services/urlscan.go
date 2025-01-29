package services

import (
	"bufio"
	"fmt"
	"net/http"
	"strings"
)

func CheckOpenPhish(domain string) (bool, error) {
	resp, err := http.Get("https://openphish.com/feed.txt")
	if err != nil {
		return false, fmt.Errorf("failed to fetch OpenPhish feed: %v", err)
	}
	defer resp.Body.Close()

	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), domain) {
			return true, nil
		}
	}

	return false, nil
}

func CheckURLhaus(domain string) (bool, error) {
	resp, err := http.Get("https://urlhaus.abuse.ch/downloads/text_online/")
	if err != nil {
		return false, fmt.Errorf("failed to fetch URLhaus feed: %v", err)
	}
	defer resp.Body.Close()

	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), domain) {
			return true, nil
		}
	}

	return false, nil
}

func ScanURL(domain string) (map[string]interface{}, error) {
	openPhishResult, err := CheckOpenPhish(domain)
	if err != nil {
		return nil, fmt.Errorf("OpenPhish scan failed: %v", err)
	}

	urlhausResult, err := CheckURLhaus(domain)
	if err != nil {
		return nil, fmt.Errorf("URLhaus scan failed: %v", err)
	}

	result := map[string]interface{}{
		"openphish": openPhishResult,
		"urlhaus":   urlhausResult,
	}

	return result, nil
}
