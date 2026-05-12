package middleware

import (
	"context"
	"net/http"

	"github.com/Blue-Onion/RestApi-Go/handler"
	"github.com/Blue-Onion/RestApi-Go/internal/database"
	"github.com/Blue-Onion/RestApi-Go/utlis"
	"github.com/google/uuid"
)

type Handler struct {
	Repo database.UserRepository
}

func (h Handler) MiddlewareAuth(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token, err := r.Cookie("authToken")
		if err != nil {
			handler.RespondWithError(w, 401, "Error Occured in getting Token")
			return
		}
		userId, err := utlis.GetUserIdJwt(token)
		if err != nil {
			handler.RespondWithError(w, 400, "Invalid User id")
			return
		}
		id, err := uuid.Parse(userId)
		user, err := h.Repo.GetUser(r.Context(), id)
		if err != nil {
			handler.RespondWithError(w, 401, "user not found")
			return
		}
		ctx := context.WithValue(r.Context(), "user", user)

		next.ServeHTTP(w, r.WithContext(ctx))

	}
}
