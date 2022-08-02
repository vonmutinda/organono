package utils

import (
	"context"
	"database/sql"

	"github.com/smartystreets/goconvey/convey"
	"github.com/vonmutinda/organono/app/db"
)

func WithTestDB(
	ctx context.Context,
	testDB db.DB,
	f func(context.Context, db.DB),
) func() {

	return func() {

		var dbTx *sql.Tx

		dB := db.NewTestDB(dbTx)

		if testDB.Valid() {
			_, err := testDB.ExecContext(ctx, "SET TRANSACTION ISOLATION LEVEL SERIALIZABLE")
			convey.So(err, convey.ShouldBeNil)

			dbTx, err = testDB.Begin()
			convey.So(err, convey.ShouldBeNil)
			dB = db.NewTestDB(dbTx)
		}

		convey.Reset(func() {
			if dbTx != nil {
				err := dbTx.Rollback()
				convey.So(err, convey.ShouldBeNil)
			}
		})

		f(ctx, dB)
	}
}
