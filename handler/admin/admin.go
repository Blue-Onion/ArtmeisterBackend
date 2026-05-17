package admin

import (
	"encoding/json"
	"net/http"

	"github.com/Blue-Onion/ArtmeisterBackend/handler"
	"github.com/Blue-Onion/ArtmeisterBackend/internal/database"
	"github.com/Blue-Onion/ArtmeisterBackend/model"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
)

type UserHandler struct {
	Repo database.UserRepository
}
type ArtHandler struct {
	Repo database.ArtRepository
}

func (h *UserHandler) HandlerUserStatus(w http.ResponseWriter, r *http.Request) {
	userId := chi.URLParam(r, "id")
	id, err := uuid.Parse(userId)
	if err != nil {
		handler.RespondWithError(w, 400, err.Error())
		return
	}
	req := model.PatchUserStatus{}
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&req)
	if err != nil {
		handler.RespondWithError(w, 400, err.Error())
		return
	}
	if req.Role == "" && req.Status == "" {
		handler.RespondWithError(w, 400, "Ehhhh Wrong Madfaka")
		return
	}
	if req.Role == "" && req.Role != string(database.UserRoleModerator) && req.Role != string(database.UserRoleUser) {
		handler.RespondWithError(w, 400, "Ehhhh Wrong Madfaka")
		return
	}
	if req.Status == "" && req.Status != string(database.AccountStatusApproved) && req.Status != string(database.AccountStatusBanned) && req.Status != string(database.AccountStatusPending) {
		handler.RespondWithError(w, 400, "Ehhhh Wrong Madfaka")
		return
	}
	params := database.PatchUserAdminParams{
		ID:     id,
		Role:   database.UserRole(req.Role),
		Status: database.AccountStatus(req.Status),
	}
	user, err := h.Repo.PatchUserAdmin(r.Context(), params)
	handler.RespondWithJson(w, 200, user)
}
func (h *ArtHandler) HandlerArtStatus(w http.ResponseWriter, r *http.Request) {
	artId := chi.URLParam(r, "art_id")
	id, err := uuid.Parse(artId)
	if err != nil {
		handler.RespondWithError(w, 400, err.Error())
		return
	}
	status := r.URL.Query().Get("status")
	if status != string(database.ArtStatusApproved) && status != string(database.AccountStatusBanned) && status != string(database.AccountStatusPending) {

		handler.RespondWithError(w, 400, "Ehhh Wrong Status")
		return
	}
	params := database.UpdateArtStatusParams{
		ID:     id,
		Status: database.ArtStatus(status),
	}
	art, err := h.Repo.UpdateArtStatus(r.Context(), params)
	handler.RespondWithJson(w, 200, art)
}
