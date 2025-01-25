package services

import (
	"bufio"
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

func FindSubs(domain string) (map[string]interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	if _, err := exec.LookPath("subfinder"); err != nil {
		return nil, fmt.Errorf("subfinder is not installed: %v", err)
	}

	cmd := exec.CommandContext(ctx, "subfinder", "-d", domain, "-silent", "-cs")

	output, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to create output pipe: %v", err)
	}

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("failed to start subfinder: %v", err)
	}

	var subdomains []string
	scanner := bufio.NewScanner(output)
	for scanner.Scan() {
		subdomain := strings.TrimSpace(scanner.Text())
		if subdomain != "" {
			subdomains = append(subdomains, subdomain)
		}
	}

	if err := cmd.Wait(); err != nil {
		return nil, fmt.Errorf("subfinder command failed: %v", err)
	}

	return map[string]interface{}{
		"subdomains": subdomains,
	}, nil
}
