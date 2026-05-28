package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"rest-api-yt/internal/models"
)

func main() {
	req := models.CreatUserRequest{
		Name:  "Richardt",
		Email: "rj.justke@gmail.com",
	}
	b, err := json.Marshal(req)
	if err != nil {
		panic(err)
	}

	http.Post("http://localhost:8080/users", "application/json", bytes.NewReader(b))
}
