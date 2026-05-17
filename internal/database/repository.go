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
	PatchUserAdmin(ctx context.Context, arg PatchUserAdminParams) (PatchUserAdminRow, error)
	PatchUserPassword(ctx context.Context, arg PatchUserPasswordParams) (PatchUserPasswordRow, error)
}
type ArtRepository interface {
	DeleteArt(ctx context.Context, arg DeleteArtParams) error
	GetArtByID(ctx context.Context, id uuid.UUID) (Art, error)
	GetArtByUser(ctx context.Context, userID uuid.UUID) ([]Art, error)
	ListArt(ctx context.Context) ([]Art, error)
	ListArtByTag(ctx context.Context, tags []string) ([]Art, error)
	ListArtByTags(ctx context.Context, dollar_1 []string) ([]Art, error)
	UpdateArt(ctx context.Context, arg UpdateArtParams) (Art, error)
	UpdateArtStatus(ctx context.Context, arg UpdateArtStatusParams) (Art, error)
	CreateArt(ctx context.Context, arg CreateArtParams) (Art, error)
}
