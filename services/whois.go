package services

import (
	"github.com/likexian/whois"
	whois_parser_go "github.com/likexian/whois-parser"
)

func Whois(url string) (map[string]interface{}, error) {
	whoisData, err := whois.Whois(url)
	if err != nil {
		return nil, err
	}

	parsedData, err := whois_parser_go.Parse(whoisData)
	if err != nil {
		return map[string]interface{}{
			"whois": map[string]interface{}{
				"parse_failed": true,
				"raw_whois":    whoisData,
				"domain_name":  url,
			},
		}, nil
	}

	result := map[string]interface{}{
		"whois": map[string]interface{}{
			"domain_name":   parsedData.Domain.Domain,
			"puny_code":     parsedData.Domain.Punycode,
			"whois_server":  parsedData.Domain.WhoisServer,
			"updated_date":  parsedData.Domain.UpdatedDate,
			"creation_date": parsedData.Domain.CreatedDate,
			"expiry_date":   parsedData.Domain.ExpirationDate,
			"dnssec":        parsedData.Domain.DNSSec,
			"registrar":     parsedData.Registrar.Name,
			"id":            parsedData.Domain.ID,
		},
	}
	return result, nil
}
