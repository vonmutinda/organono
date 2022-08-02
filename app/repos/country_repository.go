package repos

import (
	"context"
	"strings"

	"github.com/vonmutinda/organono/app/db"
	"github.com/vonmutinda/organono/app/entities"
)

const (
	getCountryByNameSQL = "SELECT id, country_code, currency, name, dialling_code FROM countries WHERE LOWER(name) = $1"
)

type (
	CountryRepository interface {
		CountryByName(ctx context.Context, operations db.SQLOperations, countryName string) (*entities.Country, error)
	}

	AppCountryRepository struct{}
)

func NewCountryRepository() *AppCountryRepository {
	return &AppCountryRepository{}
}

func (r *AppCountryRepository) CountryByName(
	ctx context.Context,
	operations db.SQLOperations,
	countryName string,
) (*entities.Country, error) {

	var country entities.Country

	err := operations.QueryRowContext(
		ctx,
		getCountryByNameSQL,
		strings.ToLower(countryName),
	).Scan(
		&country.ID,
		&country.CountryCode,
		&country.Currency,
		&country.Name,
		&country.DiallingCode,
	)
	if err != nil {
		return &entities.Country{}, err
	}

	return &country, nil
}
