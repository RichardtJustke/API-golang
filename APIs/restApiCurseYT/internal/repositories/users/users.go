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

func (u *Users) GetAll() []models.User {
	return u.users
}

func (u *Users) EmailExists(email string) bool {
	for _, v := range u.users {
		if v.Email == email {
			return true
		}
	}
	return false
}

//invenção minha função de verificação de idade
/*
func (u Users) AgeValid(age int)bool{
	if age <= 0 | age >= 100{
		print("idade não valida!")
		return
	}
}
*/

func (u *Users) Add(newUser models.User) {
	u.users = append(u.users, newUser)
}
