package art

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/Blue-Onion/ArtmeisterBackend/handler"
	"github.com/Blue-Onion/ArtmeisterBackend/internal/database"
	"github.com/Blue-Onion/ArtmeisterBackend/middleware"
	"github.com/Blue-Onion/ArtmeisterBackend/utlis"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
)

type Handler struct {
	Repo database.ArtRepository
}

func (h *Handler) HandleArtCreation(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.GetUser(r.Context())
	if !ok {
		handler.RespondWithError(w, http.StatusUnauthorized, "Authentication required")
		return
	}

	err := r.ParseMultipartForm(20 << 20)
	if err != nil {
		handler.RespondWithError(w, http.StatusBadRequest, "Failed to parse form data")
		return
	}
	file, fileHeader, err := r.FormFile("image")
	if fileHeader.Size > 5<<20 {
		handler.RespondWithError(w, http.StatusRequestEntityTooLarge, "File too large")
	}
	if err != nil {
		handler.RespondWithError(w, http.StatusBadRequest, "Image file is required")
		return
	}
	defer file.Close()
	name := r.FormValue("name")
	if len(name) < 3 {
		handler.RespondWithError(w, http.StatusBadRequest, "Name is too Short")
		return
	}
	desc := r.FormValue("description")
	tags := r.MultipartForm.Value["tags"]
	path := fmt.Sprintf("uploads/%s/art", user.ID.String())
	id := uuid.New()
	url, err := utlis.SaveLocal(file, id.String(), path)
	if err != nil {
		handler.RespondWithError(w, http.StatusInternalServerError, "Failed to save image")
		return
	}
	params := database.CreateArtParams{
		ID:          id,
		Name:        name,
		Description: utlis.ToNilStr(&desc),
		Tags:        tags,
		UserID:      user.ID,
		Image:       url,
	}
	art, err := h.Repo.CreateArt(r.Context(), params)
	if err != nil {
		handler.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	handler.RespondWithJson(w, 200, art)
}
func (h *Handler) HandleArtDeletion(w http.ResponseWriter, r *http.Request) {

	user, ok := middleware.GetUser(r.Context())
	if !ok {
		handler.RespondWithError(w, http.StatusUnauthorized, "Authentication required")
		return
	}
	id := chi.URLParam(r, "id")
	if id == "" {
		handler.RespondWithError(w, 400, "Sahi Id bhej Bhawde")
		return
	}
	artId, err := uuid.Parse(id)
	if err != nil {
		handler.RespondWithError(w, 400, err.Error())
		return
	}
	param := database.DeleteArtParams{
		ID:     artId,
		UserID: user.ID,
	}
	err = h.Repo.DeleteArt(r.Context(), param)
	if err != nil {

		if errors.Is(err, sql.ErrNoRows) {

			handler.RespondWithError(w, 404, "Art not found")

			return

		}

		handler.RespondWithError(w, 500, "Something went wrong")
		return
	}
	handler.RespondWithJson(w, 200, "Art Work Deleted")

}
