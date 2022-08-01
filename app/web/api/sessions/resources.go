package sessions

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vonmutinda/organono/app/db"
	"github.com/vonmutinda/organono/app/forms"
	"github.com/vonmutinda/organono/app/services"
	"github.com/vonmutinda/organono/app/utils"
	"github.com/vonmutinda/organono/app/web/auth"
	"github.com/vonmutinda/organono/app/web/ctxhelper"
	"github.com/vonmutinda/organono/app/web/webutils"
)

func login(
	dB db.DB,
	sessionAuthenticator auth.SessionAuthenticator,
	sessionService services.SessionService,
) func(c *gin.Context) {
	return func(c *gin.Context) {

		var form forms.UserLoginForm

		err := c.BindJSON(&form)
		if err != nil {
			wrappedError := utils.NewErrorWithCode(
				err,
				utils.ErrorCodeInvalidForm,
				"Failed to bind login form",
			)

			webutils.HandleError(c, wrappedError)
			return
		}

		ctx := c.Request.Context()

		user, session, err := sessionService.Login(ctx, dB, &form)
		if err != nil {
			wrappedError := utils.NewError(
				err,
				"Failed to log in username = [%v]",
				form.Username,
			)

			webutils.HandleError(c, wrappedError)
			return
		}

		sessionAuthenticator.SetUserSessionInResponse(c.Writer, user, session)

		c.JSON(http.StatusOK, gin.H{"success": true})
	}
}

func logout(
	dB db.DB,
	sessionService services.SessionService,
) func(c *gin.Context) {
	return func(c *gin.Context) {

		ctx := c.Request.Context()

		tokenInfo := ctxhelper.TokenInfo(ctx)

		err := sessionService.Logout(ctx, dB, tokenInfo.SessionID)
		if err != nil {
			wrappedError := utils.NewError(
				err,
				"Failed to log out userID=[%v]",
				tokenInfo.UserID,
			)

			webutils.HandleError(c, wrappedError)
			return
		}

		c.JSON(http.StatusOK, gin.H{"success": true})
	}
}
