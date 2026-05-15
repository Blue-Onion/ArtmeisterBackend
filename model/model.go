package model

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
type CreateUser struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Batch    string `json:"batch"`
}
type AuthenticateUser struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type PatchUserProfileRequest struct {
	Name  *string `json:"name"`
	Email *string `json:"email"`
	Batch *string `json:"batch"`
}
type PatchUserPassword struct {
	OldPassword string `json:"old_password"`
	Password    string `json:"password"`
}
