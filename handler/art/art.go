package art

import (
	"fmt"
	"net/http"

	"github.com/Blue-Onion/ArtmeisterBackend/handler"
	"github.com/Blue-Onion/ArtmeisterBackend/handler/logger"
	"github.com/Blue-Onion/ArtmeisterBackend/internal/database"
	"github.com/Blue-Onion/ArtmeisterBackend/middleware"
	"github.com/Blue-Onion/ArtmeisterBackend/utlis"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
)

type Handler struct {
	Repo database.ArtRepository
}
type ProfileHandler struct {
	ArtRepo  database.ArtRepository
	UserRepo database.UserRepository
}

type profile struct {
	User database.GetUserRow
	Art  []database.Art
}

func (h *Handler) HandleArtCreation(w http.ResponseWriter, r *http.Request) {
	log, _ := logger.GetLogger()
	user, ok := middleware.GetUser(r.Context())
	if !ok {
		if log != nil {
			log.Error("HandleArtCreation: unauthenticated request")
		}
		handler.RespondWithError(w, http.StatusUnauthorized, "Authentication required")
		return
	}
	name := r.FormValue("name")
	if len(name) < 3 {
		if log != nil {
			log.Error(fmt.Sprintf("HandleArtCreation: art name too short for user %s: '%s'", user.ID, name))
		}
		handler.RespondWithError(w, http.StatusBadRequest, "Name is too Short")
		return
	}

	url := r.FormValue("url")
	if len(url) < 3 {

		if log != nil {
			log.Error("HandleArtCreation: Invalid Url")
		}
		handler.RespondWithError(w, http.StatusBadRequest, "Invalid Url")
		return
	}
	desc := r.FormValue("description")
	tags := r.MultipartForm.Value["tags"]
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {

		handler.RespondWithError(w, http.StatusBadRequest, "Invalid form data")

		return

	}
	if tags == nil {
		tags = []string{}
	}
	id := uuid.New()

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
		if log != nil {
			log.Error(fmt.Sprintf("HandleArtCreation: failed to create art for user %s: %v", user.ID, err))
		}
		handler.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	if log != nil {
		log.Info(fmt.Sprintf("HandleArtCreation: art %s created by user %s", id, user.ID))
	}
	handler.RespondWithJson(w, http.StatusOK, art)
}
func (h *Handler) HandleGetArts(w http.ResponseWriter, r *http.Request) {
	log, _ := logger.GetLogger()
	userId := chi.URLParam(r, "user_id")
	if userId == "" {
		if log != nil {
			log.Error("HandleGetArts: user ID is empty")
		}
		handler.RespondWithError(w, http.StatusBadRequest, "User ID is required")
		return
	}
	fmt.Printf("userId raw = %q\n", userId)

	fmt.Printf("userId len = %d\n", len(userId))

	id, err := uuid.Parse(userId)
	if err != nil {
		if log != nil {
			log.Error(fmt.Sprintf("userId=%q err=%v", userId, err))
		}
		handler.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	fmt.Println("Got here")
	arts, err := h.Repo.GetArtByUser(r.Context(), id)
	if err != nil {
		if log != nil {
			log.Error(fmt.Sprintf("HandleGetArts: failed to get arts for user %s: %v", id, err))
		}
		handler.RespondWithError(w, http.StatusInternalServerError, "Failed to get user arts")
		return
	}
	if log != nil {
		log.Info(fmt.Sprintf("HandleGetArts: retrieved arts for user %s successfully", id))
	}
	fmt.Println("Got here 1")
	handler.RespondWithJson(w, http.StatusOK, arts)
}
func (h *Handler) HandleGetArtById(w http.ResponseWriter, r *http.Request) {
	log, _ := logger.GetLogger()
	Id := chi.URLParam(r, "id")
	if Id == "" {
		if log != nil {
			log.Error("HandleGetArtById: art ID is empty")
		}
		handler.RespondWithError(w, http.StatusBadRequest, "Art ID is required")
		return
	}
	artId, err := uuid.Parse(Id)
	if err != nil {
		if log != nil {
			log.Error(fmt.Sprintf("HandleGetArtById: invalid art ID format '%s': %v", Id, err))
		}
		handler.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	art, err := h.Repo.GetArtByID(r.Context(), artId)
	if err != nil {
		if utlis.IsNotFound(err) {
			handler.RespondWithError(w, http.StatusNotFound, "Art not found")
			return
		}
		if log != nil {
			log.Error(fmt.Sprintf("HandleGetArtById: failed to get art %s: %v", artId, err))
		}
		handler.RespondWithError(w, http.StatusInternalServerError, "Failed to get art")
		return
	}
	if log != nil {
		log.Info(fmt.Sprintf("HandleGetArtById: retrieved art %s successfully", artId))
	}
	handler.RespondWithJson(w, http.StatusOK, art)
}
func (h *Handler) HandleArtDeletion(w http.ResponseWriter, r *http.Request) {
	log, _ := logger.GetLogger()
	user, ok := middleware.GetUser(r.Context())
	if !ok {
		if log != nil {
			log.Error("HandleArtDeletion: unauthenticated request")
		}
		handler.RespondWithError(w, http.StatusUnauthorized, "Authentication required")
		return
	}
	userId := user.ID.String()
	id := chi.URLParam(r, "id")
	if id == "" {
		if log != nil {
			log.Error(fmt.Sprintf("HandleArtDeletion: art ID is empty (user %s)", user.ID))
		}
		handler.RespondWithError(w, http.StatusBadRequest, "Art ID is required")
		return
	}
	artId, err := uuid.Parse(id)
	if err != nil {
		if log != nil {
			log.Error(fmt.Sprintf("HandleArtDeletion: invalid art ID format '%s' for user %s: %v", id, user.ID, err))
		}
		handler.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	param := database.DeleteArtParams{
		ID:     artId,
		UserID: user.ID,
	}
	_, err = h.Repo.DeleteArt(r.Context(), param)
	if err != nil {
		if utlis.IsNotFound(err) {
			handler.RespondWithError(w, http.StatusNotFound, "Art not found")
			return
		}
		if log != nil {
			log.Error(fmt.Sprintf("HandleArtDeletion: failed to delete art %s: %v", artId, err))
		}
		handler.RespondWithError(w, http.StatusInternalServerError, "Failed to delete art")
		return
	}
	path := fmt.Sprintf("%s/art/%s.png", userId, id)
	err = utlis.DeleteLocal(path)
	if err != nil {
		if log != nil {
			log.Error(fmt.Sprintf("HandleArtDeletion: failed to delete local file for art %s: %v", id, err))
		}
	}
	if log != nil {
		log.Info(fmt.Sprintf("HandleArtDeletion: art %s deleted by user %s", artId, user.ID))
	}
	handler.RespondWithJson(w, http.StatusOK, "Art Work Deleted")

}
func (h *Handler) HandlerArtUpdation(w http.ResponseWriter, r *http.Request) {
	log, _ := logger.GetLogger()
	user, ok := middleware.GetUser(r.Context())
	if !ok {
		if log != nil {
			log.Error("HandlerArtUpdation: unauthenticated request")
		}
		handler.RespondWithError(w, http.StatusUnauthorized, "Not Authorized")
		return
	}
	id := chi.URLParam(r, "id")
	if id == "" {
		if log != nil {
			log.Error(fmt.Sprintf("HandlerArtUpdation: art ID is empty (user %s)", user.ID))
		}
		handler.RespondWithError(w, http.StatusBadRequest, "Art ID is required")
		return
	}
	artId, err := uuid.Parse(id)
	if err != nil {
		if log != nil {
			log.Error(fmt.Sprintf("HandlerArtUpdation: invalid art ID format '%s' for user %s: %v", id, user.ID, err))
		}
		handler.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	name := r.FormValue("name")
	if len(name) < 3 {
		if log != nil {
			log.Error(fmt.Sprintf("HandlerArtUpdation: art name too short for user %s: '%s'", user.ID, name))
		}
		handler.RespondWithError(w, http.StatusBadRequest, "Name is too Short")
		return
	}
	desc := r.FormValue("description")
	var tags []string
	if r.MultipartForm != nil {
		tags = r.MultipartForm.Value["tags"]
	}
	if tags == nil {
		tags = []string{}
	}

	params := database.UpdateArtParams{
		ID:          artId,
		UserID:      user.ID,
		Name:        name,
		Tags:        tags,
		Description: utlis.ToNilStr(&desc),
	}
	updatedWork, err := h.Repo.UpdateArt(r.Context(), params)
	if err != nil {
		if utlis.IsNotFound(err) {
			handler.RespondWithError(w, http.StatusNotFound, "Art not found")
			return
		}
		if log != nil {
			log.Error(fmt.Sprintf("HandlerArtUpdation: failed to update art %s: %v", artId, err))
		}
		handler.RespondWithError(w, http.StatusInternalServerError, "Failed to update art")
		return
	}
	if log != nil {
		log.Info(fmt.Sprintf("HandlerArtUpdation: art %s updated by user %s", artId, user.ID))
	}
	handler.RespondWithJson(w, http.StatusOK, updatedWork)
}
func (h *ProfileHandler) HandlerGetArtistProfile(w http.ResponseWriter, r *http.Request) {
	log, _ := logger.GetLogger()
	userId := chi.URLParam(r, "id")
	if userId == "" {
		if log != nil {
			log.Error("HandleGetArts: user ID is empty")
		}
		handler.RespondWithError(w, http.StatusBadRequest, "User ID is required")
		return
	}
	Id, err := uuid.Parse(userId)
	if err != nil {
		if log != nil {
			log.Error(fmt.Sprintf("HandleGetArtById: invalid art ID format '%s': %v", Id, err))
		}
		handler.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	user, err := h.UserRepo.GetUser(r.Context(), Id)
	if err != nil {
		if log != nil {
			log.Error(fmt.Sprintf("HandleGetArtById: invalid art ID format '%s': %v", Id, err))
		}
		handler.RespondWithError(w, http.StatusNotFound, err.Error())
		return
	}
	artWork, err := h.ArtRepo.GetArtByUser(r.Context(), Id)
	if err != nil {
		if log != nil {
			log.Error(fmt.Sprintf("HandleGetArtById: invalid art ID format '%s': %v", Id, err))
		}
		handler.RespondWithError(w, http.StatusNotFound, err.Error())
		return
	}
	res := profile{
		User: user,
		Art:  artWork,
	}
	fmt.Println(res)
	handler.RespondWithJson(w, 200, res)
}
func (h *Handler) HandleGetPendingArt(w http.ResponseWriter, r *http.Request) {
	log, _ := logger.GetLogger()
	arts, err := h.Repo.ListPendingArt(r.Context())
	if err != nil {
		if utlis.IsNotFound(err) {
			handler.RespondWithError(w, http.StatusNotFound, "Art not found")
			return
		}
		if log != nil {
			log.Error("HandleGetPendingArt: err Occurred")
		}
		handler.RespondWithError(w, http.StatusInternalServerError, "Failed to get art")
		return
	}
	if log != nil {
		log.Info("HandleGetPendingArt: retrieved art")
	}
	handler.RespondWithJson(w, http.StatusOK, arts)
}
