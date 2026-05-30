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
	blockedURLs := parseCSVEnv("URL_BLOCK", "")

	allowAnyOrigin := len(allowedOrigins) == 1 && allowedOrigins[0] == "*"
	allowedOriginsSet := make(map[string]struct{}, len(allowedOrigins))
	for _, origin := range allowedOrigins {
		allowedOriginsSet[origin] = struct{}{}
	}

	blockedSet := make(map[string]struct{}, len(blockedURLs))
	for _, url := range blockedURLs {
		blockedSet[strings.ToLower(strings.TrimSpace(url))] = struct{}{}
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := strings.TrimSpace(r.Header.Get("Origin"))

		if isBlockedURL(origin, blockedSet) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusForbidden)
			_, _ = w.Write([]byte(`{"error":"origem bloqueada por URL_BLOCK"}`))
			return
		}

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

func isBlockedURL(origin string, blockedSet map[string]struct{}) bool {
	origin = strings.ToLower(strings.TrimSpace(origin))
	if origin == "" || len(blockedSet) == 0 {
		return false
	}

	if _, ok := blockedSet[origin]; ok {
		return true
	}

	for blocked := range blockedSet {
		if blocked != "" && strings.HasPrefix(origin, blocked) {
			return true
		}
	}

	return false
}
