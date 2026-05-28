package models

import "github.com/google/uuid"

type User struct {
	ID    uuid.UUID
	Name  string
	Email string
}

type CreatUserRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type CreatUserResponse struct {
	NewUserID uuid.UUID `json:"newUserId"`
}
