package entities

type UserStatus string

const (
	UserStatusActive      UserStatus = "active"
	UserStatusDeactivated UserStatus = "deactivated"
	UserStatusUnverified  UserStatus = "unverified"
)

func (s UserStatus) String() string {
	return string(s)
}

func (s UserStatus) IsActive() bool {
	return s == UserStatusActive
}
