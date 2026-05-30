package services

import (
	"context"
	"fmt"
	"strings"

	"solo/internal/models"
)

func (s *Services) FindUsers(ctx context.Context, params models.SearchUserParams) ([]models.User, error) {
	query := strings.Builder{}
	query.WriteString(`
		SELECT id, nome, cpf, data_nascimento, email, senha_hash, status_contrato, id_contrato, ativo, ultimo_login, criado_em, atualizado_em
		FROM users
	`)

	args := make([]any, 0, 8)
	where := make([]string, 0, 8)
	addFilter := func(condition string, value any) {
		args = append(args, value)
		where = append(where, fmt.Sprintf(condition, len(args)))
	}

	if params.ID != nil {
		addFilter("id = $%d", *params.ID)
	}
	if nome := strings.TrimSpace(params.Nome); nome != "" {
		addFilter("nome ILIKE $%d", "%"+nome+"%")
	}
	if cpf := strings.TrimSpace(params.CPF); cpf != "" {
		addFilter("cpf = $%d", cpf)
	}
	if email := strings.TrimSpace(params.Email); email != "" {
		addFilter("email ILIKE $%d", "%"+email+"%")
	}
	if idContrato := strings.TrimSpace(params.IDContrato); idContrato != "" {
		addFilter("id_contrato = $%d", idContrato)
	}
	if status := models.NormalizeStatusContrato(params.StatusContrato); status != "" {
		if !models.IsStatusContratoValido(status) {
			return nil, fmt.Errorf("%w: %s", ErrStatusInvalido, status)
		}
		addFilter("status_contrato = $%d", status)
	}
	if params.Ativo != nil {
		addFilter("ativo = $%d", *params.Ativo)
	}

	if len(where) > 0 {
		query.WriteString(" WHERE ")
		query.WriteString(strings.Join(where, " AND "))
	}

	limit := params.Limit
	if limit <= 0 || limit > 200 {
		limit = 50
	}
	args = append(args, limit)
	query.WriteString(fmt.Sprintf(" ORDER BY id DESC LIMIT $%d", len(args)))

	rows, err := s.repos.DB.QueryContext(ctx, query.String(), args...)
	if err != nil {
		s.logger.Error("erro ao buscar usuarios", "err", err)
		return nil, err
	}
	defer rows.Close()

	users := make([]models.User, 0)
	for rows.Next() {
		var user models.User
		if err := rows.Scan(
			&user.ID,
			&user.Nome,
			&user.CPF,
			&user.DataNascimento,
			&user.Email,
			&user.SenhaHash,
			&user.StatusContrato,
			&user.IDContrato,
			&user.Ativo,
			&user.UltimoLogin,
			&user.CriadoEm,
			&user.AtualizadoEm,
		); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}
