package entities

import "syreclabs.com/go/faker"

type Company struct {
	SequentialIdentifier
	Name        string      `json:"name"`
	Code        string      `json:"code"`
	Country     string      `json:"country"`
	Website     string      `json:"website"`
	Phone       string      `json:"phone"`
	PhoneNumber PhoneNumber `json:"phone_number"`
	Timestamps
}

type CompanyList struct {
	Companies  []*Company  `json:"companies"`
	Pagination *Pagination `json:"pagination"`
}

func BuildCompany(country *Country) *Company {

	phoneNumber := fakePhoneNumber()

	return &Company{
		Name:        faker.Company().Name(),
		Country:     country.Name,
		Website:     faker.Internet().Url(),
		Phone:       phoneNumber.Phone(),
		PhoneNumber: phoneNumber,
	}
}
