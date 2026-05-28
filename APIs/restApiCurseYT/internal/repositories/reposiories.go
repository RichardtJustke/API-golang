package repositories

import (
	"rest-api-yt/internal/models"
	users "rest-api-yt/internal/repositories/users"
)

type Repositories struct {
	User interface {
		GetAll() []models.User
		Add(newUser models.User)
		EmailExists(email string) bool
	}
}

func New() *Repositories {
	return &Repositories{
		User: users.New(),
	}
}
