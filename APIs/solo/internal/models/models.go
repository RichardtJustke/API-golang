package models

import (
	"database/sql"
	"strings"
	"time"
)

const (
	StatusAtivo       = "ativo"
	StatusSuspenso    = "suspenso"
	StatusCancelado   = "cancelado"
	StatusInutilizado = "inutilizado"
	StatusProduzido   = "produzido"
	StatusExcluido    = "excluido"
)

var AllowedStatusContrato = map[string]struct{}{
	StatusAtivo:       {},
	StatusSuspenso:    {},
	StatusCancelado:   {},
	StatusInutilizado: {},
	StatusProduzido:   {},
	StatusExcluido:    {},
}

func NormalizeStatusContrato(status string) string {
	return strings.ToLower(strings.TrimSpace(status))
}

func IsStatusContratoValido(status string) bool {
	_, ok := AllowedStatusContrato[NormalizeStatusContrato(status)]
	return ok
}

type User struct {
	ID             int64        `json:"id" db:"id"`
	Nome           string       `json:"nome" db:"nome"`
	CPF            string       `json:"cpf" db:"cpf"`
	DataNascimento time.Time    `json:"data_nascimento" db:"data_nascimento"`
	Email          string       `json:"email" db:"email"`
	SenhaHash      string       `json:"senha_hash" db:"senha_hash"`
	StatusContrato string       `json:"status_contrato" db:"status_contrato"`
	IDContrato     string       `json:"id_contrato" db:"id_contrato"`
	Ativo          bool         `json:"ativo" db:"ativo"`
	UltimoLogin    sql.NullTime `json:"ultimo_login" db:"ultimo_login"`
	CriadoEm       time.Time    `json:"criado_em" db:"criado_em"`
	AtualizadoEm   time.Time    `json:"atualizado_em" db:"atualizado_em"`
}

type CreateUserRequest struct {
	Nome           string    `json:"nome"`
	CPF            string    `json:"cpf"`
	DataNascimento time.Time `json:"data_nascimento"`
	Email          string    `json:"email"`
	SenhaHash      string    `json:"senha_hash"`
	StatusContrato string    `json:"status_contrato"`
	IDContrato     string    `json:"id_contrato"`
}

type SearchUserParams struct {
	ID             *int64
	Nome           string
	CPF            string
	Email          string
	IDContrato     string
	StatusContrato string
	Ativo          *bool
	Limit          int
}

type CancelContractRequest struct {
	IDContrato string `json:"id_contrato"`
}
