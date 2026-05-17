package art

import (
	"net/http"

	"github.com/Blue-Onion/ArtmeisterBackend/middleware"
	"github.com/go-chi/chi"
)

// ArtRouter defines routes for the art package.
func ArtRouter(artHandler *Handler, middlewareHandler *middleware.Handler) *chi.Mux {
	r := chi.NewRouter()
	auth := middlewareHandler.MiddlewareAuth

	// Public routes
	r.Get("/user/{user_id}", artHandler.HandleGetArts)
	r.Get("/{id}", artHandler.HandleGetArtById)

	// Protected routes (require user authentication)
	r.Post("/", auth(http.HandlerFunc(artHandler.HandleArtCreation)))
	r.Delete("/{id}", auth(http.HandlerFunc(artHandler.HandleArtDeletion)))
	r.Patch("/{id}", auth(http.HandlerFunc(artHandler.HandlerArtUpdation)))

	return r
}
