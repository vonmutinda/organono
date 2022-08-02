package auth

import (
	"errors"

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

		isCyprus, err := sessionAuthenticator.IsCyrpusIPAddress(ctx)
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

		if isCyprus {
			c.Next()
		}

		err = validateSession(c, dB, sessionAuthenticator, sessionService, "user")
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

func validateSession(
	c *gin.Context,
	dB db.DB,
	sessionAuthenticator SessionAuthenticator,
	sessionService services.SessionService,
	role string,
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
