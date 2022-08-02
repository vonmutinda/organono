package repos

import (
	"context"
	"testing"

	"github.com/vonmutinda/organono/app/db"
	"github.com/vonmutinda/organono/app/entities"
	"github.com/vonmutinda/organono/app/utils"

	. "github.com/smartystreets/goconvey/convey"
)

func TestUserRepository(t *testing.T) {

	testDB := db.InitDB()
	defer testDB.Close()

	userRepository := NewUserRepository()

	ctx := context.Background()

	Convey("User Repository", t, utils.WithTestDB(ctx, testDB, func(ctx context.Context, dB db.DB) {

		Convey("can save user", func() {

			user := entities.BuildUser()

			err := userRepository.Save(ctx, dB, user)
			So(err, ShouldBeNil)

			So(user.ID, ShouldNotBeZeroValue)
		})

		Convey("can get user by id", func() {

			user, err := CreateUser(ctx, dB)
			So(err, ShouldBeNil)

			foundUser, err := userRepository.UserByID(ctx, dB, user.ID)
			So(err, ShouldBeNil)

			So(foundUser.ID, ShouldEqual, user.ID)
			So(foundUser.FirstName, ShouldEqual, user.FirstName)
			So(foundUser.LastName, ShouldEqual, user.LastName)
			So(foundUser.Username, ShouldEqual, user.Username)
			So(foundUser.Status, ShouldEqual, user.Status)
			So(foundUser.CreatedAt, ShouldEqual, user.CreatedAt)
			So(foundUser.UpdatedAt, ShouldEqual, user.UpdatedAt)
			So(foundUser.PasswordHash, ShouldEqual, user.PasswordHash)
		})

		Convey("can update a user", func() {

			user, err := CreateUser(ctx, dB)
			So(err, ShouldBeNil)

			user.FirstName = "Trading"
			user.LastName = "Point"
			user.Username = "tradingpoint"
			user.Status = entities.UserStatusDeactivated

			err = userRepository.Save(ctx, dB, user)
			So(err, ShouldBeNil)

			foundUser, err := userRepository.UserByID(ctx, dB, user.ID)
			So(err, ShouldBeNil)

			So(foundUser.FirstName, ShouldEqual, user.FirstName)
			So(foundUser.LastName, ShouldEqual, user.LastName)
			So(foundUser.Username, ShouldEqual, user.Username)
			So(foundUser.Status, ShouldEqual, user.Status)
		})

		Convey("can get user by username", func() {

			user, err := CreateUser(ctx, dB)
			So(err, ShouldBeNil)

			foundUser, err := userRepository.UserByUsername(ctx, dB, user.Username)
			So(err, ShouldBeNil)

			So(foundUser.ID, ShouldEqual, user.ID)
		})
	}))
}
