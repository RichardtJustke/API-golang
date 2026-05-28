package usecases

import (
	"errors"
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
	exist := u.repos.User.EmailExists(newUser.Email)
	if exist{
		slog.Error("thiss user already exists", "email", newUser.Email)

		return uuid.Nil, errors.New("user already exist")
	}
	repoReq := models.User{
		ID:   uuid.New(),
		Name: newUser.Name,
		Age:  newUser.Age,
	}

	u.repos.User.Add(repoReq)

	return repoReq.ID, nil
}
