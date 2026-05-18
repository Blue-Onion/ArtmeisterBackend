package event

import (
	"net/http"

	"github.com/Blue-Onion/ArtmeisterBackend/middleware"
	"github.com/go-chi/chi"
)

func ArtRouter(md *middleware.Handler) *chi.Mux {
	r := chi.NewRouter()
	adminAuth := md.MiddlewareAdminAuth
	r.Use(
		func(h http.Handler) http.Handler {
			return adminAuth(h)
		},
	)
	return r
}
