package services

import (
	"context"
	"errors"

	"github.com/vonmutinda/organono/app/db"
	"github.com/vonmutinda/organono/app/entities"
	"github.com/vonmutinda/organono/app/forms"
	"github.com/vonmutinda/organono/app/repos"
	"github.com/vonmutinda/organono/app/utils"
)

type (
	CompanyService interface {
		CreateCompany(ctx context.Context, dB db.DB, form *forms.CreateCompanyForm) (*entities.Company, error)
		DeleteCompany(ctx context.Context, dB db.DB, companyID int64) (*entities.Company, error)
		GetCompany(ctx context.Context, dB db.DB, companyID int64) (*entities.Company, error)
		ListCompanies(ctx context.Context, dB db.DB, filter *forms.Filter) (*entities.CompanyList, error)
		UpdateCompany(ctx context.Context, dB db.DB, companyID int64, form *forms.UpdateCompanyForm) (*entities.Company, error)
	}

	AppCompanyService struct {
		companyCountryRepository repos.CompanyCountryRepository
		companyRepository        repos.CompanyRepository
		countryReposistory       repos.CountryRepository
	}
)

func NewCompanyService(
	companyCountryRepository repos.CompanyCountryRepository,
	companyRepository repos.CompanyRepository,
	countryReposistory repos.CountryRepository,
) *AppCompanyService {
	return &AppCompanyService{
		companyCountryRepository: companyCountryRepository,
		companyRepository:        companyRepository,
		countryReposistory:       countryReposistory,
	}
}

func NewTestCompanyService() *AppCompanyService {
	return &AppCompanyService{
		companyCountryRepository: repos.NewCompanyCountryRepository(),
		companyRepository:        repos.NewCompanyRepository(),
		countryReposistory:       repos.NewCountryRepository(),
	}
}

func (s *AppCompanyService) CreateCompany(
	ctx context.Context,
	dB db.DB,
	form *forms.CreateCompanyForm,
) (*entities.Company, error) {

	err := s.validateCreateCompany(ctx, dB, form)
	if err != nil {
		return &entities.Company{}, err
	}

	country, err := s.getCountryByName(ctx, dB, form.Country)
	if err != nil {
		return &entities.Company{}, err
	}

	phoneNumber, err := utils.ParsePhoneNumber(form.Phone)
	if err != nil {
		return &entities.Company{}, err
	}

	company := &entities.Company{
		Name:        form.Name,
		Code:        form.Code,
		Website:     form.Website,
		Country:     country.Name,
		PhoneNumber: phoneNumber,
		Phone:       phoneNumber.Phone(),
	}

	err = dB.InTransaction(ctx, func(ctx context.Context, operations db.SQLOperations) error {

		err = s.companyRepository.Save(ctx, operations, company)
		if err != nil {
			return err
		}

		companyCountry := entities.CompanyCountry{
			CompanyID:       company.ID,
			CountryID:       country.ID,
			OperationStatus: entities.OperationStatusTypeActive,
		}

		return s.companyCountryRepository.Save(ctx, operations, &companyCountry)
	})
	if err != nil {
		return &entities.Company{}, err
	}

	return company, nil
}

func (s *AppCompanyService) DeleteCompany(
	ctx context.Context,
	dB db.DB,
	companyID int64,
) (*entities.Company, error) {

	company, err := s.companyRepository.CompanyByID(ctx, dB, companyID)
	if err != nil {
		if !utils.IsErrNoRows(err) {
			return &entities.Company{}, err
		}

		return &entities.Company{}, utils.NewErrorWithCode(
			err,
			utils.ErrorCodeNotFound,
			"company not found",
		)
	}

	err = s.companyRepository.DeleteCompany(ctx, dB, company.ID)
	if err != nil {
		return &entities.Company{}, err
	}

	return company, nil
}

func (s *AppCompanyService) GetCompany(
	ctx context.Context,
	dB db.DB,
	companyID int64,
) (*entities.Company, error) {

	company, err := s.companyRepository.CompanyByID(ctx, dB, companyID)
	if err != nil {
		if !utils.IsErrNoRows(err) {
			return &entities.Company{}, err
		}

		return &entities.Company{}, utils.NewErrorWithCode(
			err,
			utils.ErrorCodeNotFound,
			"company not found",
		)
	}

	return company, nil
}

func (s *AppCompanyService) ListCompanies(
	ctx context.Context,
	dB db.DB,
	filter *forms.Filter,
) (*entities.CompanyList, error) {

	companies, err := s.companyRepository.ListCompanies(ctx, dB, filter)
	if err != nil {
		return &entities.CompanyList{}, err
	}

	count, err := s.companyRepository.CompanyCount(ctx, dB, filter)
	if err != nil {
		return &entities.CompanyList{}, err
	}

	companyList := &entities.CompanyList{
		Companies:  companies,
		Pagination: entities.NewPagination(count, filter.Page, filter.Per),
	}

	return companyList, nil
}

func (s *AppCompanyService) UpdateCompany(
	ctx context.Context,
	dB db.DB,
	companyID int64,
	form *forms.UpdateCompanyForm,
) (*entities.Company, error) {

	company, err := s.companyRepository.CompanyByID(ctx, dB, companyID)
	if err != nil {
		return &entities.Company{}, err
	}

	err = s.validateUpdateCompany(ctx, dB, form)
	if err != nil {
		return &entities.Company{}, err
	}

	if form.Phone.Valid {
		phoneNumber, err := utils.ParsePhoneNumber(form.Phone.String)
		if err != nil {
			return &entities.Company{}, err
		}

		company.PhoneNumber = phoneNumber
		company.Phone = company.PhoneNumber.Phone()
	}

	if form.Name.Valid {
		company.Name = form.Name.String
	}

	if form.Code.Valid {
		company.Code = form.Code.String
	}

	if form.Website.Valid {
		company.Website = form.Website.String
	}

	err = s.companyRepository.Save(ctx, dB, company)
	if err != nil {
		return &entities.Company{}, err
	}

	return company, nil
}

func (s *AppCompanyService) validateCreateCompany(
	ctx context.Context,
	dB db.DB,
	form *forms.CreateCompanyForm,
) error {

	err := s.validateDuplicateCompanyCode(ctx, dB, form.Code)
	if err != nil {
		return err
	}

	err = s.validateDuplicateCompanyName(ctx, dB, form.Name)
	if err != nil {
		return err
	}

	err = s.validateDuplicateCompanyPhoneNumber(ctx, dB, form.Phone)
	if err != nil {
		return err
	}

	err = s.validateDuplicateCompanyWebsite(ctx, dB, form.Website)
	if err != nil {
		return err
	}

	return nil
}

func (s *AppCompanyService) validateUpdateCompany(
	ctx context.Context,
	dB db.DB,
	form *forms.UpdateCompanyForm,
) error {

	if form.Code.Valid {
		err := s.validateDuplicateCompanyCode(ctx, dB, form.Code.String)
		if err != nil {
			return err
		}
	}

	if form.Name.Valid {
		err := s.validateDuplicateCompanyName(ctx, dB, form.Name.String)
		if err != nil {
			return err
		}
	}

	if form.Phone.Valid {
		err := s.validateDuplicateCompanyPhoneNumber(ctx, dB, form.Phone.String)
		if err != nil {
			return err
		}
	}

	if form.Website.Valid {
		err := s.validateDuplicateCompanyWebsite(ctx, dB, form.Website.String)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *AppCompanyService) validateDuplicateCompanyCode(
	ctx context.Context,
	dB db.DB,
	code string,
) error {

	_, err := s.companyRepository.CompanyByCode(ctx, dB, code)
	if err != nil {
		if !utils.IsErrNoRows(err) {
			return err
		}
		return nil
	}

	return utils.NewErrorWithCode(
		errors.New("company code already exists"),
		utils.ErrorCodeResourceExists,
		"duplicate company code",
	)
}

func (s *AppCompanyService) validateDuplicateCompanyName(
	ctx context.Context,
	dB db.DB,
	name string,
) error {

	_, err := s.companyRepository.CompanyByName(ctx, dB, name)
	if err != nil {
		if !utils.IsErrNoRows(err) {
			return err
		}
		return nil
	}

	return utils.NewErrorWithCode(
		errors.New("company name already exists"),
		utils.ErrorCodeResourceExists,
		"duplicate company name",
	)
}

func (s *AppCompanyService) validateDuplicateCompanyPhoneNumber(
	ctx context.Context,
	dB db.DB,
	phone string,
) error {

	phoneNumber, err := utils.ParsePhoneNumber(phone)
	if err != nil {
		return err
	}

	_, err = s.companyRepository.CompanyByPhoneNumber(ctx, dB, phoneNumber)
	if err != nil {
		if !utils.IsErrNoRows(err) {
			return err
		}
		return nil
	}

	return utils.NewErrorWithCode(
		errors.New("company phone already exists"),
		utils.ErrorCodeResourceExists,
		"duplicate company phone",
	)
}

func (s *AppCompanyService) validateDuplicateCompanyWebsite(
	ctx context.Context,
	dB db.DB,
	code string,
) error {

	_, err := s.companyRepository.CompanyByWebsite(ctx, dB, code)
	if err != nil {
		if !utils.IsErrNoRows(err) {
			return err
		}
		return nil
	}

	return utils.NewErrorWithCode(
		errors.New("company website already exists"),
		utils.ErrorCodeResourceExists,
		"duplicate website code",
	)
}

func (s *AppCompanyService) getCountryByName(
	ctx context.Context,
	dB db.DB,
	countryName string,
) (*entities.Country, error) {

	country, err := s.countryReposistory.CountryByName(ctx, dB, countryName)
	if err != nil {
		if !utils.IsErrNoRows(err) {
			return &entities.Country{}, err
		}

		return &entities.Country{}, utils.NewErrorWithCode(
			err,
			utils.ErrorCodeNotFound,
			"country not found",
		)
	}

	return country, nil
}
