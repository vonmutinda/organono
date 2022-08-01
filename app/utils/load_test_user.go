package utils

import (
	"context"

	"github.com/vonmutinda/organono/app/db"
	"github.com/vonmutinda/organono/app/entities"
	"github.com/vonmutinda/organono/app/logger"
)

var (
	user = &entities.User{
		FirstName: "Trading",
		LastName:  "Point",
		Username:  "xm",
		Status:    "active",
	}

	defaultPassword = "password"

	saveUserSQL = `
		INSERT INTO users 
			(first_name, last_name, username, password_hash, status, created_at, updated_at) 
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (username) DO NOTHING
		`
)

func LoadTestData(dB db.SQLOperations) error {

	hash, err := GeneratePasswordHash(defaultPassword)
	if err != nil {
		logger.Fatalf("load_test_data: generate passsword hash err = %v", err)
	}

	user.PasswordHash = string(hash)

	err = dB.QueryRowContext(
		context.Background(),
		saveUserSQL,
		user.FirstName,
		user.LastName,
		user.Username,
		user.PasswordHash,
		user.Status,
		user.CreatedAt,
		user.UpdatedAt,
	).Err()
	if err != nil {
		logger.Fatalf("load_test_data: save user query row error = %v", err)

	}

	return nil
}
