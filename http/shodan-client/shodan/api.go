package shodan

import (
	"encoding/json"
	"fmt"
	"net/http"
)

//APIInfo represents an api info response struct
type APIInfo struct {
	QueryCredits int    `json:"query_credits"`
	ScanCredits  int    `json:"scan_credits"`
	Telnet       bool   `json:"telnet"`
	Plan         string `json:"plan"`
	HTTPS        bool   `json:"https"`
	Unlocked     bool   `json:"unlocked"`
}

//APIInfo Returns information about the API plan belonging to the given API key
func (s *Client) APIInfo() (*APIInfo, error) {
	res, err := http.Get(fmt.Sprintf("%s/api-info?key=%s", BaseURL, s.apiKey))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var apiInfo APIInfo
	if err = json.NewDecoder(res.Body).Decode(&apiInfo); err != nil {
		return nil, err
	}

	return &apiInfo, nil
}
