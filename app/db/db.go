package db

import (
	"context"
	"database/sql"
	"os"

	_ "github.com/lib/pq"
	"github.com/vonmutinda/organono/app/logger"
)

var db DB

type DB interface {
	SQLOperations
	Begin() (*sql.Tx, error)
	Close() error
	Ping() error
	InTransaction(ctx context.Context, operations func(context.Context, SQLOperations) error) error
	Valid() bool
}

type RowScanner interface {
	Scan(dest ...interface{}) error
}

type SQLOperations interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
}

type AppDB struct {
	*sql.DB
	valid bool
}

type pgSQLOperations struct {
	*sql.Tx
}

func InitDB() DB {
	return InitDBWithURL(os.Getenv("DATABASE_URL"))
}

func InitDBWithURL(databaseURL string) DB {

	appDB := newPostgresDBWithURL(databaseURL)

	db = &AppDB{
		DB:    appDB,
		valid: true,
	}

	return db
}

func (db *AppDB) InTransaction(
	ctx context.Context,
	operations func(context.Context, SQLOperations) error,
) error {

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	sqlOperations := &pgSQLOperations{
		Tx: tx,
	}

	if err = operations(ctx, sqlOperations); err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return rollbackErr
		}

		return err
	}

	return tx.Commit()
}

func (db *AppDB) Valid() bool {
	return db.valid
}

func newPostgresDBWithURL(databaseURL string) *sql.DB {

	if databaseURL == "" {
		logger.Fatal("database url is required")
	}

	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		logger.Fatalf("sql.Open failed: %v", err)
	}

	return db
}

func GetDB() DB {
	return db
}
