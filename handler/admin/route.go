package admin

import (
	"net/http"

	"github.com/Blue-Onion/ArtmeisterBackend/handler/art"
	"github.com/Blue-Onion/ArtmeisterBackend/handler/user"
	"github.com/Blue-Onion/ArtmeisterBackend/middleware"
	"github.com/go-chi/chi"
)

// AdminRoute constructs and returns the subrouter for admin-only routes.
// It applies the MiddlewareAdminAuth middleware to all endpoints.
func AdminRoute(userHandler *user.Handler, artHandler *art.Handler, middlewareHandler *middleware.Handler) *chi.Mux {
	r := chi.NewRouter()

	// Apply MiddlewareAdminAuth to all routes in this subrouter.
	// Since MiddlewareAdminAuth returns http.HandlerFunc, we adapt it to http.Handler.
	adminAuth := middlewareHandler.MiddlewareAdminAuth
	r.Use(func(next http.Handler) http.Handler {
		return adminAuth(next)
	})

	// Instantiate the admin handlers with the repositories from user and art handlers.
	adminUserHandler := &UserHandler{Repo: userHandler.Repo}
	adminArtHandler := &ArtHandler{Repo: artHandler.Repo}

	// Admin action routes
	r.Patch("/users/{id}/status", adminUserHandler.HandlerUserStatus)
	r.Patch("/arts/{art_id}/status", adminArtHandler.HandlerArtStatus)

	return r
}
