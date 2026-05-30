package services

import (
	"context"
	"errors"
	"fmt"
	"strconv"
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
		req.IDContrato = s.generateContractID(ctx)
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
		INSERT INTO usuarios (nome, cpf, data_nascimento, email, senha_hash, status_contrato, id_contrato, ativo)
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
	// kept for backward compatibility if ever used elsewhere
	now := time.Now().UTC()
	return fmt.Sprintf("CTR-%s-%04d", now.Format("2006"), now.Nanosecond()%10000)
}

// generateContractID generates a contract id in the format CTR-YYYY-0001
// It queries the database for the current max id for the year and increments it.
func (s *Services) generateContractID(ctx context.Context) string {
	year := time.Now().UTC().Format("2006")
	pattern := fmt.Sprintf("CTR-%s-%%", year)

	var last string
	err := s.repos.DB.QueryRowContext(ctx, `
		SELECT id_contrato
		FROM usuarios
		WHERE id_contrato LIKE $1
		ORDER BY id_contrato DESC
		LIMIT 1
	`, pattern).Scan(&last)
	next := 1
	if err == nil && last != "" {
		parts := strings.Split(last, "-")
		if len(parts) == 3 {
			if n, perr := strconv.Atoi(parts[2]); perr == nil {
				next = n + 1
			}
		}
	}

	return fmt.Sprintf("CTR-%s-%04d", year, next)
}
