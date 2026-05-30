package services

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"solo/internal/models"
)

func (s *Services) CancelContract(ctx context.Context, req models.CancelContractRequest) error {
	idContrato := strings.TrimSpace(req.IDContrato)
	if idContrato == "" {
		return fmt.Errorf("%w: id_contrato obrigatorio", ErrEntradaInvalida)
	}

	result, err := s.repos.DB.ExecContext(ctx, `
		UPDATE users
		SET status_contrato = $1,
			ativo = false,
			atualizado_em = NOW()
		WHERE id_contrato = $2
	`, models.StatusCancelado, idContrato)
	if err != nil {
		s.logger.Error("erro ao cancelar contrato", "id_contrato", idContrato, "err", err)
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

func (s *Services) ContractStatus(ctx context.Context, idContrato string) (string, error) {
	idContrato = strings.TrimSpace(idContrato)
	if idContrato == "" {
		return "", fmt.Errorf("%w: id_contrato obrigatorio", ErrEntradaInvalida)
	}

	var status string
	err := s.repos.DB.QueryRowContext(ctx, `
		SELECT status_contrato
		FROM users
		WHERE id_contrato = $1
	`, idContrato).Scan(&status)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", ErrUsuarioNaoEncontrado
		}
		return "", err
	}

	return models.NormalizeStatusContrato(status), nil
}
