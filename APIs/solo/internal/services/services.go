package services

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"solo/internal/models"
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

// UpdateContractStatus atualiza o status do contrato e marca ativo conforme o status.
func (s *Services) UpdateContractStatus(ctx context.Context, idContrato, status string) error {
	idContrato = strings.TrimSpace(idContrato)
	if idContrato == "" {
		return fmt.Errorf("%w: id_contrato obrigatorio", ErrEntradaInvalida)
	}

	status = models.NormalizeStatusContrato(status)
	if status == "" || !models.IsStatusContratoValido(status) {
		return fmt.Errorf("%w: %s", ErrStatusInvalido, status)
	}

	ativo := status == models.StatusAtivo

	result, err := s.repos.DB.ExecContext(ctx, `
		UPDATE usuarios
		SET status_contrato = $1,
			ativo = $2,
			atualizado_em = NOW()
		WHERE id_contrato = $3
	`, status, ativo, idContrato)
	if err != nil {
		s.logger.Error("erro ao atualizar status contrato", "id_contrato", idContrato, "err", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrUsuarioNaoEncontrado
	}

	return nil
}
