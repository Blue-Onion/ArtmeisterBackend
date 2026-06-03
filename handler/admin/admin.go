package admin

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Blue-Onion/ArtmeisterBackend/handler"
	"github.com/Blue-Onion/ArtmeisterBackend/handler/logger"
	"github.com/Blue-Onion/ArtmeisterBackend/internal/database"
	"github.com/Blue-Onion/ArtmeisterBackend/model"
	"github.com/Blue-Onion/ArtmeisterBackend/utlis"
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
	log, _ := logger.GetLogger()
	userId := chi.URLParam(r, "id")
	id, err := uuid.Parse(userId)
	if err != nil {
		handler.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	req := model.PatchUserStatus{}
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&req)
	if err != nil {
		handler.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	if req.Role == "" && req.Status == "" {
		handler.RespondWithError(w, http.StatusBadRequest, "Role and status cannot both be empty")
		return
	}
	if req.Role != "" && req.Role != string(database.UserRoleUser) {
		handler.RespondWithError(w, http.StatusBadRequest, "Invalid role provided")
		return
	}
	if req.Status != "" && req.Status != string(database.AccountStatusApproved) && req.Status != string(database.AccountStatusBanned) && req.Status != string(database.AccountStatusPending) {
		handler.RespondWithError(w, http.StatusBadRequest, "Invalid status provided")
		return
	}
	params := database.PatchUserAdminParams{
		ID:     id,
		Role:   database.UserRole(req.Role),
		Status: database.AccountStatus(req.Status),
	}
	user, err := h.Repo.PatchUserAdmin(r.Context(), params)
	if err != nil {
		if utlis.IsNotFound(err) {
			handler.RespondWithError(w, http.StatusNotFound, "User not found")
			return
		}
		if log != nil {
			log.Error(fmt.Sprintf("HandlerUserStatus: failed to update status for user %s: %v", id, err))
		}
		handler.RespondWithError(w, http.StatusInternalServerError, "Failed to update user status")
		return
	}
	if log != nil {
		log.Info(fmt.Sprintf("HandlerUserStatus: user %s status updated (role=%s, status=%s)", id, req.Role, req.Status))
	}
	handler.RespondWithJson(w, http.StatusOK, user)
}
func (h *ArtHandler) HandlerArtStatus(w http.ResponseWriter, r *http.Request) {
	log, _ := logger.GetLogger()
	artId := chi.URLParam(r, "art_id")
	id, err := uuid.Parse(artId)
	if err != nil {
		handler.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	status := r.URL.Query().Get("status")
	if status != string(database.ArtStatusApproved) && status != string(database.AccountStatusBanned) && status != string(database.AccountStatusPending) {

		handler.RespondWithError(w, http.StatusBadRequest, "Invalid status provided")
		return
	}
	params := database.UpdateArtStatusParams{
		ID:     id,
		Status: database.ArtStatus(status),
	}
	art, err := h.Repo.UpdateArtStatus(r.Context(), params)
	if err != nil {
		if utlis.IsNotFound(err) {
			handler.RespondWithError(w, http.StatusNotFound, "Art not found")
			return
		}
		if log != nil {
			log.Error(fmt.Sprintf("HandlerArtStatus: failed to update art %s status: %v", id, err))
		}
		handler.RespondWithError(w, http.StatusInternalServerError, "Failed to update art status")
		return
	}
	if log != nil {
		log.Info(fmt.Sprintf("HandlerArtStatus: art %s status updated to %s", id, status))
	}
	handler.RespondWithJson(w, http.StatusOK, art)
}
