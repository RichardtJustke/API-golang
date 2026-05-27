package users

import (
	"rest-api-yt/internal/models"
)

type Users struct {
	users []models.User
}

func New() *Users {
	return &Users{users: make([]models.User, 0)}
}

func (u Users) GetAll() []models.User {
	return u.users
}

func (u Users) Add(newUser models.User) {
	u.users = append(u.users, newUser)
}
