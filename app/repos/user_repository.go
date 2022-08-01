package repos

import (
	"context"

	"github.com/vonmutinda/organono/app/db"
	"github.com/vonmutinda/organono/app/entities"
	"github.com/vonmutinda/organono/app/utils"
)

const (
	getUserByIDSQL       = "SELECT id, first_name, last_name, username, password_hash, status, created_at, updated_at FROM users WHERE id = $1"
	getUserByUsernameSQL = "SELECT id, username, password, email, created_at, updated_at FROM users WHERE username = $1"
	saveUserSQL          = "INSERT INTO users (first_name, last_name, username, password_hash, status, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id"
	updateUserSQL        = "UPDATE users SET first_name = $1, last_name = $2, username = $3, password_hash = $4, status = $5, updated_at = $6 WHERE id = $7"
)

type (
	UserRepository interface {
		UserByID(ctx context.Context, operations db.SQLOperations, userID int64) (*entities.User, error)
		UserByUsername(ctx context.Context, operations db.SQLOperations, username string) (*entities.User, error)
		Save(ctx context.Context, operations db.SQLOperations, user *entities.User) error
	}

	AppUserRepository struct{}
)

func NewUserRepository() *AppUserRepository {
	return &AppUserRepository{}
}

func (r *AppUserRepository) UserByID(
	ctx context.Context,
	operations db.SQLOperations,
	userID int64,
) (*entities.User, error) {

	row := operations.QueryRowContext(
		ctx,
		getUserByIDSQL,
		userID,
	)

	user, err := r.scanRow(row)
	if err != nil {
		return &entities.User{}, utils.NewError(
			err,
			"user by id query row error",
		)
	}

	return user, nil
}

func (r *AppUserRepository) UserByUsername(
	ctx context.Context,
	operations db.SQLOperations,
	username string,
) (*entities.User, error) {

	row := operations.QueryRowContext(
		ctx,
		getUserByUsernameSQL,
		username,
	)

	user, err := r.scanRow(row)
	if err != nil {
		return user, utils.NewError(
			err,
			"user by username query row error",
		)
	}

	return user, nil
}

func (r *AppUserRepository) Save(
	ctx context.Context,
	operations db.SQLOperations,
	user *entities.User,
) error {

	user.Touch()

	if user.IsNew() {

		err := operations.QueryRowContext(
			ctx,
			saveUserSQL,
			user.FirstName,
			user.LastName,
			user.Username,
			user.PasswordHash,
			user.Status,
			user.CreatedAt,
			user.UpdatedAt,
		).Scan(
			&user.ID,
		)
		if err != nil {
			return utils.NewError(
				err,
				"save user query row error",
			)
		}

		return nil
	}

	_, err := operations.ExecContext(
		ctx,
		updateUserSQL,
		user.FirstName,
		user.LastName,
		user.Username,
		user.PasswordHash,
		user.Status,
		user.UpdatedAt,
		user.ID,
	)
	if err != nil {
		return utils.NewError(
			err,
			"update user exec context error",
		)
	}

	return nil
}

func (r *AppUserRepository) scanRow(
	rowScanner db.RowScanner,
) (*entities.User, error) {

	var user entities.User

	err := rowScanner.Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Username,
		&user.PasswordHash,
		&user.Status,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return &user, utils.NewError(
			err,
			"scan user row error",
		)
	}

	return &user, nil
}
