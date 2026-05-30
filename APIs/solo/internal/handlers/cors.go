package handlers

import (
	"net/http"
	"os"
	"strings"
)

func (h *Handlers) withCORS(next http.Handler) http.Handler {
	allowedOrigins := parseCSVEnv("CORS_ALLOWED_ORIGINS", "*")
	allowedMethods := parseCSVEnv("CORS_ALLOWED_METHODS", "GET,POST,PATCH,OPTIONS")
	allowedHeaders := parseCSVEnv("CORS_ALLOWED_HEADERS", "Content-Type,Authorization")

	allowAnyOrigin := len(allowedOrigins) == 1 && allowedOrigins[0] == "*"
	allowedOriginsSet := make(map[string]struct{}, len(allowedOrigins))
	for _, origin := range allowedOrigins {
		allowedOriginsSet[origin] = struct{}{}
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := strings.TrimSpace(r.Header.Get("Origin"))
		if allowAnyOrigin {
			w.Header().Set("Access-Control-Allow-Origin", "*")
		} else if origin != "" {
			if _, ok := allowedOriginsSet[origin]; ok {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Vary", "Origin")
			}
		}

		w.Header().Set("Access-Control-Allow-Methods", strings.Join(allowedMethods, ", "))
		w.Header().Set("Access-Control-Allow-Headers", strings.Join(allowedHeaders, ", "))
		w.Header().Set("Access-Control-Max-Age", "600")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func parseCSVEnv(key, fallback string) []string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		value = fallback
	}

	parts := strings.Split(value, ",")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}

	return result
}
