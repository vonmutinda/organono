package companies

import (
	"github.com/gin-gonic/gin"
	"github.com/vonmutinda/organono/app/db"
	"github.com/vonmutinda/organono/app/services"
)

func AddEndpoints(
	r *gin.RouterGroup,
	dB db.DB,
	companyService services.CompanyService,
) {
	r.POST("/companies", createCompany(dB, companyService))
	r.GET("/companies", listCompanies(dB, companyService))
	r.GET("/companies/:id", getCompany(dB, companyService))
	r.PUT("/companies/:id", updateCompany(dB, companyService))
	r.DELETE("/companies/:id", deleteCompany(dB, companyService))
}
