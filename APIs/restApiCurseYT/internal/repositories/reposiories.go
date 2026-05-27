package repositories

import (
	"rest-api-yt/internal/models"
	users "rest-api-yt/internal/repositories/users"
)

type Repositories struct {
	User interface {
		GetAll() []models.User
		Add(newUser models.User)
	}
}

func New() *Repositories {
	return &Repositories{
		User: users.New(),
	}
}
