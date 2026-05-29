package usecases

import (
	"errors"
	"log/slog"
	"rest-api-yt/internal/models"
	"rest-api-yt/internal/repositories"

	"github.com/google/uuid"
)

type UseCases struct {
	repos *repositories.Repositories
}

var ErrUserAlreadyExists = errors.New("user already exist")

func New(repos *repositories.Repositories) *UseCases {
	return &UseCases{
		repos: repos,
	}
}

func (u UseCases) GetAll() []models.User {
	users := u.repos.User.GetAll()

	return users
}

func (u UseCases) Add(newUser models.CreatUserRequest) (uuid.UUID, error) {
	exist := u.repos.User.EmailExists(newUser.Email)
	if exist {
		slog.Error("thiss user already exists", "email", newUser.Email)

		return uuid.Nil, ErrUserAlreadyExists
	}
	repoReq := models.User{
		ID:    uuid.New(),
		Name:  newUser.Name,
		Email: newUser.Email,
	}

	u.repos.User.Add(repoReq)

	return repoReq.ID, nil
}
