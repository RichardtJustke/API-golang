package handlers

import (
	"errors"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type authTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int64  `json:"expires_in"`
}

func (h *Handlers) registerAuthEndpoints() {
	h.mux.HandleFunc("/auth/token", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		h.createToken(w, r)
	})
}

func (h *Handlers) createToken(w http.ResponseWriter, r *http.Request) {
	secret := strings.TrimSpace(os.Getenv("AUTH_SECRET"))
	if secret == "" {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "AUTH_SECRET nao configurado"})
		return
	}

	ttl := 60 * time.Minute
	if rawTTL := strings.TrimSpace(os.Getenv("AUTH_TOKEN_TTL")); rawTTL != "" {
		parsedTTL, err := time.ParseDuration(rawTTL)
		if err == nil && parsedTTL > 0 {
			ttl = parsedTTL
		}
	}

	now := time.Now().UTC()
	expiresAt := now.Add(ttl)
	claims := jwt.MapClaims{
		"sub": "api-client",
		"iat": now.Unix(),
		"exp": expiresAt.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(secret))
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "erro ao gerar token"})
		return
	}

	writeJSON(w, http.StatusOK, authTokenResponse{
		AccessToken: signed,
		TokenType:   "Bearer",
		ExpiresIn:   int64(ttl.Seconds()),
	})
}

func (h *Handlers) withAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions || r.URL.Path == "/auth/token" {
			next.ServeHTTP(w, r)
			return
		}

		secret := strings.TrimSpace(os.Getenv("AUTH_SECRET"))
		if secret == "" {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "AUTH_SECRET nao configurado"})
			return
		}

		authHeader := strings.TrimSpace(r.Header.Get("Authorization"))
		tokenString, err := parseBearerToken(authHeader)
		if err != nil {
			writeJSON(w, http.StatusUnauthorized, map[string]string{"error": err.Error()})
			return
		}

		_, err = jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("metodo de assinatura invalido")
			}
			return []byte(secret), nil
		})
		if err != nil {
			writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "token invalido ou expirado"})
			return
		}

		next.ServeHTTP(w, r)
	})
}

func parseBearerToken(authHeader string) (string, error) {
	if authHeader == "" {
		return "", errors.New("header Authorization obrigatorio")
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return "", errors.New("use Authorization: Bearer <token>")
	}

	token := strings.TrimSpace(parts[1])
	if token == "" {
		return "", errors.New("token vazio")
	}

	return token, nil
}
