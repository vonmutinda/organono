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
		ListCompanies(ctx context.Context, dB db.DB, filter *forms.Filter) (*entities.CompanyList, error)
		UpdateCompany(ctx context.Context, dB db.DB, companyID int64, form *forms.UpdateCompanyForm) (*entities.Company, error)
	}

	AppCompanyService struct {
		companyCountryRepository repos.CompanyCountryRepository
		companyRepository        repos.CompanyRepository
		countryReposistory       repos.CountryRepository
	}
)

func NewCompanyCountryService(
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

func NewTestCompanyCountryService() *AppCompanyService {
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

	company := &entities.Company{
		Name:    form.Name,
		Code:    form.Code,
		Website: form.Website,
	}

	err = dB.InTransaction(ctx, func(ctx context.Context, operations db.SQLOperations) error {

		err = s.companyRepository.Save(ctx, operations, company)
		if err != nil {
			return err
		}

		companyCountry := entities.CompanyCountry{
			CompanyID: company.ID,
			CountryID: form.CountryID,
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

	_, err := s.countryReposistory.CountryByID(ctx, dB, form.CountryID)
	if err != nil {
		if !utils.IsErrNoRows(err) {
			return err
		}

		return utils.NewErrorWithCode(
			err,
			utils.ErrorCodeNotFound,
			"country not found",
		)
	}

	err = s.validateCompanyCodeExists(ctx, dB, form.Code)
	if err != nil {
		return err
	}

	err = s.validateCompanyNameExists(ctx, dB, form.Name)
	if err != nil {
		return err
	}

	err = s.validateCompanyPhoneNumberExists(ctx, dB, form.Phone)
	if err != nil {
		return err
	}

	err = s.validateCompanyWebsiteExists(ctx, dB, form.Website)
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
		err := s.validateCompanyCodeExists(ctx, dB, form.Code.String)
		if err != nil {
			return err
		}
	}

	if form.Name.Valid {
		err := s.validateCompanyNameExists(ctx, dB, form.Name.String)
		if err != nil {
			return err
		}
	}

	if form.Phone.Valid {
		err := s.validateCompanyPhoneNumberExists(ctx, dB, form.Phone.String)
		if err != nil {
			return err
		}
	}

	if form.Website.Valid {
		err := s.validateCompanyWebsiteExists(ctx, dB, form.Website.String)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *AppCompanyService) validateCompanyCodeExists(
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
		utils.ErrorCodeInvalidForm,
		"duplicate company code",
	)
}

func (s *AppCompanyService) validateCompanyNameExists(
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
		utils.ErrorCodeInvalidForm,
		"duplicate company name",
	)
}

func (s *AppCompanyService) validateCompanyPhoneNumberExists(
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
		utils.ErrorCodeInvalidForm,
		"duplicate company phone",
	)
}

func (s *AppCompanyService) validateCompanyWebsiteExists(
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
		utils.ErrorCodeInvalidForm,
		"duplicate website code",
	)
}
