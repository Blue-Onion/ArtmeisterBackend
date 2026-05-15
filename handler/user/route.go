package user

import (
	"net/http"

	"github.com/Blue-Onion/ArtmeisterBackend/middleware"
	"github.com/go-chi/chi"
)

func UserMethod(userHandler *Handler, middlewareHandler *middleware.Handler) *chi.Mux {

	userRoute := chi.NewRouter()

	// Auth
	userRoute.Post("/auth/register", userHandler.HandleCreateUser)
	userRoute.Post("/auth/login", userHandler.HandleLogin)
	userRoute.Post("/auth/logout", middlewareHandler.MiddlewareAuth(http.HandlerFunc(userHandler.HandleLogOut)))

	// User profile
	userRoute.Patch("/users/me", middlewareHandler.MiddlewareAuth(http.HandlerFunc(userHandler.HandleUpdateUserProfile)))
	userRoute.Patch("/users/me/password", middlewareHandler.MiddlewareAuth(http.HandlerFunc(userHandler.HandlePasswordChange)))
	userRoute.Patch("/users/me/avatar", middlewareHandler.MiddlewareAuth(http.HandlerFunc(userHandler.HandleUpdateImg)))

	return userRoute
}
