package main

import (
	"log"
	"os"
	"strconv"

	"solo/internal/handlers"
	"solo/internal/repositories"
	"solo/internal/services"
)

func main() {
	port := 8080
	if rawPort := os.Getenv("PORT"); rawPort != "" {
		if parsedPort, err := strconv.Atoi(rawPort); err == nil {
			port = parsedPort
		}
	}

	repos, err := repositories.New()
	if err != nil {
		log.Fatal(err)
	}
	defer repos.DB.Close()

	svc := services.New(repos)
	h := handlers.New(svc)

	if err := h.Listen(port); err != nil {
		log.Fatal(err)
	}
}
