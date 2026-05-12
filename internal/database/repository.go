package database

import (
	"context"
	"github.com/google/uuid"
)

type UserRepository interface {
	CreateUser(ctx context.Context, arg CreateUserParams) (CreateUserRow, error)
	GetUser(ctx context.Context, id uuid.UUID) (GetUserRow, error)
	GetUserByEmail(ctx context.Context, email string) (GetUserByEmailRow, error)
}
