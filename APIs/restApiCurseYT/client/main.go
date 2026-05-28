package main

import (
	"bytes"
	"encoding/json"
	"fmt"
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

	resp, err := http.Post("http://localhost:8080/users", "application/json", bytes.NewReader(b))
	if err != nil {
		panic(err)
	}
	if resp.StatusCode != http.StatusCreated {
		panic("error to creat user")
	}

	var responseApi models.CreatUserResponse
	if err := json.NewDecoder(resp.Body).Decode(&responseApi); err != nil {
		panic(err)
	}
	fmt.Println("new user created", responseApi)

}
