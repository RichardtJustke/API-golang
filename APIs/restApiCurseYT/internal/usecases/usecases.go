package usecases

import (
	"rest-api-yt/internal/models"
	"rest-api-yt/internal/repositories"

	"github.com/google/uuid"
)

type UseCases struct {
	repos *repositories.Repositories
}

func New(repos *repositories.Repositories) {
	return &UseCases{
		repos: repos,
	}
}

func (u UseCases) GetAll() []models.User {
	users := u.repos.User.GetAll()

	return users
}

func (u UseCases) Add(newUser models.User) uuid.UUID {
	repoReq := models.User{
		ID:   uuid.New(),
		Name: newUser.Name,
		Age:  newUser.Age,
	}

	u.repos.User.Add(repoReq)
	return repoReq.ID
}
