package repos

import (
	"context"
	"fmt"
	"strings"

	"github.com/vonmutinda/organono/app/db"
	"github.com/vonmutinda/organono/app/entities"
	"github.com/vonmutinda/organono/app/forms"
	"github.com/vonmutinda/organono/app/utils"
)

const (
	deleteCompanySQL           = "DELETE FROM companies WHERE id = $1"
	getCompaniesSQL            = "SELECT co.id, co.name, co.code, c.name, co.website, co.country_code, co.number, cc.operation_status, co.created_at, co.updated_at FROM companies co JOIN company_countries cc ON cc.company_id = co.id JOIN countries c ON c.id = cc.country_id"
	getCompanyByCodeSQL        = getCompaniesSQL + " WHERE co.code = $1"
	getCompanyByIDSQL          = getCompaniesSQL + " WHERE co.id = $1"
	getCompanyByNameSQL        = getCompaniesSQL + " WHERE co.name = $1"
	getCompanyByPhoneNumberSQL = getCompaniesSQL + " WHERE co.country_code = $1 AND co.number = $2"
	getCompanyByWebsiteSQL     = getCompaniesSQL + " WHERE co.website = $1"
	getCompanyCountSQL         = "SELECT COUNT(co.id) FROM companies co"
	saveCompanySQL             = "INSERT INTO companies (name, code, website, country_code, number, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id"
	updateCompanySQL           = "UPDATE companies SET name = $1, code = $2, website = $3, country_code = $4, number = $5, updated_at = $6 WHERE id = $7"
)

type (
	CompanyRepository interface {
		CompanyByCode(ctx context.Context, operations db.SQLOperations, code string) (*entities.Company, error)
		CompanyByID(ctx context.Context, operations db.SQLOperations, companyID int64) (*entities.Company, error)
		CompanyByName(ctx context.Context, operations db.SQLOperations, companyName string) (*entities.Company, error)
		CompanyByPhoneNumber(ctx context.Context, operations db.SQLOperations, phoneNumber entities.PhoneNumber) (*entities.Company, error)
		CompanyByWebsite(ctx context.Context, operations db.SQLOperations, website string) (*entities.Company, error)
		CompanyCount(ctx context.Context, operations db.SQLOperations, filter *forms.Filter) (int, error)
		DeleteCompany(ctx context.Context, operations db.SQLOperations, companyID int64) error
		ListCompanies(ctx context.Context, operations db.SQLOperations, filter *forms.Filter) ([]*entities.Company, error)
		Save(ctx context.Context, operations db.SQLOperations, company *entities.Company) error
	}

	AppCompanyRepository struct{}
)

func NewCompanyRepository() *AppCompanyRepository {
	return &AppCompanyRepository{}
}

func (r *AppCompanyRepository) CompanyByCode(
	ctx context.Context,
	operations db.SQLOperations,
	code string,
) (*entities.Company, error) {

	row := operations.QueryRowContext(
		ctx,
		getCompanyByCodeSQL,
		code,
	)

	return r.scanRow(row)
}

func (r *AppCompanyRepository) CompanyByID(
	ctx context.Context,
	operations db.SQLOperations,
	companyID int64,
) (*entities.Company, error) {

	row := operations.QueryRowContext(
		ctx,
		getCompanyByIDSQL,
		companyID,
	)

	return r.scanRow(row)
}

func (r *AppCompanyRepository) CompanyByName(
	ctx context.Context,
	operations db.SQLOperations,
	companyName string,
) (*entities.Company, error) {

	row := operations.QueryRowContext(
		ctx,
		getCompanyByNameSQL,
		companyName,
	)

	return r.scanRow(row)
}

func (r *AppCompanyRepository) CompanyByPhoneNumber(
	ctx context.Context,
	operations db.SQLOperations,
	phoneNumber entities.PhoneNumber,
) (*entities.Company, error) {

	row := operations.QueryRowContext(
		ctx,
		getCompanyByPhoneNumberSQL,
		phoneNumber.CountryCode,
		phoneNumber.Number,
	)

	return r.scanRow(row)
}

func (r *AppCompanyRepository) CompanyByWebsite(
	ctx context.Context,
	operations db.SQLOperations,
	website string,
) (*entities.Company, error) {

	row := operations.QueryRowContext(
		ctx,
		getCompanyByWebsiteSQL,
		website,
	)

	return r.scanRow(row)
}

func (r *AppCompanyRepository) CompanyCount(
	ctx context.Context,
	operations db.SQLOperations,
	filter *forms.Filter,
) (int, error) {

	var count int

	query, args := r.buildQuery(getCompanyCountSQL, filter.NoPagination())

	err := operations.QueryRowContext(
		ctx,
		query,
		args...,
	).Scan(&count)
	if err != nil {
		return 0, utils.NewError(
			err,
			"company count query row error",
		)
	}

	return count, nil
}

func (r *AppCompanyRepository) DeleteCompany(
	ctx context.Context,
	operations db.SQLOperations,
	companyID int64,
) error {

	_, err := operations.ExecContext(
		ctx,
		deleteCompanySQL,
		companyID,
	)
	if err != nil {
		return utils.NewError(
			err,
			"delete company exec context error",
		)
	}

	return nil
}

func (r *AppCompanyRepository) ListCompanies(
	ctx context.Context,
	operations db.SQLOperations,
	filter *forms.Filter,
) ([]*entities.Company, error) {

	query, args := r.buildQuery(getCompaniesSQL, filter)

	fmt.Printf("list companies query: %v\n\nargs = %v\n\n", query, args)

	rows, err := operations.QueryContext(ctx, query, args...)
	if err != nil {
		return []*entities.Company{}, utils.NewError(
			err,
			"list companies query context error",
		)
	}

	defer rows.Close()

	companies := make([]*entities.Company, 0)

	for rows.Next() {

		company, err := r.scanRow(rows)
		if err != nil {
			return []*entities.Company{}, err
		}

		companies = append(companies, company)
	}

	if rows.Err() != nil {
		return []*entities.Company{}, utils.NewError(
			rows.Err(),
			"list companies rows error",
		)
	}

	return companies, nil
}

func (r *AppCompanyRepository) Save(
	ctx context.Context,
	operations db.SQLOperations,
	company *entities.Company,
) error {

	company.Touch()

	if company.IsNew() {

		err := operations.QueryRowContext(
			ctx,
			saveCompanySQL,
			company.Name,
			company.Code,
			company.Website,
			company.PhoneNumber.CountryCode,
			company.PhoneNumber.Number,
			company.CreatedAt,
			company.UpdatedAt,
		).Scan(
			&company.ID,
		)
		if err != nil {
			return utils.NewError(
				err,
				"save company query row error",
			)
		}

		return nil
	}

	_, err := operations.ExecContext(
		ctx,
		updateCompanySQL,
		company.Name,
		company.Code,
		company.Website,
		company.PhoneNumber.CountryCode,
		company.PhoneNumber.Number,
		company.UpdatedAt,
		company.ID,
	)
	if err != nil {
		return utils.NewError(
			err,
			"update company exec error",
		)
	}

	return nil
}

func (r *AppCompanyRepository) buildQuery(
	query string,
	filter *forms.Filter,
) (string, []interface{}) {

	args := make([]interface{}, 0)
	conditions := make([]string, 0)
	counter := utils.NewPlaceholder()

	if filter.Term != "" {
		filterColumns := []string{"co.name", "co.code", "co.website", "c.name", "CONCAT(co.country_code, co.number)"}
		likeStatements := make([]string, 0)

		args = append(args, strings.ToLower(filter.Term))
		termPlaceholder := counter.Touch()
		for _, col := range filterColumns {
			stmt := fmt.Sprintf("LOWER(%s) LIKE '%%' || $%d || '%%'", col, termPlaceholder)
			likeStatements = append(likeStatements, stmt)
		}
		condition := fmt.Sprintf(" (%s)", strings.Join(likeStatements, " OR "))
		conditions = append(conditions, condition)
	}

	if filter.Status != "" {
		conditions = append(conditions, fmt.Sprintf(" cc.operation_status = $%d", counter.Touch()))
		args = append(args, filter.Status)
	}

	if len(conditions) > 0 {
		query += " WHERE" + strings.Join(conditions, " AND ")
	}

	if filter.Page > 0 && filter.Per > 0 {
		query += fmt.Sprintf(" ORDER BY co.name ASC LIMIT $%d OFFSET $%d", counter.Touch(), counter.Touch())
		args = append(args, filter.Per, (filter.Page-1)*filter.Per)
	}

	return query, args
}

func (r *AppCompanyRepository) scanRow(
	rowScanner db.RowScanner,
) (*entities.Company, error) {

	var company entities.Company

	err := rowScanner.Scan(
		&company.ID,
		&company.Name,
		&company.Code,
		&company.Country,
		&company.Website,
		&company.PhoneNumber.CountryCode,
		&company.PhoneNumber.Number,
		&company.OperationStatus,
		&company.CreatedAt,
		&company.UpdatedAt,
	)
	if err != nil {
		return &entities.Company{}, utils.NewError(
			err,
			"scan company row error",
		)
	}

	return &company, nil
}
