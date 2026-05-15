package middleware

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/Blue-Onion/ArtmeisterBackend/handler"
	"github.com/Blue-Onion/ArtmeisterBackend/internal/database"
	"github.com/Blue-Onion/ArtmeisterBackend/utlis"
	"github.com/google/uuid"
)

type Handler struct {
	Repo database.UserRepository
}

type contextKey string

const userContextKey contextKey = "user"

type User struct {
	ID          uuid.UUID
	Name        string
	Email       string
	Batch       string
	Status      database.AccountStatus
	Role        database.UserRole
	Image       sql.NullString
	BannerImage sql.NullString
}

// GetUser extracts the authenticated user from the request context.
func GetUser(ctx context.Context) (User, bool) {
	user, ok := ctx.Value(userContextKey).(User)
	return user, ok
}

func (h Handler) MiddlewareAuth(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenCookie, err := r.Cookie("authToken")
		if err != nil {
			handler.RespondWithError(w, http.StatusUnauthorized, "Unauthorized: login required")
			return
		}

		userId, err := utlis.GetUserIdJwt(tokenCookie)
		if err != nil {
			handler.RespondWithError(w, http.StatusUnauthorized, "Unauthorized: invalid or expired token")
			return
		}

		id, err := uuid.Parse(userId)
		if err != nil {
			handler.RespondWithError(w, http.StatusBadRequest, "Invalid user ID format")
			return
		}

		dbUser, err := h.Repo.GetUser(r.Context(), id)
		if err != nil {
			handler.RespondWithError(w, http.StatusUnauthorized, "Unauthorized: user not found")
			return
		}

		user := User{
			ID:          dbUser.ID,
			Name:        dbUser.Name,
			Email:       dbUser.Email,
			Batch:       dbUser.Batch,
			Status:      dbUser.Status,
			Role:        dbUser.Role,
			Image:       dbUser.Image,
			BannerImage: dbUser.BannerImage,
		}
		ctx := context.WithValue(r.Context(), userContextKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}
