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

func (h Handler) MiddlewareAuth(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		{
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
				handler.RespondWithError(w, http.StatusBadRequest, "Invalid user id format")
				return
			}

			user, err := h.Repo.GetUser(r.Context(), id)
			if err != nil {
				handler.RespondWithError(w, http.StatusUnauthorized, "Unauthorized: user not found")
				return
			}
			ctx := context.WithValue(r.Context(), "user", user)
			next.ServeHTTP(w, r.WithContext(ctx))

		}

	}
}
