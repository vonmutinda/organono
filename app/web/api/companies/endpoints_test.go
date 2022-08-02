package companies

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/vonmutinda/organono/app/db"
	"github.com/vonmutinda/organono/app/entities"
	"github.com/vonmutinda/organono/app/forms"
	"github.com/vonmutinda/organono/app/providers"
	"github.com/vonmutinda/organono/app/repos"
	"github.com/vonmutinda/organono/app/services"
	"github.com/vonmutinda/organono/app/utils"
	"github.com/vonmutinda/organono/app/web/auth"
	"github.com/vonmutinda/organono/app/web/ctxhelper"
	"github.com/vonmutinda/organono/app/web/middleware"
	"gopkg.in/guregu/null.v3"

	. "github.com/smartystreets/goconvey/convey"
)

func TestCompanyEndpointsRepository(t *testing.T) {

	testDB := db.InitDB()
	defer testDB.Close()

	ctx := context.Background()

	companyRepository := repos.NewCompanyRepository()
	sessionRepository := repos.NewSessionRepository()
	userRepository := repos.NewUserRepository()

	companyService := services.NewTestCompanyService()

	sessionAuthenticator := auth.NewSessionAuthenticator(
		providers.NewIPAPI(),
		sessionRepository,
		userRepository,
	)
	sessionService := services.NewTestSessionService()

	Convey("Company Endpoints", t, utils.WithTestDB(ctx, testDB, func(ctx context.Context, dB db.DB) {

		testRouter := gin.Default()
		testRouter.Use(middleware.DefaultMiddlewares(sessionAuthenticator)...)

		routerGroup := testRouter.Group("/v1")
		routerGroup.Use(auth.AllowOnlyActiveUser(
			dB,
			sessionAuthenticator,
			sessionService,
		))

		AddEndpoints(routerGroup, dB, companyService)

		country, err := repos.CreateCountry(ctx, dB)
		So(err, ShouldBeNil)
		So(country.Name, ShouldEqual, "Cyprus")

		Convey("protected endpoints", func() {

			user, err := repos.CreateUser(ctx, dB)
			So(err, ShouldBeNil)

			session, err := repos.CreateSession(ctx, dB, user.ID)
			So(err, ShouldBeNil)

			token, err := auth.NewJWTHandler().CreateUserToken(user, session)
			So(err, ShouldBeNil)

			Convey("can create a company", func() {

				form := &forms.CreateCompanyForm{
					Name:    "Microsoft",
					Code:    "SOFT",
					Country: "Cyprus",
					Website: "https://microsoft.com",
					Phone:   "+35790034567",
				}

				w, err := utils.DoRequest(testRouter, http.MethodPost, "/v1/companies", form, token)
				So(err, ShouldBeNil)

				So(w.Code, ShouldEqual, http.StatusCreated)
			})

			Convey("can get company by id", func() {

				company, _, err := repos.CreateCompany(ctx, dB, "Trading Point LLC", country)
				So(err, ShouldBeNil)

				w, err := utils.DoRequest(testRouter, http.MethodGet, fmt.Sprintf("/v1/companies/%d", company.ID), nil, token)
				So(err, ShouldBeNil)

				So(w.Code, ShouldEqual, http.StatusOK)

				data, err := ioutil.ReadAll(w.Body)
				So(err, ShouldBeNil)

				var foundCompany entities.Company

				err = json.Unmarshal(data, &foundCompany)
				So(err, ShouldBeNil)

				So(foundCompany.ID, ShouldEqual, company.ID)
				So(foundCompany.Name, ShouldEqual, company.Name)
				So(foundCompany.Country, ShouldEqual, country.Name)
				So(foundCompany.Website, ShouldEqual, company.Website)
				So(foundCompany.PhoneNumber.Phone(), ShouldEqual, company.Phone)
				So(foundCompany.Code, ShouldEqual, company.Code)
			})

			Convey("can update a company", func() {

				company, _, err := repos.CreateCompany(ctx, dB, "Trading Point LLC", country)
				So(err, ShouldBeNil)

				form := &forms.UpdateCompanyForm{
					Name:    null.StringFrom("Microsoft"),
					Website: null.StringFrom("https://microsoft.com"),
					Phone:   null.StringFrom("+35790034567"),
				}

				w, err := utils.DoRequest(testRouter, http.MethodPut, fmt.Sprintf("/v1/companies/%v", company.ID), form, token)
				So(err, ShouldBeNil)

				So(w.Code, ShouldEqual, http.StatusOK)
			})

			Convey("can delete a company", func() {

				company, _, err := repos.CreateCompany(ctx, dB, "Trading Point LLC", country)
				So(err, ShouldBeNil)

				w, err := utils.DoRequest(testRouter, http.MethodDelete, fmt.Sprintf("/v1/companies/%v", company.ID), nil, token)
				So(err, ShouldBeNil)

				So(w.Code, ShouldEqual, http.StatusOK)

				_, err = companyRepository.CompanyByID(ctx, dB, company.ID)
				So(err, ShouldNotBeNil)
				So(utils.IsErrNoRows(err), ShouldBeTrue)
			})

			Convey("can list companies", func() {

				_, _, err := repos.CreateCompany(ctx, dB, "Trading Point LLC", country)
				So(err, ShouldBeNil)

				_, _, err = repos.CreateCompany(ctx, dB, "XM LTD", country)
				So(err, ShouldBeNil)

				w, err := utils.DoRequest(testRouter, http.MethodGet, "/v1/companies", nil, token)
				So(err, ShouldBeNil)

				So(w.Code, ShouldEqual, http.StatusOK)

				data, err := ioutil.ReadAll(w.Body)
				So(err, ShouldBeNil)

				var companyList entities.CompanyList

				err = json.Unmarshal(data, &companyList)
				So(err, ShouldBeNil)

				So(len(companyList.Companies), ShouldEqual, 2)
			})

			Convey("can list companies based on search", func() {

				tradingPointCyprus, _, err := repos.CreateCompany(ctx, dB, "Trading Point LLC", country)
				So(err, ShouldBeNil)

				_, _, err = repos.CreateCompany(ctx, dB, "XM LTD", country)
				So(err, ShouldBeNil)

				w, err := utils.DoRequest(testRouter, http.MethodGet, "/v1/companies?term=point", nil, token)
				So(err, ShouldBeNil)

				So(w.Code, ShouldEqual, http.StatusOK)

				data, err := ioutil.ReadAll(w.Body)
				So(err, ShouldBeNil)

				var companyList entities.CompanyList

				err = json.Unmarshal(data, &companyList)
				So(err, ShouldBeNil)

				So(len(companyList.Companies), ShouldEqual, 1)
				So(companyList.Companies[0].ID, ShouldEqual, tradingPointCyprus.ID)
			})
		})

		Convey("unauthenticated endpoints", func() {

			Convey("can create a company from cyprus unauthenticated", func() {

				form := &forms.CreateCompanyForm{
					Name:    "Microsoft",
					Code:    "SOFT",
					Country: "Cyprus",
					Website: "https://microsoft.com",
					Phone:   "+35790034567",
				}

				b, err := json.Marshal(form)
				So(err, ShouldBeNil)

				ctx = ctxhelper.WithIpAddress(ctx, "176.56.168.0") // cyprus ip address

				req, err := http.NewRequest(http.MethodPost, "/v1/companies", bytes.NewReader(b))
				So(err, ShouldBeNil)

				req = req.WithContext(ctx)
				w := httptest.NewRecorder()

				testRouter.ServeHTTP(w, req)
				So(w.Code, ShouldEqual, http.StatusCreated)
			})

			Convey("cannot create a company from kenya unauthenticated", func() {

				form := &forms.CreateCompanyForm{
					Name:    "Microsoft",
					Code:    "SOFT",
					Country: "Kenya",
					Website: "https://microsoft.com",
					Phone:   "+35790034567",
				}

				b, err := json.Marshal(form)
				So(err, ShouldBeNil)

				ctx = ctxhelper.WithIpAddress(ctx, "197.156.118.251") // kenya ip address

				req, err := http.NewRequest(http.MethodPost, "/v1/companies", bytes.NewReader(b))
				So(err, ShouldBeNil)

				req = req.WithContext(ctx)
				w := httptest.NewRecorder()

				testRouter.ServeHTTP(w, req)
				So(w.Code, ShouldEqual, http.StatusBadRequest)
			})
		})
	}))
}
