package repos

import (
	"context"
	"testing"

	"github.com/vonmutinda/organono/app/db"
	"github.com/vonmutinda/organono/app/entities"
	"github.com/vonmutinda/organono/app/utils"

	. "github.com/smartystreets/goconvey/convey"
)

func TestCompanyCountryRepository(t *testing.T) {

	testDB := db.InitDB()
	defer testDB.Close()

	companyCountryRepository := NewCompanyCountryRepository()
	companyRepository := NewCompanyRepository()

	ctx := context.Background()

	Convey("Company Country Repository", t, utils.WithTestDB(ctx, testDB, func(ctx context.Context, dB db.DB) {

		country, err := CreateCountry(ctx, dB)
		So(err, ShouldBeNil)

		Convey("can save a company country", func() {

			company := entities.BuildCompany("Microsoft", country)

			err = companyRepository.Save(ctx, dB, company)
			So(err, ShouldBeNil)

			So(company.ID, ShouldNotBeZeroValue)

			companyCountry := entities.BuildCompanyCountry(company.ID, country.ID)

			err := companyCountryRepository.Save(ctx, dB, companyCountry)
			So(err, ShouldBeNil)

			So(companyCountry.ID, ShouldNotBeZeroValue)
		})
	}))
}
