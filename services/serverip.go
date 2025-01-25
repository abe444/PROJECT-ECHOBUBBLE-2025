package services

import (
	"echobubble/models"
	"net"

	"github.com/ipinfo/go/v2/ipinfo"
)

func GetServerIP(domain string) (map[string]interface{}, error) {
	ips, err := net.LookupIP(domain)
	if err != nil {
		return nil, err
	}

	client := ipinfo.NewClient(nil, nil, "") // token here
	var ipInfo []models.IPResponse

	for _, ip := range ips {
		if ipv4 := ip.To4(); ipv4 != nil {
			info, err := client.GetIPInfo(ip)
			if err != nil {
				continue
			}

			response := models.IPResponse{
				IP:           info.IP.String(),
				Hostname:     info.Hostname,
				City:         info.City,
				Region:       info.Region,
				Country:      info.Country,
				Location:     info.Location,
				Organization: info.Org,
				Postal:       info.Postal,
				Timezone:     info.Timezone,
			}
			ipInfo = append(ipInfo, response)
		}
	}

	return map[string]interface{}{"serverip": ipInfo}, nil
}
