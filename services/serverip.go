package services

import (
	"echobubble/models"
	"fmt"
	"image/color"
	"image/png"
	"net"
	"os"

	sm "github.com/flopp/go-staticmaps"
	"github.com/golang/geo/s2"
	"github.com/ipinfo/go/v2/ipinfo"
)

func GetServerIP(domain string) (map[string]interface{}, error) {
	ips, err := net.LookupIP(domain)
	if err != nil {
		return nil, err
	}

	client := ipinfo.NewClient(nil, nil, "") // token if needed
	var ipInfo []models.IPResponse

	for _, ip := range ips {
		if ipv4 := ip.To4(); ipv4 != nil {
			info, err := client.GetIPInfo(ip)
			if err != nil {
				continue
			}

			var lat, lon float64
			if info.Location != "" {
				_, err := fmt.Sscanf(info.Location, "%f,%f", &lat, &lon)
				if err != nil {
					lat, lon = 0, 0
				}
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
				Latitude:     lat,
				Longitude:    lon,
			}

			if lat != 0 && lon != 0 {
				mapPath, err := generateMapImage(lat, lon)
				if err != nil {
					fmt.Printf("Failed to generate map image: %v\n", err)
				} else {
					response.MapImage = mapPath
				}
			}

			ipInfo = append(ipInfo, response)
		}
	}

	return map[string]interface{}{"serverip": ipInfo}, nil
}

func generateMapImage(lat, lon float64) (string, error) {
	ctx := sm.NewContext()
	ctx.SetSize(400, 300)
	ctx.SetZoom(5)
	ctx.SetCenter(s2.LatLngFromDegrees(lat, lon))

	marker := sm.NewMarker(
		s2.LatLngFromDegrees(lat, lon),
		color.RGBA{R: 255, G: 0, B: 0, A: 255},
		16.0,
	)
	ctx.AddMarker(marker)

	err := os.MkdirAll("static/maps", 0755)
	if err != nil {
		return "", err
	}

	mapPath := fmt.Sprintf("static/maps/map_%.4f_%.4f.png", lat, lon)

	img, err := ctx.Render()
	if err != nil {
		return "", err
	}

	file, err := os.Create(mapPath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	if err := png.Encode(file, img); err != nil {
		return "", err
	}

	return mapPath, nil
}
