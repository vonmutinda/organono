package auth

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/vonmutinda/organono/app/db"
	"github.com/vonmutinda/organono/app/entities"
	"github.com/vonmutinda/organono/app/providers"
	"github.com/vonmutinda/organono/app/repos"
	"github.com/vonmutinda/organono/app/utils"
	"github.com/vonmutinda/organono/app/web/ctxhelper"
)

const (
	tokenHeader = "X-ORGANONO-TOKEN"
)

var ErrTokenNotProvided = errors.New("token not provided")

type SessionAuthenticator interface {
	IsCyrpusIPAddress(ctx context.Context) (bool, error)
	RefreshTokenFromRequest(ctx context.Context, dB db.DB, tokenInfo *entities.TokenInfo, w http.ResponseWriter) (string, error)
	SetUserSessionInResponse(w http.ResponseWriter, user *entities.User, session *entities.Session) (string, error)
	TokenInfoFromRequest(req *http.Request) (*entities.TokenInfo, error)
	UserByID(ctx context.Context, dB db.DB, userID int64) (*entities.User, error)
}

type AppSessionAuthenticator struct {
	ipAPIProvider     providers.IPAPI
	jwtHandler        JWTHandler
	sessionRepository repos.SessionRepository
	userRepository    repos.UserRepository
}

func NewSessionAuthenticator(
	ipAPIProvider providers.IPAPI,
	sessionRepository repos.SessionRepository,
	userRepository repos.UserRepository,
) SessionAuthenticator {
	return NewSessionAuthenticatorWithJWTHandler(
		ipAPIProvider,
		NewJWTHandler(),
		sessionRepository,
		userRepository,
	)
}

func NewSessionAuthenticatorWithJWTHandler(
	ipAPIProvider providers.IPAPI,
	jwtHandler JWTHandler,
	sessionRepository repos.SessionRepository,
	userRepository repos.UserRepository,
) SessionAuthenticator {
	return &AppSessionAuthenticator{
		ipAPIProvider:     ipAPIProvider,
		jwtHandler:        jwtHandler,
		sessionRepository: sessionRepository,
		userRepository:    userRepository,
	}
}

func (a *AppSessionAuthenticator) IsCyrpusIPAddress(
	ctx context.Context,
) (bool, error) {

	lookupIPAddress := ctxhelper.IPAddress(ctx)

	if lookupIPAddress == "" {
		return false, nil
	}

	countryWithIp, err := a.ipAPIProvider.CountryForIP(lookupIPAddress)
	if err != nil {
		return false, err
	}

	if strings.EqualFold(countryWithIp.CountryName, "Cyprus") {
		return true, nil
	}

	return false, nil
}

func (a *AppSessionAuthenticator) RefreshTokenFromRequest(
	ctx context.Context,
	dB db.DB,
	tokenInfo *entities.TokenInfo,
	w http.ResponseWriter,
) (string, error) {

	session, err := a.sessionRepository.SessionByID(ctx, dB, tokenInfo.SessionID)
	if err != nil {
		return "", utils.NewError(
			err,
			"Failed to find session by id=[%v]",
			tokenInfo.SessionID,
		)
	}

	if session.DeactivatedAt.Valid {
		return "", utils.NewErrorWithCode(
			errors.New("session expired"),
			utils.ErrorCodeSessionExpired,
			"Failed to refresh session by id=[%v]",
			tokenInfo.SessionID,
		)
	}

	user, err := a.userRepository.UserByID(ctx, dB, session.UserID)
	if err != nil {
		return "", utils.NewError(
			err,
			"Failed to find user by id=[%v]",
			session.UserID,
		)
	}

	if !user.Status.IsActive() {
		return "", utils.NewErrorWithCode(
			errors.New("invalid status"),
			utils.ErrorCodeInvalidUserStatus,
			"Failed to refresh session id=[%v] for user id=[%v]",
			tokenInfo.SessionID,
			session.UserID,
		)
	}

	session.LastRefreshedAt = time.Now()

	err = a.sessionRepository.Save(ctx, dB, session)
	if err != nil {
		utils.NewError(
			err,
			"Failed to update last refreshed at for session id =[%v]",
			session.ID,
		).LogErrorMessages()
	}

	return a.SetUserSessionInResponse(w, user, session)
}

func (a *AppSessionAuthenticator) SetUserSessionInResponse(
	w http.ResponseWriter,
	user *entities.User,
	session *entities.Session,
) (string, error) {

	tokenValue, err := a.jwtHandler.CreateUserToken(user, session)
	if err != nil {
		return "", utils.NewError(
			err,
			"Failed to generate signed user token for user=[%v]",
			user.ID,
		)
	}

	w.Header().Set(tokenHeader, tokenValue)

	return tokenValue, nil
}

func (a *AppSessionAuthenticator) TokenInfoFromRequest(
	req *http.Request,
) (*entities.TokenInfo, error) {

	tokenValue := req.Header.Get(tokenHeader)
	if tokenValue != "" {
		return a.jwtHandler.TokenInfo(tokenValue)
	}

	return &entities.TokenInfo{}, ErrTokenNotProvided
}

func (a *AppSessionAuthenticator) UserByID(
	ctx context.Context,
	dB db.DB,
	userID int64,
) (*entities.User, error) {

	user, err := a.userRepository.UserByID(ctx, dB, userID)
	if err != nil {
		return user, utils.NewError(
			err,
			"Failed to find user by id=[%v]",
			userID,
		)
	}

	return user, nil
}
