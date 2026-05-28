package main

import (
	"rest-api-yt/internal/handlers"
	"rest-api-yt/internal/repositories"
	"rest-api-yt/internal/usecases"
)

func main() {
	repos := repositories.New()
	useCases := usecases.New(repos)
	h := handlers.New(useCases)
	h.Listen(8080)
}
