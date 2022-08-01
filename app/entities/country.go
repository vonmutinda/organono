package entities

type Country struct {
	SequentialIdentifier
	CountryCode  string `json:"country_code"`
	Currency     string `json:"currency"`
	Name         string `json:"name"`
	DiallingCode string `json:"dialling_code"`
}

func BuildCountry() *Country {
	return &Country{
		CountryCode:  "CY",
		Currency:     "EUR",
		Name:         "Cyprus",
		DiallingCode: "+357",
	}
}
