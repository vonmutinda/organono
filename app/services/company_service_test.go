package services

import (
	"context"
	"testing"

	"github.com/vonmutinda/organono/app/db"
	"github.com/vonmutinda/organono/app/forms"
	"github.com/vonmutinda/organono/app/repos"
	"github.com/vonmutinda/organono/app/utils"
	"gopkg.in/guregu/null.v3"

	. "github.com/smartystreets/goconvey/convey"
)

func TestCompanyServiceRepository(t *testing.T) {

	testDB := db.InitDB()
	defer testDB.Close()

	ctx := context.Background()

	companyRepository := repos.NewCompanyRepository()

	companyService := NewTestCompanyService()

	Convey("Company Service", t, utils.WithTestDB(ctx, testDB, func(ctx context.Context, dB db.DB) {

		country, err := repos.CreateCountry(ctx, dB)
		So(err, ShouldBeNil)
		So(country.Name, ShouldEqual, "Cyprus")

		Convey("can create a company", func() {

			form := &forms.CreateCompanyForm{
				Name:    "Microsoft",
				Country: "cyprus",
				Website: "https://microsoft.com",
				Phone:   "+35790034567",
			}

			_, err := companyService.CreateCompany(ctx, dB, form)
			So(err, ShouldBeNil)
		})

		Convey("cannot create a company with an existing name", func() {

			company, _, err := repos.CreateCompany(ctx, dB, "KFC", country)
			So(err, ShouldBeNil)

			form := &forms.CreateCompanyForm{
				Name: company.Name,
			}

			_, err = companyService.CreateCompany(ctx, dB, form)
			So(err, ShouldNotBeNil)

			appError, ok := err.(*utils.Error)
			So(ok, ShouldBeTrue)
			So(appError.GetErrorCode(), ShouldEqual, utils.ErrorCodeResourceExists)
		})

		Convey("cannot create a company with an existing code", func() {

			company, _, err := repos.CreateCompany(ctx, dB, "KFC", country)
			So(err, ShouldBeNil)

			form := &forms.CreateCompanyForm{
				Code: company.Code,
			}

			_, err = companyService.CreateCompany(ctx, dB, form)
			So(err, ShouldNotBeNil)

			appError, ok := err.(*utils.Error)
			So(ok, ShouldBeTrue)
			So(appError.GetErrorCode(), ShouldEqual, utils.ErrorCodeResourceExists)
		})

		Convey("cannot create a company with an existing website", func() {

			company, _, err := repos.CreateCompany(ctx, dB, "KFC", country)
			So(err, ShouldBeNil)

			form := &forms.CreateCompanyForm{
				Name:    "Microsoft",
				Country: "cyprus",
				Website: company.Website,
				Phone:   "+35790034567",
			}

			_, err = companyService.CreateCompany(ctx, dB, form)
			So(err, ShouldNotBeNil)

			appError, ok := err.(*utils.Error)
			So(ok, ShouldBeTrue)
			So(appError.GetErrorCode(), ShouldEqual, utils.ErrorCodeResourceExists)
		})

		Convey("cannot create a company with an existing phone number", func() {

			company, _, err := repos.CreateCompany(ctx, dB, "KFC", country)
			So(err, ShouldBeNil)

			form := &forms.CreateCompanyForm{
				Name:    "Microsoft",
				Country: "cyprus",
				Website: "https://microsoft.com",
				Phone:   company.PhoneNumber.Phone(),
			}

			_, err = companyService.CreateCompany(ctx, dB, form)
			So(err, ShouldNotBeNil)

			appError, ok := err.(*utils.Error)
			So(ok, ShouldBeTrue)
			So(appError.GetErrorCode(), ShouldEqual, utils.ErrorCodeResourceExists)
		})

		Convey("can update a company", func() {

			company, _, err := repos.CreateCompany(ctx, dB, "KFC", country)
			So(err, ShouldBeNil)

			form := &forms.UpdateCompanyForm{
				Name:    null.StringFrom("Burger King"),
				Code:    null.StringFrom("BKING357"),
				Website: null.StringFrom("https://kfc.com"),
				Phone:   null.StringFrom("+35790034567"),
			}

			updatedCompany, err := companyService.UpdateCompany(ctx, dB, company.ID, form)
			So(err, ShouldBeNil)

			So(updatedCompany.ID, ShouldEqual, company.ID)
			So(updatedCompany.Name, ShouldEqual, form.Name.String)
			So(updatedCompany.Code, ShouldEqual, form.Code.String)
			So(updatedCompany.Website, ShouldEqual, form.Website.String)
			So(updatedCompany.Phone, ShouldEqual, form.Phone.String)
		})

		Convey("cannot update a company with an existing name", func() {

			company, _, err := repos.CreateCompany(ctx, dB, "Burger King", country)
			So(err, ShouldBeNil)

			form := &forms.UpdateCompanyForm{
				Name: null.StringFrom("Burger King"),
			}

			_, err = companyService.UpdateCompany(ctx, dB, company.ID, form)
			So(err, ShouldNotBeNil)

			appError, ok := err.(*utils.Error)
			So(ok, ShouldBeTrue)
			So(appError.GetErrorCode(), ShouldEqual, utils.ErrorCodeResourceExists)
		})

		Convey("cannot update a company with an existing code", func() {

			company, _, err := repos.CreateCompany(ctx, dB, "Burger King", country)
			So(err, ShouldBeNil)

			form := &forms.UpdateCompanyForm{
				Code: null.StringFrom(company.Code),
			}

			_, err = companyService.UpdateCompany(ctx, dB, company.ID, form)
			So(err, ShouldNotBeNil)

			appError, ok := err.(*utils.Error)
			So(ok, ShouldBeTrue)
			So(appError.GetErrorCode(), ShouldEqual, utils.ErrorCodeResourceExists)
		})

		Convey("cannot update a company with an existing website", func() {

			company, _, err := repos.CreateCompany(ctx, dB, "Burger King", country)
			So(err, ShouldBeNil)

			form := &forms.UpdateCompanyForm{
				Website: null.StringFrom(company.Website),
			}

			_, err = companyService.UpdateCompany(ctx, dB, company.ID, form)
			So(err, ShouldNotBeNil)

			appError, ok := err.(*utils.Error)
			So(ok, ShouldBeTrue)
			So(appError.GetErrorCode(), ShouldEqual, utils.ErrorCodeResourceExists)
		})

		Convey("cannot update a company with an existing phone number", func() {

			company, _, err := repos.CreateCompany(ctx, dB, "Burger King", country)
			So(err, ShouldBeNil)

			form := &forms.UpdateCompanyForm{
				Phone: null.StringFrom(company.Phone),
			}

			_, err = companyService.UpdateCompany(ctx, dB, company.ID, form)
			So(err, ShouldNotBeNil)

			appError, ok := err.(*utils.Error)
			So(ok, ShouldBeTrue)
			So(appError.GetErrorCode(), ShouldEqual, utils.ErrorCodeResourceExists)
		})

		Convey("can list companies based on filter", func() {

			appleCyprus, _, err := repos.CreateCompany(ctx, dB, "Apple", country)
			So(err, ShouldBeNil)

			_, _, err = repos.CreateCompany(ctx, dB, "Microsoft", country)
			So(err, ShouldBeNil)

			googleCyprus, _, err := repos.CreateCompany(ctx, dB, "Google", country)
			So(err, ShouldBeNil)

			filter := &forms.Filter{
				Page: 1,
				Per:  2,
			}

			foundCompanyList, err := companyService.ListCompanies(ctx, dB, filter)
			So(err, ShouldBeNil)

			foundCompanies := foundCompanyList.Companies
			So(len(foundCompanies), ShouldEqual, 2)
			So(foundCompanies[0].ID, ShouldEqual, appleCyprus.ID)
			So(foundCompanies[1].ID, ShouldEqual, googleCyprus.ID)

			So(foundCompanyList.Pagination.Count, ShouldEqual, 3)
		})

		Convey("can delete a company", func() {

			company, _, err := repos.CreateCompany(ctx, dB, "KFC", country)
			So(err, ShouldBeNil)

			_, err = companyService.DeleteCompany(ctx, dB, company.ID)
			So(err, ShouldBeNil)

			_, err = companyRepository.CompanyByID(ctx, dB, company.ID)
			So(err, ShouldNotBeNil)
			So(utils.IsErrNoRows(err), ShouldBeTrue)
		})
	}))
}
