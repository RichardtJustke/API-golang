package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"solo/internal/models"
	"solo/internal/services"
)

func (h *Handlers) registerUserEndpoints() {
	h.mux.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			h.createUser(w, r)
		case http.MethodGet:
			h.findUsers(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	// cancel endpoint removed; use /users/contracts/status for status changes

	h.mux.HandleFunc("/users/contracts/status", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		h.changeContractStatus(w, r)
	})
}

type createUserPayload struct {
	Nome           string `json:"nome"`
	CPF            string `json:"cpf"`
	DataNascimento string `json:"data_nascimento"`
	Email          string `json:"email"`
	SenhaHash      string `json:"senha_hash"`
	StatusContrato string `json:"status_contrato"`
	IDContrato     string `json:"id_contrato"`
}

func (h *Handlers) createUser(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var payload createUserPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "json invalido"})
		return
	}

	dataNascimento, err := parseDate(payload.DataNascimento)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "data_nascimento invalida, use YYYY-MM-DD ou RFC3339"})
		return
	}

	user, err := h.services.CreateUser(r.Context(), models.CreateUserRequest{
		Nome:           payload.Nome,
		CPF:            payload.CPF,
		DataNascimento: dataNascimento,
		Email:          payload.Email,
		SenhaHash:      payload.SenhaHash,
		StatusContrato: payload.StatusContrato,
		IDContrato:     payload.IDContrato,
	})
	if err != nil {
		h.handleServiceError(w, err)
		return
	}

	// format response dates to Brazilian format
	resp := map[string]any{
		"id":              user.ID,
		"nome":            user.Nome,
		"cpf":             user.CPF,
		"data_nascimento": formatDate(user.DataNascimento),
		"email":           user.Email,
		"senha_hash":      user.SenhaHash,
		"status_contrato": user.StatusContrato,
		"id_contrato":     user.IDContrato,
		"ativo":           user.Ativo,
		"ultimo_login":    formatNullTime(user.UltimoLogin),
		"criado_em":       formatDateTime(user.CriadoEm),
		"atualizado_em":   formatDateTime(user.AtualizadoEm),
	}

	writeJSON(w, http.StatusCreated, resp)
}

func (h *Handlers) findUsers(w http.ResponseWriter, r *http.Request) {
	params, err := buildSearchParams(r)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	users, err := h.services.FindUsers(r.Context(), params)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}

	// format dates for response
	out := make([]map[string]any, 0, len(users))
	for _, u := range users {
		out = append(out, map[string]any{
			"id":              u.ID,
			"nome":            u.Nome,
			"cpf":             u.CPF,
			"data_nascimento": formatDate(u.DataNascimento),
			"email":           u.Email,
			"senha_hash":      u.SenhaHash,
			"status_contrato": u.StatusContrato,
			"id_contrato":     u.IDContrato,
			"ativo":           u.Ativo,
			"ultimo_login":    formatNullTime(u.UltimoLogin),
			"criado_em":       formatDateTime(u.CriadoEm),
			"atualizado_em":   formatDateTime(u.AtualizadoEm),
		})
	}

	writeJSON(w, http.StatusOK, out)
}

func (h *Handlers) cancelContract(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var req models.CancelContractRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "json invalido"})
		return
	}

	if err := h.services.CancelContract(r.Context(), req); err != nil {
		h.handleServiceError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "contrato cancelado com sucesso"})
}

type changeStatusPayload struct {
	IDContrato string `json:"id_contrato"`
	Status     string `json:"status"`
}

func (h *Handlers) changeContractStatus(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var req changeStatusPayload
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "json invalido"})
		return
	}

	if err := h.services.UpdateContractStatus(r.Context(), req.IDContrato, req.Status); err != nil {
		h.handleServiceError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "status atualizado com sucesso"})
}

func buildSearchParams(r *http.Request) (models.SearchUserParams, error) {
	q := r.URL.Query()
	params := models.SearchUserParams{
		Nome:           q.Get("nome"),
		CPF:            q.Get("cpf"),
		Email:          q.Get("email"),
		IDContrato:     q.Get("id_contrato"),
		StatusContrato: q.Get("status_contrato"),
	}

	if idRaw := strings.TrimSpace(q.Get("id")); idRaw != "" {
		id, err := strconv.ParseInt(idRaw, 10, 64)
		if err != nil {
			return params, errors.New("id invalido")
		}
		params.ID = &id
	}

	if ativoRaw := strings.TrimSpace(q.Get("ativo")); ativoRaw != "" {
		ativo, err := strconv.ParseBool(ativoRaw)
		if err != nil {
			return params, errors.New("ativo invalido, use true/false")
		}
		params.Ativo = &ativo
	}

	if limitRaw := strings.TrimSpace(q.Get("limit")); limitRaw != "" {
		limit, err := strconv.Atoi(limitRaw)
		if err != nil {
			return params, errors.New("limit invalido")
		}
		params.Limit = limit
	}

	return params, nil
}

func (h *Handlers) handleServiceError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, services.ErrEntradaInvalida), errors.Is(err, services.ErrStatusInvalido):
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
	case errors.Is(err, services.ErrUsuarioDuplicado):
		writeJSON(w, http.StatusConflict, map[string]string{"error": err.Error()})
	case errors.Is(err, services.ErrUsuarioNaoEncontrado):
		writeJSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})
	default:
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "erro interno"})
	}
}

func writeJSON(w http.ResponseWriter, statusCode int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(payload)
}

func parseDate(raw string) (time.Time, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return time.Time{}, errors.New("data vazia")
	}

	if t, err := time.Parse("2006-01-02", raw); err == nil {
		return t, nil
	}

	return time.Parse(time.RFC3339, raw)
}

func formatDate(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format("02/01/2006")
}

func formatDateTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format("02/01/2006 15:04:05")
}

func formatNullTime(nt sql.NullTime) string {
	if !nt.Valid {
		return ""
	}
	return formatDateTime(nt.Time)
}
