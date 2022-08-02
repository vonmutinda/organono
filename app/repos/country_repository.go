package repos

import (
	"context"
	"errors"
	"strings"

	"github.com/vonmutinda/organono/app/db"
	"github.com/vonmutinda/organono/app/entities"
	"github.com/vonmutinda/organono/app/utils"
)

const (
	getCountryByNameSQL = "SELECT id, code, currency, name, dialling_code FROM countries WHERE LOWER(name) = $1"
	saveCountrySQL      = "INSERT INTO countries (code, currency, name, dialling_code) VALUES ($1, $2, $3, $4) RETURNING id"
)

type (
	CountryRepository interface {
		CountryByName(ctx context.Context, operations db.SQLOperations, countryName string) (*entities.Country, error)
		Save(ctx context.Context, operations db.SQLOperations, country *entities.Country) error
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

func (r *AppCountryRepository) Save(
	ctx context.Context,
	operations db.SQLOperations,
	country *entities.Country,
) error {

	if country.IsNew() {

		err := operations.QueryRowContext(
			ctx,
			saveCountrySQL,
			country.CountryCode,
			country.Currency,
			country.Name,
			country.DiallingCode,
		).Scan(&country.ID)
		if err != nil {
			return utils.NewError(
				err,
				"save country query row context error",
			)
		}

		return nil
	}

	return errors.New("cannot update acountry")
}
