package user

import (
	"net/http"

	"github.com/Blue-Onion/ArtmeisterBackend/middleware"
	"github.com/go-chi/chi"
)

func UserMethod(userHandler *Handler, middlewareHandler *middleware.Handler) *chi.Mux {

	userRoute := chi.NewRouter()
	userRoute.Post("/users", userHandler.HandleCreateUser)
	userRoute.Post("/login", userHandler.HandleLogin)
	userRoute.Post("/update-img", middlewareHandler.MiddlewareAuth(http.HandlerFunc(userHandler.HandleUpdateImg)))
	userRoute.Post("/logOut", middlewareHandler.MiddlewareAuth(http.HandlerFunc(userHandler.HandleLogOut)))
	return userRoute
}
