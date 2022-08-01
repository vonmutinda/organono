package entities

import (
	"time"

	"gopkg.in/guregu/null.v3"
	"syreclabs.com/go/faker"
)

type Session struct {
	SequentialIdentifier
	DeactivatedAt   null.Time  `json:"deactivated_at"`
	IPAddress       string     `json:"ip_address"`
	LastRefreshedAt time.Time  `json:"last_refreshed_at"`
	UserAgent       string     `json:"user_agent"`
	UserID          int64      `json:"user_id"`
	UserStatus      UserStatus `json:"user_status"`
	Timestamps
}

func BuildSession(UserID int64) *Session {
	return &Session{
		IPAddress:       faker.Internet().IpV4Address(),
		LastRefreshedAt: time.Now(),
		UserAgent:       faker.Internet().Slug(),
		UserID:          UserID,
	}
}
