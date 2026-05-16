package art

import (
	"github.com/Blue-Onion/ArtmeisterBackend/middleware"
	"github.com/go-chi/chi"
)

func ArtRouter(userHandler *Handler, middlewareHandler *middleware.Handler) *chi.Mux {
	r := chi.NewRouter()

	return r
}
