package handlers

import (
	"encoding/json"
	"net/http"
	"rest-api-yt/internal/models"

	"github.com/google/uuid"
)

func (h Handlers) registerUserEndPoints() {
	http.HandleFunc("GET /users", h.getAllUsers)
	http.HandleFunc("POST /users", h.addUser)
}

func (h Handlers) getAllUsers(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode([]models.User{})
}
func (h Handlers) addUser(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(models.CreatUserResponse{NewUserID: uuid.New()})
}
