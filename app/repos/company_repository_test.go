package repos

import (
	"context"
	"testing"

	"github.com/vonmutinda/organono/app/db"
	"github.com/vonmutinda/organono/app/entities"
	"github.com/vonmutinda/organono/app/forms"
	"github.com/vonmutinda/organono/app/utils"
	"gopkg.in/guregu/null.v3"

	. "github.com/smartystreets/goconvey/convey"
)

func TestCompanyRepository(t *testing.T) {

	testDB := db.InitDB()
	defer testDB.Close()

	companyCountryRepository := NewCompanyCountryRepository()
	companyRepository := NewCompanyRepository()

	ctx := context.Background()

	Convey("Company Repository", t, utils.WithTestDB(ctx, testDB, func(ctx context.Context, dB db.DB) {

		country, err := CreateCountry(ctx, dB)
		So(err, ShouldBeNil)

		Convey("can save a company", func() {

			company := entities.BuildCompany("Trading Point LLC", country)

			err := companyRepository.Save(ctx, dB, company)
			So(err, ShouldBeNil)

			So(company.ID, ShouldNotBeZeroValue)
		})

		Convey("can get a country by id", func() {

			company, _, err := CreateCompany(ctx, dB, "Trading Point LLC", country)
			So(err, ShouldBeNil)

			foundCompany, err := companyRepository.CompanyByID(ctx, dB, company.ID)
			So(err, ShouldBeNil)

			So(foundCompany.ID, ShouldEqual, company.ID)
			So(foundCompany.Name, ShouldEqual, company.Name)
			So(foundCompany.Country, ShouldEqual, country.Name)
			So(foundCompany.Name, ShouldEqual, company.Name)
			So(foundCompany.Code, ShouldEqual, company.Code)
			So(foundCompany.CreatedAt, ShouldNotBeZeroValue)
			So(foundCompany.UpdatedAt, ShouldNotBeZeroValue)
		})

		Convey("can get a company by code", func() {

			company, _, err := CreateCompany(ctx, dB, "Trading Point LLC", country)
			So(err, ShouldBeNil)

			foundCompany, err := companyRepository.CompanyByCode(ctx, dB, company.Code)
			So(err, ShouldBeNil)

			So(foundCompany.ID, ShouldEqual, company.ID)
		})

		Convey("can get a company by name", func() {

			company, _, err := CreateCompany(ctx, dB, "Trading Point LLC", country)
			So(err, ShouldBeNil)

			foundCompany, err := companyRepository.CompanyByCode(ctx, dB, company.Code)
			So(err, ShouldBeNil)

			So(foundCompany.ID, ShouldEqual, company.ID)
		})

		Convey("can get a company by phone number", func() {

			company, _, err := CreateCompany(ctx, dB, "Trading Point LLC", country)
			So(err, ShouldBeNil)

			foundCompany, err := companyRepository.CompanyByPhoneNumber(ctx, dB, company.PhoneNumber)
			So(err, ShouldBeNil)

			So(foundCompany.ID, ShouldEqual, company.ID)
		})

		Convey("can get a company by website", func() {

			company, _, err := CreateCompany(ctx, dB, "Trading Point LLC", country)
			So(err, ShouldBeNil)

			foundCompany, err := companyRepository.CompanyByWebsite(ctx, dB, company.Website)
			So(err, ShouldBeNil)

			So(foundCompany.ID, ShouldEqual, company.ID)
		})

		Convey("can update a company", func() {

			company, _, err := CreateCompany(ctx, dB, "Trading Point LLC", country)
			So(err, ShouldBeNil)

			company.Name = "Trading"
			company.Code = "TP"
			company.Website = "https://tradingpoint.com"
			company.PhoneNumber = entities.PhoneNumber{
				CountryCode: null.StringFrom("+357"),
				Number:      null.StringFrom("93456789"),
			}

			err = companyRepository.Save(ctx, dB, company)
			So(err, ShouldBeNil)

			foundCompany, err := companyRepository.CompanyByID(ctx, dB, company.ID)
			So(err, ShouldBeNil)

			So(foundCompany.ID, ShouldEqual, company.ID)
			So(foundCompany.Name, ShouldEqual, company.Name)
			So(foundCompany.Code, ShouldEqual, company.Code)
			So(foundCompany.Website, ShouldEqual, company.Website)
			So(foundCompany.PhoneNumber.CountryCode.String, ShouldEqual, company.PhoneNumber.CountryCode.String)
			So(foundCompany.PhoneNumber.Number.String, ShouldEqual, company.PhoneNumber.Number.String)
		})

		Convey("can list companies based on filter", func() {

			tradingPoint, _, err := CreateCompany(ctx, dB, "Trading Point LLC", country)
			So(err, ShouldBeNil)

			xmLTD, _, err := CreateCompany(ctx, dB, "XM LTD", country)
			So(err, ShouldBeNil)

			foundCompanies, err := companyRepository.ListCompanies(ctx, dB, &forms.Filter{})
			So(err, ShouldBeNil)
			So(len(foundCompanies), ShouldEqual, 2)

			So(foundCompanies[0].ID, ShouldEqual, tradingPoint.ID)
			So(foundCompanies[1].ID, ShouldEqual, xmLTD.ID)
		})

		Convey("can list companies based on filter - page 1 per 1", func() {

			tradingPoint, _, err := CreateCompany(ctx, dB, "Trading Point LLC", country)
			So(err, ShouldBeNil)

			_, _, err = CreateCompany(ctx, dB, "XM LTD", country)
			So(err, ShouldBeNil)

			filter := &forms.Filter{
				Page: 1,
				Per:  1,
			}

			foundCompanies, err := companyRepository.ListCompanies(ctx, dB, filter)
			So(err, ShouldBeNil)
			So(len(foundCompanies), ShouldEqual, 1)
			So(foundCompanies[0].ID, ShouldEqual, tradingPoint.ID)
		})

		Convey("can list companies based on filter - page 2 per 1", func() {

			_, _, err := CreateCompany(ctx, dB, "Trading Point LLC", country)
			So(err, ShouldBeNil)

			xmLTD, _, err := CreateCompany(ctx, dB, "XM LTD", country)
			So(err, ShouldBeNil)

			filter := &forms.Filter{
				Page: 2,
				Per:  1,
			}

			foundCompanies, err := companyRepository.ListCompanies(ctx, dB, filter)
			So(err, ShouldBeNil)
			So(len(foundCompanies), ShouldEqual, 1)
			So(foundCompanies[0].ID, ShouldEqual, xmLTD.ID)
		})

		Convey("can list companies based on term filter - name", func() {

			_, _, err := CreateCompany(ctx, dB, "Google", country)
			So(err, ShouldBeNil)

			appleCyprus, _, err := CreateCompany(ctx, dB, "Apple", country)
			So(err, ShouldBeNil)

			filter := &forms.Filter{
				Term: "app",
			}

			foundCompanies, err := companyRepository.ListCompanies(ctx, dB, filter)
			So(err, ShouldBeNil)
			So(len(foundCompanies), ShouldEqual, 1)
			So(foundCompanies[0].ID, ShouldEqual, appleCyprus.ID)
		})

		Convey("can list companies based on term filter - code", func() {

			_, _, err := CreateCompany(ctx, dB, "Google", country)
			So(err, ShouldBeNil)

			appleCyprus, _, err := CreateCompany(ctx, dB, "Apple", country)
			So(err, ShouldBeNil)

			filter := &forms.Filter{
				Term: appleCyprus.Code[:len(appleCyprus.Code)-1],
			}

			foundCompanies, err := companyRepository.ListCompanies(ctx, dB, filter)
			So(err, ShouldBeNil)
			So(len(foundCompanies), ShouldEqual, 1)
			So(foundCompanies[0].ID, ShouldEqual, appleCyprus.ID)
		})

		Convey("can list companies based on status filter", func() {

			closedCompany := entities.BuildCompany("Safari", country)
			err := companyRepository.Save(ctx, dB, closedCompany)
			So(err, ShouldBeNil)

			closedCyprusCountry := entities.BuildCompanyCountry(closedCompany.ID, country.ID)
			closedCyprusCountry.OperationStatus = entities.OperationStatusTypeClosed

			err = companyCountryRepository.Save(ctx, dB, closedCyprusCountry)
			So(err, ShouldBeNil)

			activeCompany := entities.BuildCompany("Redbull", country)
			err = companyRepository.Save(ctx, dB, activeCompany)
			So(err, ShouldBeNil)

			filter := &forms.Filter{
				Status: "closed",
			}

			foundCompanies, err := companyRepository.ListCompanies(ctx, dB, filter)
			So(err, ShouldBeNil)
			So(len(foundCompanies), ShouldEqual, 1)
			So(foundCompanies[0].ID, ShouldEqual, closedCompany.ID)
		})

		Convey("can count companies", func() {

			_, _, err := CreateCompany(ctx, dB, "Google", country)
			So(err, ShouldBeNil)

			_, _, err = CreateCompany(ctx, dB, "Apple", country)
			So(err, ShouldBeNil)

			count, err := companyRepository.CompanyCount(ctx, dB, &forms.Filter{})
			So(err, ShouldBeNil)

			So(count, ShouldEqual, 2)
		})
	}))
}
