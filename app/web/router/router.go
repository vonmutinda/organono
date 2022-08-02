package router

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/vonmutinda/organono/app/db"
	"github.com/vonmutinda/organono/app/repos"
	"github.com/vonmutinda/organono/app/services"
	"github.com/vonmutinda/organono/app/web/api/companies"
	"github.com/vonmutinda/organono/app/web/api/sessions"
	"github.com/vonmutinda/organono/app/web/auth"
	"github.com/vonmutinda/organono/app/web/middleware"
)

type AppRouter struct {
	*gin.Engine
}

func BuildRouter(
	dB db.DB,
) *AppRouter {

	if os.Getenv("ENVIRONMENT") == "development" {
		gin.SetMode(gin.DebugMode)
	}

	router := gin.Default()

	// router versions
	appV1Router := router.Group("/v1")

	sessionRepository := repos.NewSessionRepository()
	userRepository := repos.NewUserRepository()
	sessionAuthenticator := auth.NewSessionAuthenticator(sessionRepository, userRepository)

	defaultMiddlewares := middleware.DefaultMiddlewares(sessionAuthenticator)
	router.Use(defaultMiddlewares...)

	// Repositories
	companyCountryRepository := repos.NewCompanyCountryRepository()
	companyRepository := repos.NewCompanyRepository()
	countryRepository := repos.NewCountryRepository()

	// Services
	companyService := services.NewCompanyService(
		companyCountryRepository,
		companyRepository,
		countryRepository,
	)
	sessionService := services.NewSessionService(sessionRepository, userRepository)

	// Open endpoints
	unauthenticatedUsers := appV1Router.Group("")
	sessions.AddOpenEndpoints(unauthenticatedUsers, dB, sessionAuthenticator, sessionService)

	// User endpoints
	activeUsers := appV1Router.Group("")
	activeUsers.Use(auth.AllowOnlyActiveUser(dB, sessionAuthenticator, sessionService))

	sessions.AddEndpoints(activeUsers, dB, sessionService)
	companies.AddEndpoints(activeUsers, dB, companyService)

	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"error_message": "Endpoint not found"})
	})

	return &AppRouter{router}
}
