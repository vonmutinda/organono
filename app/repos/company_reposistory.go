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
	getCompaniesSQL   = "SELECT co.id, co.name, c.name, co.website, co.country_code, co.number, co.created_at, co.updated_at FROM companies co JOIN company_countries cc ON cc.company_id = comp.id JOIN countries c ON c.id = cc.country_id"
	getCompanyByIDSQL = getCompaniesSQL + " WHERE id = $1"
	saveCompanySQL    = "INSERT INTO companies (name, created_at, updated_at) VALUES ($1, $2, $3)"
	updateCompanySQL  = "UPDATE companies SET name = $1, website = $2, country_code = $3, number = $4, updated_at = $5 WHERE id = $6"
)

type (
	CompanyRepository interface {
		CompanyByID(ctx context.Context, operations db.SQLOperations, companyID int64) (*entities.Company, error)
		ListCompanies(ctx context.Context, operations db.SQLOperations, filter *forms.Filter) ([]*entities.Company, error)
		Save(ctx context.Context, operations db.SQLOperations, company *entities.Company) error
	}

	AppCompanyRepository struct{}
)

func NewCompanyRepository() *AppCompanyRepository {
	return &AppCompanyRepository{}
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

func (r *AppCompanyRepository) ListCompanies(
	ctx context.Context,
	operations db.SQLOperations,
	filter *forms.Filter,
) ([]*entities.Company, error) {

	query, args := r.buildQuery(getCompaniesSQL, filter)

	rows, err := operations.QueryContext(
		ctx,
		query,
		args...,
	)
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
		company.Website,
		company.PhoneNumber.CountryCode,
		company.PhoneNumber.Number,
		company.UpdatedAt,
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

	conditions := make([]string, 0)
	args := make([]interface{}, 0)
	counter := utils.NewPlaceholder()

	if filter.Term != "" {
		filterColumns := []string{"co.name", "c.name", "co.website", "CONCAT(co.country_code, co.number)"}
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
		conditions = append(conditions, fmt.Sprintf(" cc.status = $%d", counter.Touch()))
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
		&company.PhoneNumber.CountryCode,
		&company.PhoneNumber.Number,
		&company.Website,
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
