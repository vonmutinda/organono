package repos

import (
	"context"
	"errors"

	"github.com/vonmutinda/organono/app/db"
	"github.com/vonmutinda/organono/app/entities"
	"github.com/vonmutinda/organono/app/utils"
)

const (
	saveCompanyCountrySQL = "INSERT INTO company_countries (company_id, country_id, operation_status, created_at, updated_at) VALUES ($1, $2, $3, $4, $5) RETURNING id"
)

type (
	CompanyCountryRepository interface {
		Save(ctx context.Context, operations db.SQLOperations, companyCountry *entities.CompanyCountry) error
	}

	AppCompanyCountryRepository struct{}
)

func NewCompanyCountryRepository() *AppCompanyCountryRepository {
	return &AppCompanyCountryRepository{}
}

func (r *AppCompanyCountryRepository) Save(
	ctx context.Context,
	operations db.SQLOperations,
	companyCountry *entities.CompanyCountry,
) error {

	companyCountry.Touch()

	if companyCountry.IsNew() {

		err := operations.QueryRowContext(
			ctx,
			saveCompanyCountrySQL,
			companyCountry.CompanyID,
			companyCountry.CountryID,
			companyCountry.OperationStatus,
			companyCountry.CreatedAt,
			companyCountry.UpdatedAt,
		).Scan(
			&companyCountry.ID,
		)
		if err != nil {
			return utils.NewError(
				err,
				"save company country query row error",
			)
		}

		return nil
	}

	return errors.New("cannot update company country")
}
