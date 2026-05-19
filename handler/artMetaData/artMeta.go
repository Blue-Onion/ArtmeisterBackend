package artmetadata

import (
	"net/http"

	"github.com/Blue-Onion/ArtmeisterBackend/handler"
	"github.com/Blue-Onion/ArtmeisterBackend/internal/database"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
)

type Handler struct {
	Repo database.ArtMetaDataRepository
}

func (h *Handler) HandleGetArtComments(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	artId, err := uuid.Parse(id)
	if err != nil {
		handler.RespondWithError(w, 400, err.Error())
		return
	}
	comments, err := h.Repo.GetArtCommentsByArtID(r.Context(), artId)
	if err != nil {
		handler.RespondWithError(w, 400, err.Error())
		return
	}
	handler.RespondWithJson(w, 200, comments)
}
func (h *Handler) HandleGetArtCommentsCount(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	artId, err := uuid.Parse(id)
	if err != nil {
		handler.RespondWithError(w, 400, err.Error())
		return
	}
	commentsCount, err := h.Repo.GetArtCommentsCount(r.Context(), artId)
	if err != nil {
		handler.RespondWithError(w, 400, err.Error())
		return
	}
	handler.RespondWithJson(w, 200, commentsCount)
}
func (h *Handler) HandleGetArtLikeCount(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	artId, err := uuid.Parse(id)
	if err != nil {
		handler.RespondWithError(w, 400, err.Error())
		return
	}
	likeCount, err := h.Repo.GetArtLikesCount(r.Context(), artId)
	if err != nil {
		handler.RespondWithError(w, 400, err.Error())
		return
	}
	handler.RespondWithJson(w, 200, likeCount)
}
