package services

import (
	"errors"
	"log/slog"

	"solo/internal/repositories"
)

var (
	ErrEntradaInvalida      = errors.New("entrada invalida")
	ErrStatusInvalido       = errors.New("status de contrato invalido")
	ErrUsuarioNaoEncontrado = errors.New("usuario nao encontrado")
	ErrUsuarioDuplicado     = errors.New("usuario ja cadastrado")
)

type Services struct {
	repos  *repositories.Repositories
	logger *slog.Logger
}

func New(repos *repositories.Repositories) *Services {
	return &Services{
		repos:  repos,
		logger: slog.Default(),
	}
}
