package models

type IPResponse struct {
	IP           string  `json:"ip"`
	Hostname     string  `json:"hostname"`
	City         string  `json:"city"`
	Region       string  `json:"region"`
	Country      string  `json:"country"`
	Location     string  `json:"location"`
	Organization string  `json:"org"`
	Postal       string  `json:"postal"`
	Timezone     string  `json:"timezone"`
	Latitude     float64 `json:"latitude"`
	Longitude    float64 `json:"longitude"`
	MapImage     string  `json:"map_image"`
}
