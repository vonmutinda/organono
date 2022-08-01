package providers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const ipAPIURL = "https://ipapi.co/%v/json/"

type (
	IPAPI interface {
		CountryForIP(ipAddress string) (*CountryWithIP, error)
	}

	AppIPAPI struct {
		client *http.Client
	}

	CountryWithIP struct {
		IP                 string `json:"ip"`
		City               string `json:"city"`
		CountryCode        string `json:"country_code"`
		CountryName        string `json:"country_name"`
		Currency           string `json:"currency"`
		CountryCallingCode string `json:"country_calling_code"`
	}
)

func NewIPAPI() *AppIPAPI {
	return &AppIPAPI{
		client: &http.Client{
			Timeout: time.Second * 10,
		},
	}
}

func (p *AppIPAPI) CountryForIP(ipAddress string) (*CountryWithIP, error) {

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf(ipAPIURL, ipAddress), nil)
	if err != nil {
		return &CountryWithIP{}, err
	}

	req.Header.Add("Content-Type", "application/json")

	response, err := p.client.Do(req)
	if err != nil {
		return &CountryWithIP{}, err
	}

	defer response.Body.Close()

	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return &CountryWithIP{}, err
	}

	var countryWithIP CountryWithIP

	err = json.Unmarshal(data, &countryWithIP)
	if err != nil {
		return &CountryWithIP{}, err
	}

	return &countryWithIP, nil
}
