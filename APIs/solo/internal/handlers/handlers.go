package handlers

import (
	"fmt"
	"log/slog"
	"net/http"

	"solo/internal/services"
)

type Handlers struct {
	services *services.Services
	mux      *http.ServeMux
}

func New(services *services.Services) *Handlers {
	h := &Handlers{
		services: services,
		mux:      http.NewServeMux(),
	}

	h.registerUserEndpoints()

	return h
}

func (h *Handlers) Listen(port int) error {
	slog.Info("api iniciada", "port", port)
	return http.ListenAndServe(fmt.Sprintf(":%d", port), h.mux)
}
