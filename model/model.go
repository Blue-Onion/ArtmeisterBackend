package model

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID
	Name      string
	Email     string
	Password  string
	Createdat time.Time
	Updatedat time.Time
}
type CreateUser struct {
	Name     string
	Email    string
	Password string
	Batch    string
}
type AutheticateUser struct {
	Email    string
	Password string
}
type PatchUserProfileRequest struct {
	Name  *string `json:"name"`
	Email *string `json:"email"`
	Batch *string `json:"batch"`
}
type PatchUserPassword struct {
	Password string `json:"password"`
}
