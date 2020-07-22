package shodan

import (
	"encoding/json"
	"fmt"
	"net/http"
)

//HostLocation struct contains host location info
type HostLocation struct {
	City         string  `json:"city"`
	RegionCode   string  `json:"region_code"`
	AreaCode     int     `json:"area_code"`
	Longitude    float32 `json:"longitude"`
	CountryCode3 string  `json:"country_code3"`
	CountryName  string  `json:"country_name"`
	PostalCode   string  `json:"postal_code"`
	DMACode      int     `json:"dma_code"`
	CountryCode  string  `json:"country_code"`
	Latitude     float32 `json:"latitude"`
}

//Host struct contains host info
type Host struct {
	OS        string       `json:"os"`
	Timestamp string       `json:"timestamp"`
	ISP       string       `json:"isp"`
	ASN       string       `json:"asn"`
	Hostname  []string     `json:"hostnames"`
	IP        int64        `json:"ip"`
	Domains   []string     `json:"domains"`
	Org       string       `json:"org"`
	Data      string       `json:"data"`
	Port      int          `json:"port"`
	IPString  string       `json:"ip_str"`
	Location  HostLocation `json:"location"`
}

//HostSearch struct is the top level matches response
type HostSearch struct {
	Matches []Host `json:"matches"`
}

//HostSearch search Shodan for host info and use facets
//to get summary information for different properties.
func (s *Client) HostSearch(q string) (*HostSearch, error) {
	res, err := http.Get(fmt.Sprintf("%s/shodan/host/search?key=%s&query=%s", BaseURL, s.apiKey, q))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var hostinfo HostSearch
	if err = json.NewDecoder(res.Body).Decode(&hostinfo); err != nil {
		return nil, err
	}

	return &hostinfo, nil
}
