package database

import (
	"context"
	"github.com/google/uuid"
)

type UserRepository interface {
	CreateUser(ctx context.Context, arg CreateUserParams) (CreateUserRow, error)
	GetUser(ctx context.Context, id uuid.UUID) (GetUserRow, error)
	GetUserByEmail(ctx context.Context, email string) (GetUserByEmailRow, error)
	PatchUserImages(ctx context.Context, arg PatchUserImagesParams) (PatchUserImagesRow, error)
	PatchUserProfile(ctx context.Context, arg PatchUserProfileParams) (PatchUserProfileRow, error)
	PatchUserPassword(ctx context.Context, arg PatchUserPasswordParams) (PatchUserPasswordRow, error)
}
