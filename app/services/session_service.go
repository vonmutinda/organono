package services

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/vonmutinda/organono/app/db"
	"github.com/vonmutinda/organono/app/entities"
	"github.com/vonmutinda/organono/app/forms"
	"github.com/vonmutinda/organono/app/repos"
	"github.com/vonmutinda/organono/app/utils"
	"github.com/vonmutinda/organono/app/web/ctxhelper"
	"gopkg.in/guregu/null.v3"
)

type (
	SessionService interface {
		Login(ctx context.Context, dB db.DB, form *forms.UserLoginForm) (*entities.User, *entities.Session, error)
		Logout(ctx context.Context, dB db.DB, sessionID int64) error
		SessionByID(ctx context.Context, dB db.DB, sessionID int64) (*entities.Session, error)
	}

	AppSessionService struct {
		sessionRepository repos.SessionRepository
		userRepository    repos.UserRepository
	}
)

func NewSessionService(
	sessionRepository repos.SessionRepository,
	userRepository repos.UserRepository,
) *AppSessionService {
	return &AppSessionService{
		sessionRepository: sessionRepository,
		userRepository:    userRepository,
	}
}

func NewTestSessionService() *AppSessionService {
	return NewSessionService(
		repos.NewSessionRepository(),
		repos.NewUserRepository(),
	)
}

func (s *AppSessionService) SessionByID(
	ctx context.Context,
	dB db.DB,
	sessionID int64,
) (*entities.Session, error) {

	return s.sessionRepository.SessionByID(ctx, dB, sessionID)
}

func (s *AppSessionService) Login(
	ctx context.Context,
	dB db.DB,
	form *forms.UserLoginForm,
) (*entities.User, *entities.Session, error) {

	username := strings.TrimSpace(form.Username)
	password := strings.TrimSpace(form.Password)

	user, err := s.userRepository.UserByUsername(ctx, dB, username)
	if err != nil {
		if !utils.IsErrNoRows(err) {
			return &entities.User{}, &entities.Session{}, utils.NewError(
				err,
				"user by username =[%v]",
				username,
			)
		}

		return &entities.User{}, &entities.Session{}, utils.NewErrorWithCode(
			errors.New("invalid username"),
			utils.ErrorCodeInvalidCredentials,
			"username does not exist username=[%v]",
			username,
		)
	}

	err = utils.VerifyPassword(user.PasswordHash, password)
	if err != nil {
		return &entities.User{}, &entities.Session{}, utils.NewErrorWithCode(
			err,
			utils.ErrorCodeInvalidCredentials,
			"verify user password",
		)
	}

	if !user.Status.IsActive() {
		return &entities.User{}, &entities.Session{}, utils.NewErrorWithCode(
			errors.New("invalid user status"),
			utils.ErrorCodeInvalidUserStatus,
			"invalid status for user = %v",
			user.ID,
		)
	}

	session := &entities.Session{
		IPAddress:       ctxhelper.IPAddress(ctx),
		LastRefreshedAt: time.Now(),
		UserAgent:       ctxhelper.UserAgent(ctx),
		UserID:          user.ID,
	}

	err = s.sessionRepository.Save(ctx, dB, session)
	if err != nil {
		return &entities.User{}, &entities.Session{}, err
	}

	return user, session, nil

}

func (s *AppSessionService) Logout(
	ctx context.Context,
	dB db.DB,
	sessionID int64,
) error {

	tokenInfo := ctxhelper.TokenInfo(ctx)

	if sessionID != tokenInfo.SessionID {
		return utils.NewErrorWithCode(
			errors.New("unauthorized"),
			utils.ErrorCodeInvalidCredentials,
			"Cannot logout session=[%v] by current session=[%v]",
			sessionID,
			tokenInfo.SessionID,
		)
	}

	session, err := s.sessionRepository.SessionByID(ctx, dB, sessionID)
	if err != nil {
		return err
	}

	if session.UserID != tokenInfo.UserID {
		return utils.NewErrorWithCode(
			errors.New("unauthorized"),
			utils.ErrorCodeInvalidCredentials,
			"Cannot logout session=[%v] by user=[%v]",
			sessionID,
			tokenInfo.UserID,
		)
	}

	if session.DeactivatedAt.Valid {
		return nil
	}

	session.DeactivatedAt = null.TimeFrom(time.Now())

	return s.sessionRepository.Save(ctx, dB, session)
}
