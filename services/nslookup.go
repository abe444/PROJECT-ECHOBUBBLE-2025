package services

import (
	"fmt"
	"net"
	"strings"
)

type NSLookupResult struct {
	A     []string `json:"a_records"`
	AAAA  []string `json:"aaaa_records"`
	MX    []string `json:"mx_records"`
	TXT   []string `json:"txt_records"`
	NS    []string `json:"ns_records"`
	CNAME []string `json:"cname_records"`
	PTR   []string `json:"ptr_records"`
}

func NSLookup(domain string) (map[string]interface{}, error) {
	result := NSLookupResult{}

	if records, err := net.LookupIP(domain); err == nil {
		for _, ip := range records {
			if ipv4 := ip.To4(); ipv4 != nil {
				result.A = append(result.A, ipv4.String())
			}
		}
	}

	if records, err := net.LookupIP(domain); err == nil {
		for _, ip := range records {
			if ipv4 := ip.To4(); ipv4 == nil {
				result.AAAA = append(result.AAAA, ip.String())
			}
		}
	}

	if records, err := net.LookupMX(domain); err == nil {
		for _, mx := range records {
			result.MX = append(result.MX, fmt.Sprintf("%s (priority: %d)", mx.Host, mx.Pref))
		}
	}

	if records, err := net.LookupTXT(domain); err == nil {
		result.TXT = records
	}

	if records, err := net.LookupNS(domain); err == nil {
		for _, ns := range records {
			result.NS = append(result.NS, ns.Host)
		}
	}

	if cname, err := net.LookupCNAME(domain); err == nil {
		result.CNAME = append(result.CNAME, strings.TrimSuffix(cname, "."))
	}

	if addrs, err := net.LookupAddr(domain); err == nil {
		result.PTR = addrs
	}

	return map[string]interface{}{"nslookup": result}, nil
}
