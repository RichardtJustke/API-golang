package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"rest-api-yt/internal/models"
	"rest-api-yt/internal/usecases"
)

func (h Handlers) registerUserEndPoints() {
	http.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
 		switch r.Method {
 		case http.MethodGet:
 			h.getAllUsers(w, r)
 		case http.MethodPost:
 			h.addUser(w, r)
 		default:
 			w.WriteHeader(http.StatusMethodNotAllowed)
 		}
 	})
}

func (h Handlers) getAllUsers(w http.ResponseWriter, r *http.Request) {
	users := h.useCases.GetAll()
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(users)
}

func (h Handlers) addUser(w http.ResponseWriter, r *http.Request) {
	var req models.CreatUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.ErrorResponse{Reason: err.Error()})

		return

	}

	id, err := h.useCases.Add(req)
	if err != nil {
		if errors.Is(err, usecases.ErrUserAlreadyExists) {
			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode(models.ErrorResponse{Reason: err.Error()})

			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.ErrorResponse{Reason: err.Error()})

		return

	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(models.CreatUserResponse{NewUserID: id})

}
