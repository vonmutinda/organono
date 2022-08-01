package sessions

import (
	"github.com/gin-gonic/gin"
	"github.com/vonmutinda/organono/app/db"
	"github.com/vonmutinda/organono/app/services"
	"github.com/vonmutinda/organono/app/web/auth"
)

func AddEndpoints(
	r *gin.RouterGroup,
	dB db.DB,
	sessionsService services.SessionService,
) {
	r.DELETE("/auth", logout(dB, sessionsService))
}

func AddOpenEndpoints(
	r *gin.RouterGroup,
	dB db.DB,
	sessionAuthenticator auth.SessionAuthenticator,
	sessionService services.SessionService,
) {
	r.POST("/auth", login(dB, sessionAuthenticator, sessionService))
}
