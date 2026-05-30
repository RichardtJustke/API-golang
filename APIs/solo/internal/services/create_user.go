package services

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgconn"

	"solo/internal/models"
)

func (s *Services) CreateUser(ctx context.Context, req models.CreateUserRequest) (*models.User, error) {
	req.Nome = strings.TrimSpace(req.Nome)
	req.CPF = strings.TrimSpace(req.CPF)
	req.Email = strings.TrimSpace(req.Email)
	req.SenhaHash = strings.TrimSpace(req.SenhaHash)
	req.IDContrato = strings.TrimSpace(req.IDContrato)
	req.StatusContrato = models.NormalizeStatusContrato(req.StatusContrato)

	if req.Nome == "" || req.CPF == "" || req.Email == "" || req.SenhaHash == "" {
		return nil, fmt.Errorf("%w: nome, cpf, email e senha_hash sao obrigatorios", ErrEntradaInvalida)
	}

	if req.DataNascimento.IsZero() {
		return nil, fmt.Errorf("%w: data_nascimento obrigatoria", ErrEntradaInvalida)
	}

	if req.StatusContrato == "" {
		req.StatusContrato = models.StatusAtivo
	}

	if !models.IsStatusContratoValido(req.StatusContrato) {
		return nil, fmt.Errorf("%w: %s", ErrStatusInvalido, req.StatusContrato)
	}

	if req.IDContrato == "" {
		req.IDContrato = generateContractID()
	}

	user := &models.User{
		Nome:           req.Nome,
		CPF:            req.CPF,
		DataNascimento: req.DataNascimento,
		Email:          req.Email,
		SenhaHash:      req.SenhaHash,
		StatusContrato: req.StatusContrato,
		IDContrato:     req.IDContrato,
		Ativo:          req.StatusContrato == models.StatusAtivo,
	}

	err := s.repos.DB.QueryRowContext(ctx, `
		INSERT INTO users (nome, cpf, data_nascimento, email, senha_hash, status_contrato, id_contrato, ativo)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, criado_em, atualizado_em
	`, user.Nome, user.CPF, user.DataNascimento, user.Email, user.SenhaHash, user.StatusContrato, user.IDContrato, user.Ativo).
		Scan(&user.ID, &user.CriadoEm, &user.AtualizadoEm)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			s.logger.Error("tentativa de criar usuario duplicado", "email", user.Email, "cpf", user.CPF)
			return nil, ErrUsuarioDuplicado
		}

		s.logger.Error("erro ao criar usuario", "err", err)
		return nil, err
	}

	return user, nil
}

func generateContractID() string {
	now := time.Now().UTC()
	return fmt.Sprintf("CTR-%s-%04d", now.Format("200601"), now.Nanosecond()%10000)
}
