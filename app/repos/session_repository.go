package repos

import (
	"context"

	"github.com/vonmutinda/organono/app/db"
	"github.com/vonmutinda/organono/app/entities"
	"github.com/vonmutinda/organono/app/utils"
)

const (
	getSessionByIDSQL = "SELECT s.id, s.deactivated_at, s.ip_address, s.last_refreshed_at, s.user_agent, u.status, s.user_id, s.created_at, s.updated_at FROM sessions s JOIN users u ON u.id = s.user_id WHERE s.id = $1"
	saveSessionSQL    = "INSERT INTO sessions (deactivated_at, ip_address, last_refreshed_at, user_agent, user_id, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id"
	updateSessionSQL  = "UPDATE sessions SET deactivated_at = $1, last_refreshed_at = $2, updated_at = $3 WHERE id = $4"
)

type (
	SessionRepository interface {
		Save(ctx context.Context, operations db.SQLOperations, session *entities.Session) error
		SessionByID(ctx context.Context, operations db.SQLOperations, sessionID int64) (*entities.Session, error)
	}

	AppSessionRepository struct{}
)

func NewSessionRepository() *AppSessionRepository {
	return &AppSessionRepository{}
}

func (r *AppSessionRepository) Save(
	ctx context.Context,
	operations db.SQLOperations,
	session *entities.Session,
) error {

	session.Touch()

	if session.IsNew() {

		err := operations.QueryRowContext(
			ctx,
			saveSessionSQL,
			session.DeactivatedAt,
			session.IPAddress,
			session.LastRefreshedAt,
			session.UserAgent,
			session.UserID,
			session.CreatedAt,
			session.UpdatedAt,
		).Scan(
			&session.ID,
		)
		if err != nil {
			return utils.NewError(
				err,
				"save session query row error",
			)
		}

		return nil
	}

	_, err := operations.ExecContext(
		ctx,
		updateSessionSQL,
		session.DeactivatedAt,
		session.LastRefreshedAt,
		session.UpdatedAt,
		session.ID,
	)
	if err != nil {
		return utils.NewError(
			err,
			"update session exec context error",
		)
	}

	return nil
}

func (r *AppSessionRepository) SessionByID(
	ctx context.Context,
	operations db.SQLOperations,
	sessionID int64,
) (*entities.Session, error) {

	row := operations.QueryRowContext(
		ctx,
		getSessionByIDSQL,
		sessionID,
	)

	session, err := r.scanRow(row)
	if err != nil {
		return &entities.Session{}, utils.NewError(
			err,
			"session by id query row error",
		)
	}

	return session, nil
}

func (r *AppSessionRepository) scanRow(
	rowScanner db.RowScanner,
) (*entities.Session, error) {

	var session entities.Session

	err := rowScanner.Scan(
		&session.ID,
		&session.DeactivatedAt,
		&session.IPAddress,
		&session.LastRefreshedAt,
		&session.UserAgent,
		&session.UserStatus,
		&session.UserID,
		&session.CreatedAt,
		&session.UpdatedAt,
	)

	if err != nil {
		return &session, utils.NewError(
			err,
			"session row scan",
		)
	}

	return &session, nil
}
