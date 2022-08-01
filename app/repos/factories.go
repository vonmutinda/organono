package repos

import (
	"context"

	"github.com/vonmutinda/organono/app/db"
	"github.com/vonmutinda/organono/app/entities"
)

func CreateCompany(ctx context.Context, dB db.DB, country *entities.Country) (*entities.Company, *entities.CompanyCountry, error) {
	company := entities.BuildCompany(country)
	err := NewCompanyRepository().Save(ctx, dB, company)
	if err != nil {
		return &entities.Company{}, &entities.CompanyCountry{}, err
	}
	companyCountry, err := CreateCompanyCountry(ctx, dB, company.ID, country.ID)
	return company, companyCountry, err
}

func CreateCompanyCountry(ctx context.Context, dB db.DB, companyID, countryID int64) (*entities.CompanyCountry, error) {
	companyCountry := entities.BuildCompanyCountry(companyID, countryID)
	err := NewCompanyCountryRepository().Save(ctx, dB, companyCountry)
	return companyCountry, err
}

func CreateSession(ctx context.Context, dB db.DB, userID int64) (*entities.Session, error) {
	session := entities.BuildSession(userID)
	err := NewSessionRepository().Save(ctx, dB, session)
	return session, err
}

func CreateUser(ctx context.Context, dB db.DB) (*entities.User, error) {
	user := entities.BuildUser()
	err := NewUserRepository().Save(ctx, dB, user)
	return user, err
}
