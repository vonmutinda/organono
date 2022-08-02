package companies

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/vonmutinda/organono/app/db"
	"github.com/vonmutinda/organono/app/forms"
	"github.com/vonmutinda/organono/app/services"
	"github.com/vonmutinda/organono/app/utils"
	"github.com/vonmutinda/organono/app/web/webutils"
)

func createCompany(
	dB db.DB,
	companyService services.CompanyService,
) func(c *gin.Context) {
	return func(c *gin.Context) {

		var form forms.CreateCompanyForm

		err := c.BindJSON(&form)
		if err != nil {
			wrappedError := utils.NewErrorWithCode(
				err,
				utils.ErrorCodeInvalidForm,
				"Failed to bind create company form",
			)

			webutils.HandleError(c, wrappedError)
			return
		}

		ctx := c.Request.Context()

		company, err := companyService.CreateCompany(ctx, dB, &form)
		if err != nil {
			wrappedError := utils.NewError(
				err,
				"Failed to create company = [%+v]",
				form,
			)

			webutils.HandleError(c, wrappedError)
			return
		}

		c.JSON(http.StatusCreated, company)
	}
}

func deleteCompany(
	dB db.DB,
	companyService services.CompanyService,
) func(c *gin.Context) {
	return func(c *gin.Context) {

		companyID, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			wrappedError := utils.NewErrorWithCode(
				err,
				utils.ErrorCodeInvalidArgument,
				"Failed to parse company id = %v",
				c.Param("id"),
			)

			webutils.HandleError(c, wrappedError)
			return
		}

		ctx := c.Request.Context()

		company, err := companyService.DeleteCompany(ctx, dB, companyID)
		if err != nil {
			wrappedError := utils.NewError(
				err,
				"Failed to delete company id = %v",
				companyID,
			)

			webutils.HandleError(c, wrappedError)
			return
		}

		c.JSON(http.StatusOK, company)
	}
}

func getCompany(
	dB db.DB,
	companyService services.CompanyService,
) func(c *gin.Context) {
	return func(c *gin.Context) {

		companyID, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			wrappedError := utils.NewErrorWithCode(
				err,
				utils.ErrorCodeInvalidArgument,
				"Failed to parse company id = %v",
				c.Param("id"),
			)

			webutils.HandleError(c, wrappedError)
			return
		}

		ctx := c.Request.Context()

		company, err := companyService.GetCompany(ctx, dB, companyID)
		if err != nil {
			wrappedError := utils.NewError(
				err,
				"Failed to get company by id = %v",
				companyID,
			)

			webutils.HandleError(c, wrappedError)
			return
		}

		c.JSON(http.StatusOK, company)
	}
}

func listCompanies(
	dB db.DB,
	companyService services.CompanyService,
) func(c *gin.Context) {
	return func(c *gin.Context) {

		filter, err := webutils.FilterFromContext(c)
		if err != nil {
			appError := utils.NewError(
				err,
				"Failed to parse filter params from context",
			)

			webutils.HandleError(c, appError)
			return
		}

		ctx := c.Request.Context()

		companyList, err := companyService.ListCompanies(ctx, dB, filter)
		if err != nil {
			wrappedError := utils.NewError(
				err,
				"Failed to fetch companies",
			)

			webutils.HandleError(c, wrappedError)
			return
		}

		c.JSON(http.StatusOK, companyList)
	}
}

func updateCompany(
	dB db.DB,
	companyService services.CompanyService,
) func(c *gin.Context) {
	return func(c *gin.Context) {

		companyID, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			wrappedError := utils.NewErrorWithCode(
				err,
				utils.ErrorCodeInvalidArgument,
				"Failed to parse company id = %v",
				c.Param("id"),
			)

			webutils.HandleError(c, wrappedError)
			return
		}

		var form forms.UpdateCompanyForm

		err = c.BindJSON(&form)
		if err != nil {
			wrappedError := utils.NewErrorWithCode(
				err,
				utils.ErrorCodeInvalidForm,
				"Failed to bind update company form",
			)

			webutils.HandleError(c, wrappedError)
			return
		}

		ctx := c.Request.Context()

		company, err := companyService.UpdateCompany(ctx, dB, companyID, &form)
		if err != nil {
			wrappedError := utils.NewError(
				err,
				"Failed to update company id = %v form = [%+v]",
				companyID,
				form,
			)

			webutils.HandleError(c, wrappedError)
			return
		}

		c.JSON(http.StatusOK, company)
	}
}
