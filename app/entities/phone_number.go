package entities

import (
	"strings"

	"gopkg.in/guregu/null.v3"
	"syreclabs.com/go/faker"
)

type PhoneNumber struct {
	CountryCode null.String `json:"country_code"`
	Number      null.String `json:"number"`
}

func (pn PhoneNumber) Phone() string {
	return strings.TrimSpace(pn.CountryCode.String + pn.Number.String)
}

func fakePhoneNumber() PhoneNumber {
	return PhoneNumber{
		CountryCode: null.StringFrom("+357"),
		Number:      null.StringFrom("9003" + faker.Number().Number(4)),
	}
}
