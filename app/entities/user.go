package entities

import "syreclabs.com/go/faker"

type User struct {
	SequentialIdentifier
	FirstName    string     `json:"first_name"`
	LastName     string     `json:"last_name"`
	Username     string     `json:"useranme"`
	PasswordHash string     `json:"-"`
	Status       UserStatus `json:"status"`
	Timestamps
}

func BuildUser() *User {
	return &User{
		FirstName: faker.Name().FirstName(),
		LastName:  faker.Name().LastName(),
		Username:  faker.RandomString(10),
		Status:    UserStatusActive,
	}
}
