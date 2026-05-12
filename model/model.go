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
type CreateUser struct{
	Name      string
	Email     string
	Password  string
}
type AutheticateUser struct{
	Name      string
	Email     string
	Password  string
}
