package handlers

import (
	"fmt"
	"log/slog"
	"net/http"
	"rest-api-yt/internal/usecases"
)

type Handlers struct {
	useCases *usecases.UseCases
}

func New(useCases *usecases.UseCases) *Handlers {
	return &Handlers{
		useCases: useCases,
	}
}

func (h Handlers) Listen(port int) error {
	h.registerUserEndPoints()

	slog.Info("listening on", "port", port)

	return http.ListenAndServe(
		fmt.Sprintf(":%v", port),
		nil,
	)
}
