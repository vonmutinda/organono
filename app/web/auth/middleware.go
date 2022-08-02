package auth

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/vonmutinda/organono/app/db"
	"github.com/vonmutinda/organono/app/services"
	"github.com/vonmutinda/organono/app/utils"
	"github.com/vonmutinda/organono/app/web/ctxhelper"
)

func AllowOnlyActiveUser(
	dB db.DB,
	sessionAuthenticator SessionAuthenticator,
	sessionService services.SessionService,
) func(c *gin.Context) {

	return func(c *gin.Context) {

		ctx := c.Request.Context()

		err := validateCyprusIPAddress(c, sessionAuthenticator)
		if err != nil {
			wrappedError := utils.NewError(
				err,
				"Failed to check is cyprus ip address",
			).Notify()

			wrappedError.WithContext(c.Request.Context())
			wrappedError.LogErrorMessages()
			c.JSON(wrappedError.HttpStatus(), wrappedError.JsonResponse())
			c.Abort()
			return
		}

		err = validateSession(c, dB, sessionAuthenticator, sessionService)
		if err != nil {
			wrappedError := utils.NewError(
				err,
				"Failed to validate session",
			).Notify()

			wrappedError.WithContext(c.Request.Context())
			wrappedError.LogErrorMessages()
			c.JSON(wrappedError.HttpStatus(), wrappedError.JsonResponse())
			c.Abort()
			return
		}

		tokenInfo := ctxhelper.TokenInfo(ctx)

		if tokenInfo.RequiresRefresh() {
			_, err = sessionAuthenticator.RefreshTokenFromRequest(ctx, dB, tokenInfo, c.Writer)
			if err != nil {
				wrappedError := utils.NewError(
					err,
					"Failed to refresh token from request",
				).Notify()

				wrappedError.WithContext(c.Request.Context())
				wrappedError.LogErrorMessages()
				c.JSON(wrappedError.HttpStatus(), wrappedError.JsonResponse())
				c.Abort()
				return
			}
		}
	}
}

func validateCyprusIPAddress(
	c *gin.Context,
	sessionAuthenticator SessionAuthenticator,
) error {

	isCyprus, err := sessionAuthenticator.IsCyrpusIPAddress(c.Request.Context())
	if err != nil {
		return err
	}

	if !isCyprus {
		return nil
	}

	exemptedMethodURLMap := map[string]string{
		http.MethodPost:   "/v1/companies",
		http.MethodDelete: "/v1/company/%v",
	}

	url, ok := exemptedMethodURLMap[c.Request.Method]
	if !ok {
		return nil
	}

	if strings.Contains(url, "%v") {
		url = fmt.Sprintf(url, c.Param("id"))
	}

	if url == c.Request.URL.Path {
		c.Next()
	}

	return nil
}

func validateSession(
	c *gin.Context,
	dB db.DB,
	sessionAuthenticator SessionAuthenticator,
	sessionService services.SessionService,
) error {

	ctx := c.Request.Context()
	tokenInfo := ctxhelper.TokenInfo(ctx)

	session, err := sessionService.SessionByID(ctx, dB, tokenInfo.SessionID)
	if err != nil {
		return utils.NewError(
			err,
			"Failed to get sessionID =[%v]",
			tokenInfo.SessionID,
		)
	}

	if session.UserID != tokenInfo.UserID {
		return utils.NewErrorWithCode(
			errors.New("invalid user id"),
			utils.ErrorCodeRoleForbidden,
			"Failed to check user role for sessionID=[%v] with userID=[%v], and user=[%v]",
			tokenInfo.SessionID,
			session.UserID,
			tokenInfo.UserID,
		)
	}

	user, err := sessionAuthenticator.UserByID(ctx, dB, session.UserID)
	if err != nil {
		return err
	}

	if !user.Status.IsActive() {
		return utils.NewErrorWithCode(
			errors.New("inactive user status"),
			utils.ErrorCodeInvalidUserStatus,
			"Failed to check for active user=[%v]",
			tokenInfo.UserID,
		)
	}

	if session.DeactivatedAt.Valid {
		return utils.NewErrorWithCode(
			errors.New("inactive session"),
			utils.ErrorCodeSessionExpired,
			"Failed to check for active user for sessionID=[%v]",
			tokenInfo.SessionID,
		)
	}

	return nil
}
