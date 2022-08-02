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
			Timeout: time.Second * 5,
		},
	}
}

func (p *AppIPAPI) CountryForIP(ipAddress string) (*CountryWithIP, error) {

	response, err := http.Get(fmt.Sprintf(ipAPIURL, ipAddress))
	if err != nil {
		return &CountryWithIP{}, fmt.Errorf("failed to get response err = %v", err)
	}

	defer response.Body.Close()

	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return &CountryWithIP{}, fmt.Errorf("failed to read response body = %v", err)
	}

	var countryWithIP CountryWithIP

	err = json.Unmarshal(data, &countryWithIP)
	if err != nil {
		return &CountryWithIP{}, fmt.Errorf("failed to unmarshal response body err = %v", err)
	}

	return &countryWithIP, nil
}
