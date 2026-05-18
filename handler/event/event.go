package event

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Blue-Onion/ArtmeisterBackend/handler"
	"github.com/Blue-Onion/ArtmeisterBackend/internal/database"
	"github.com/Blue-Onion/ArtmeisterBackend/utlis"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
)

type Handler struct {
	Repo database.EventRepository
}

func (h *Handler) HandleCreateEvent(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(20 << 20)
	if err != nil {
		handler.RespondWithError(w, http.StatusBadRequest, "Failed to parse form data")
		return
	}
	id := uuid.New()
	path := fmt.Sprintf("uploads/event/%s", id.String())
	name := r.FormValue("name")
	if len(name) < 3 {
		handler.RespondWithError(w, http.StatusBadRequest, "Too Short Name")
		return
	}
	form_date := r.FormValue("date")

	eventDate, err := time.Parse("2006-01-02", form_date)

	if err != nil {
		handler.RespondWithError(w, 400, "invalid date format (expected YYYY-MM-DD)")
		return
	}

	desc := r.FormValue("description")
	venue := r.FormValue("venue")

	status := r.FormValue("status")
	mode := database.ModeOfConduct(status)
	if mode == "" {
		handler.RespondWithError(w, http.StatusBadRequest, "Unknown Mode")
		return
	}
	hasUpdate := false
	params := database.CreateEventParams{
		ID:          id,
		Name:        name,
		Description: utlis.ToNilStr(&desc),
		Venue:       utlis.ToNilStr(&venue),
		Status:      mode,
		EventDate:   eventDate,
	}
	userfile, _, err := r.FormFile("image")
	if err == nil && userfile != nil {
		defer userfile.Close()
		userImageFilePath, saveErr := utlis.SaveLocal(userfile, "event-logo", path)
		if saveErr != nil {
			handler.RespondWithError(w, http.StatusInternalServerError, "Failed to save Event image")
			return
		}
		params.Image = utlis.ToNilStr(&userImageFilePath)
		hasUpdate = true
	}

	bannerFile, _, err := r.FormFile("banner_image")
	if err == nil && bannerFile != nil {
		defer bannerFile.Close()
		bannerImageFilePath, saveErr := utlis.SaveLocal(bannerFile, "banner_image", path)
		if saveErr != nil {
			handler.RespondWithError(w, http.StatusInternalServerError, "Failed to save banner image")
			return
		}
		params.BannerImage = utlis.ToNilStr(&bannerImageFilePath)
		hasUpdate = true
	}

	if !hasUpdate {
		handler.RespondWithError(w, http.StatusBadRequest, "At least one image (image or banner_image) is required")
		return
	}
	res, err := h.Repo.CreateEvent(r.Context(), params)
	if err != nil {
		handler.RespondWithError(w, http.StatusInternalServerError, "Failed to update images")
		return
	}
	handler.RespondWithJson(w, http.StatusOK, res)
}
func (h *Handler) HandleDeleteEvent(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	eventId, err := uuid.Parse(id)
	if err != nil {
		handler.RespondWithError(w, 400, "Invalid Id")
		return
	}
	err = h.Repo.DeleteEvent(r.Context(), eventId)
	if err != nil {
		handler.RespondWithError(w, 400, err.Error())
		return
	}
	handler.RespondWithJson(w, 200, "ok")
}
func (h *Handler) HandleGetEventById(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	eventId, err := uuid.Parse(id)
	if err != nil {
		handler.RespondWithError(w, 400, "Invalid Id")
		return
	}
	res, err := h.Repo.GetEventByID(r.Context(), eventId)
	if err != nil {
		handler.RespondWithError(w, 400, err.Error())
		return
	}
	handler.RespondWithJson(w, 200, res)
}
func (h *Handler) HandleGetAllEvent(w http.ResponseWriter, r *http.Request) {
	res, err := h.Repo.ListEvents(r.Context())
	if err != nil {
		handler.RespondWithError(w, 400, err.Error())
		return
	}
	handler.RespondWithJson(w, 200, res)
}
