package models

import "github.com/google/uuid"

type User struct {
	ID   uuid.UUID
	Name string
	Age  int
}

type CreatUserRequest struct{
	Name string
	Age int 
}
