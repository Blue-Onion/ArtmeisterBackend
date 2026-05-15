package model

import (
	"database/sql"
	"time"

	"github.com/Blue-Onion/ArtmeisterBackend/internal/database"
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
	Name     string
	Email    string
	Password string
}
type UpdateUser struct {
	Name     string
	Email    string
	Password string
	Batch    sql.NullString
	Status   database.AccountStatus
	Role     database.UserRole
}
