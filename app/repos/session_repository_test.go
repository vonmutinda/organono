package repos

import (
	"context"
	"testing"
	"time"

	"github.com/vonmutinda/organono/app/db"
	"github.com/vonmutinda/organono/app/entities"
	"github.com/vonmutinda/organono/app/utils"
	"gopkg.in/guregu/null.v3"
	"syreclabs.com/go/faker"

	. "github.com/smartystreets/goconvey/convey"
)

func TestSessionRepository(t *testing.T) {

	testDB := db.InitDB()
	defer testDB.Close()

	sessionRepository := NewSessionRepository()

	ctx := context.Background()

	Convey("Session Repository", t, utils.WithTestDB(ctx, testDB, func(ctx context.Context, dB db.DB) {

		user, err := CreateUser(ctx, dB)
		So(err, ShouldBeNil)

		Convey("can save session", func() {

			session := entities.BuildSession(user.ID)

			err := sessionRepository.Save(ctx, dB, session)
			So(err, ShouldBeNil)

			So(session.ID, ShouldNotBeZeroValue)
		})

		Convey("can get session by id", func() {

			session, err := CreateSession(ctx, dB, user.ID)
			So(err, ShouldBeNil)

			foundSession, err := sessionRepository.SessionByID(ctx, dB, session.ID)
			So(err, ShouldBeNil)

			So(foundSession.ID, ShouldEqual, session.ID)
			So(foundSession.IPAddress, ShouldEqual, session.IPAddress)
			So(foundSession.UserAgent, ShouldEqual, session.UserAgent)
			So(foundSession.LastRefreshedAt, ShouldEqual, session.LastRefreshedAt)
			So(foundSession.DeactivatedAt.Valid, ShouldBeFalse)
			So(foundSession.UserID, ShouldEqual, session.UserID)
			So(foundSession.CreatedAt, ShouldEqual, session.CreatedAt)
			So(foundSession.UpdatedAt, ShouldEqual, session.UpdatedAt)
		})

		Convey("can update a session", func() {

			session, err := CreateSession(ctx, dB, user.ID)
			So(err, ShouldBeNil)

			session.IPAddress = faker.Internet().IpV4Address()
			session.DeactivatedAt = null.TimeFrom(time.Now())

			err = sessionRepository.Save(ctx, dB, session)
			So(err, ShouldBeNil)

			foundSession, err := sessionRepository.SessionByID(ctx, dB, session.ID)
			So(err, ShouldBeNil)

			So(foundSession.ID, ShouldEqual, session.ID)
			So(foundSession.UserAgent, ShouldEqual, session.UserAgent)
			So(foundSession.DeactivatedAt.Valid, ShouldBeTrue)
		})
	}))
}
