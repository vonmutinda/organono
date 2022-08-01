package repos

import (
	"context"

	"github.com/vonmutinda/organono/app/db"
	"github.com/vonmutinda/organono/app/entities"
)

const (
	getCountryByIDSQL = "SELECT id, country_code, currency, name, dialling_code FROM countries WHERE id = $1"
)

type (
	CountryRepository interface {
		CountryByID(ctx context.Context, operations db.SQLOperations, countryID int64) (*entities.Country, error)
	}

	AppCountryRepository struct{}
)

func NewCountryRepository() *AppCountryRepository {
	return &AppCountryRepository{}
}

func (r *AppCountryRepository) CountryByID(
	ctx context.Context,
	operations db.SQLOperations,
	countryID int64,
) (*entities.Country, error) {

	var country entities.Country

	err := operations.QueryRowContext(
		ctx,
		getCountryByIDSQL,
		countryID,
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
