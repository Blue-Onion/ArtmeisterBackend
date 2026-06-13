package event

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Blue-Onion/ArtmeisterBackend/handler"
	"github.com/Blue-Onion/ArtmeisterBackend/handler/logger"
	"github.com/Blue-Onion/ArtmeisterBackend/internal/database"
	"github.com/Blue-Onion/ArtmeisterBackend/middleware"
	"github.com/Blue-Onion/ArtmeisterBackend/utlis"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
)

type EventHandler struct {
	Repo database.EventRepository
}
type EventAttendeeHandler struct {
	Repo database.EventAttendeesRepository
}

func (h *EventHandler) HandleCreateEvent(w http.ResponseWriter, r *http.Request) {
	log, _ := logger.GetLogger()
	err := r.ParseMultipartForm(20 << 20)

	if err != nil {
		if log != nil {
			log.Error(fmt.Sprintf("HandleCreateEvent: failed to parse form data: %v", err))
		}
		handler.RespondWithError(w, http.StatusBadRequest, "Failed to parse form data")
		return
	}
	id := uuid.New()
	name := r.FormValue("name")
	if len(name) < 3 {
		if log != nil {
			log.Error(fmt.Sprintf("HandleCreateEvent: name too short: '%s'", name))
		}
		handler.RespondWithError(w, http.StatusBadRequest, "Too Short Name")
		return
	}
	form_date := r.FormValue("date")
	fmt.Println(form_date)
	eventDate, err := time.Parse("2006-01-02", form_date)

	if err != nil {
		if log != nil {
			log.Error(fmt.Sprintf("HandleCreateEvent: invalid date format '%s': %v", form_date, err))
		}
		handler.RespondWithError(w, 400, "invalid date format (expected YYYY-MM-DD)")
		return
	}

	desc := r.FormValue("description")
	venue := r.FormValue("venue")
	LogoUrl := r.FormValue("LogoUrl")
	bannerUrl := r.FormValue("bannerUrl")

	status := r.FormValue("status")
	mode := database.ModeOfConduct(status)
	if mode == "" {
		if log != nil {
			log.Error(fmt.Sprintf("HandleCreateEvent: unknown mode status '%s'", status))
		}
		handler.RespondWithError(w, http.StatusBadRequest, "Unknown Mode")
		return
	}
	params := database.CreateEventParams{
		ID:          id,
		Name:        name,
		Description: utlis.ToNilStr(&desc),
		Venue:       utlis.ToNilStr(&venue),
		Status:      mode,
		Image:       utlis.ToNilStr(&LogoUrl),
		BannerImage: utlis.ToNilStr(&bannerUrl),
		EventDate:   eventDate,
	}
	res, err := h.Repo.CreateEvent(r.Context(), params)
	if err != nil {
		if log != nil {
			log.Error(fmt.Sprintf("HandleCreateEvent: failed to create event %s: %v", id, err))
		}
		handler.RespondWithError(w, http.StatusInternalServerError, "Failed to update images")
		return
	}
	if log != nil {
		log.Info(fmt.Sprintf("HandleCreateEvent: event %s created", id))
	}
	handler.RespondWithJson(w, http.StatusOK, res)
}
func (h *EventHandler) HandleDeleteEvent(w http.ResponseWriter, r *http.Request) {
	log, _ := logger.GetLogger()
	id := chi.URLParam(r, "id")
	eventId, err := uuid.Parse(id)
	if err != nil {
		if log != nil {
			log.Error(fmt.Sprintf("HandleDeleteEvent: invalid ID format '%s': %v", id, err))
		}
		handler.RespondWithError(w, http.StatusBadRequest, "Invalid Id")
		return
	}
	_, err = h.Repo.DeleteEvent(r.Context(), eventId)
	if err != nil {
		if utlis.IsNotFound(err) {
			handler.RespondWithError(w, http.StatusNotFound, "Event not found")
			return
		}
		if log != nil {
			log.Error(fmt.Sprintf("HandleDeleteEvent: failed to delete event %s: %v", eventId, err))
		}
		handler.RespondWithError(w, http.StatusInternalServerError, "Failed to delete event")
		return
	}
	path := fmt.Sprintf("%s", id)
	err = utlis.DeleteLocal(path)
	if err != nil {
		if log != nil {
			log.Error(fmt.Sprintf("HandleDeleteEvent: failed to delete local files for event %s: %v", id, err))
		}
	}
	if log != nil {
		log.Info(fmt.Sprintf("HandleDeleteEvent: event %s deleted", eventId))
	}
	handler.RespondWithJson(w, http.StatusOK, "ok")
}
func (h *EventHandler) HandleGetEventById(w http.ResponseWriter, r *http.Request) {
	log, _ := logger.GetLogger()
	id := chi.URLParam(r, "id")
	eventId, err := uuid.Parse(id)
	if err != nil {
		if log != nil {
			log.Error(fmt.Sprintf("HandleGetEventById: invalid ID format '%s': %v", id, err))
		}
		handler.RespondWithError(w, http.StatusBadRequest, "Invalid Id")
		return
	}
	res, err := h.Repo.GetEventByID(r.Context(), eventId)
	if err != nil {
		if utlis.IsNotFound(err) {
			handler.RespondWithError(w, http.StatusNotFound, "Event not found")
			return
		}
		if log != nil {
			log.Error(fmt.Sprintf("HandleGetEventById: failed to get event %s: %v", eventId, err))
		}
		handler.RespondWithError(w, http.StatusInternalServerError, "Failed to get event")
		return
	}
	if log != nil {
		log.Info(fmt.Sprintf("HandleGetEventById: retrieved event %s successfully", eventId))
	}
	handler.RespondWithJson(w, http.StatusOK, res)
}
func (h *EventHandler) HandleGetAllEvent(w http.ResponseWriter, r *http.Request) {
	log, _ := logger.GetLogger()
	res, err := h.Repo.ListEvents(r.Context())
	if err != nil {
		if log != nil {
			log.Error(fmt.Sprintf("HandleGetAllEvent: failed to list events: %v", err))
		}
		handler.RespondWithError(w, http.StatusInternalServerError, "Failed to list events")
		return
	}
	if log != nil {
		log.Info("HandleGetAllEvent: retrieved all events successfully")
	}
	handler.RespondWithJson(w, http.StatusOK, res)
}
func (h *EventHandler) HandleUpdateEvent(w http.ResponseWriter, r *http.Request) {
	log, _ := logger.GetLogger()
	err := r.ParseMultipartForm(20 << 20)
	if err != nil {
		if log != nil {
			log.Error(fmt.Sprintf("HandleUpdateEvent: failed to parse form data: %v", err))
		}
		handler.RespondWithError(w, http.StatusBadRequest, "Failed to parse form data")
		return
	}
	eventId := chi.URLParam(r, "id")
	id, err := uuid.Parse(eventId)
	if err != nil {
		if log != nil {
			log.Error(fmt.Sprintf("HandleUpdateEvent: invalid ID format '%s': %v", eventId, err))
		}
		handler.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	path := fmt.Sprintf("uploads/event/%s", id.String())
	name := r.FormValue("name")
	if len(name) < 3 {
		if log != nil {
			log.Error(fmt.Sprintf("HandleUpdateEvent: name too short for event %s: '%s'", id, name))
		}
		handler.RespondWithError(w, http.StatusBadRequest, "Too Short Name")
		return
	}
	form_date := r.FormValue("date")

	eventDate, err := time.Parse("2006-01-02", form_date)

	if err != nil {
		if log != nil {
			log.Error(fmt.Sprintf("HandleUpdateEvent: invalid date format '%s' for event %s: %v", form_date, id, err))
		}
		handler.RespondWithError(w, 400, "invalid date format (expected YYYY-MM-DD)")
		return
	}

	desc := r.FormValue("description")
	venue := r.FormValue("venue")

	status := r.FormValue("status")
	mode := database.ModeOfConduct(status)
	if mode == "" {
		if log != nil {
			log.Error(fmt.Sprintf("HandleUpdateEvent: unknown mode status '%s' for event %s", status, id))
		}
		handler.RespondWithError(w, http.StatusBadRequest, "Unknown Mode")
		return
	}
	hasUpdate := false
	params := database.UpdateEventParams{
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
			if log != nil {
				log.Error(fmt.Sprintf("HandleUpdateEvent: failed to save event logo for event %s: %v", id, saveErr))
			}
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
			if log != nil {
				log.Error(fmt.Sprintf("HandleUpdateEvent: failed to save banner image for event %s: %v", id, saveErr))
			}
			handler.RespondWithError(w, http.StatusInternalServerError, "Failed to save banner image")
			return
		}
		params.BannerImage = utlis.ToNilStr(&bannerImageFilePath)
		hasUpdate = true
	}

	if !hasUpdate {
		if log != nil {
			log.Error(fmt.Sprintf("HandleUpdateEvent: at least one image (image or banner_image) is required for event %s", id))
		}
		handler.RespondWithError(w, http.StatusBadRequest, "At least one image (image or banner_image) is required")
		return
	}
	res, err := h.Repo.UpdateEvent(r.Context(), params)
	if err != nil {
		if utlis.IsNotFound(err) {
			handler.RespondWithError(w, http.StatusNotFound, "Event not found")
			return
		}
		if log != nil {
			log.Error(fmt.Sprintf("HandleUpdateEvent: failed to update event %s: %v", id, err))
		}
		handler.RespondWithError(w, http.StatusInternalServerError, "Failed to update event")
		return
	}
	if log != nil {
		log.Info(fmt.Sprintf("HandleUpdateEvent: event %s updated", id))
	}
	handler.RespondWithJson(w, http.StatusOK, res)
}
func (h *EventAttendeeHandler) HandleJoinEvent(w http.ResponseWriter, r *http.Request) {
	log, _ := logger.GetLogger()
	user, ok := middleware.GetUser(r.Context())
	if !ok {
		if log != nil {
			log.Error("HandleJoinEvent: unauthenticated request")
		}
		handler.RespondWithError(w, http.StatusUnauthorized, "Not Authorized")
		return
	}
	userId := user.ID
	id := chi.URLParam(r, "id")
	event_id, err := uuid.Parse(id)
	if err != nil {
		if log != nil {
			log.Error(fmt.Sprintf("HandleJoinEvent: invalid event ID format '%s' for user %s: %v", id, userId, err))
		}
		handler.RespondWithError(w, 400, err.Error())
		return
	}
	param := database.EnrollUserToEventParams{
		ID:      uuid.New(),
		EventID: event_id,
		UserID:  userId,
	}
	res, err := h.Repo.EnrollUserToEvent(r.Context(), param)
	if err != nil {
		if log != nil {
			log.Error(fmt.Sprintf("HandleJoinEvent: failed to enroll user %s to event %s: %v", userId, event_id, err))
		}
		handler.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	if log != nil {
		log.Info(fmt.Sprintf("HandleJoinEvent: user %s joined event %s", userId, event_id))
	}
	handler.RespondWithJson(w, http.StatusOK, res)
}
func (h *EventAttendeeHandler) HandleDeleteEventAttendee(w http.ResponseWriter, r *http.Request) {
	log, _ := logger.GetLogger()
	user_id := r.URL.Query().Get("user_id")
	userId, err := uuid.Parse(user_id)
	if err != nil {
		if log != nil {
			log.Error(fmt.Sprintf("HandleDeleteEventAttendee: invalid user_id format '%s': %v", user_id, err))
		}
		handler.RespondWithError(w, 400, err.Error())
		return
	}
	id := chi.URLParam(r, "id")
	event_id, err := uuid.Parse(id)
	if err != nil {
		if log != nil {
			log.Error(fmt.Sprintf("HandleDeleteEventAttendee: invalid event ID format '%s' for user %s: %v", id, userId, err))
		}
		handler.RespondWithError(w, 400, err.Error())
		return
	}
	param := database.RemoveUserFromEventParams{
		EventID: event_id,
		UserID:  userId,
	}
	_, err = h.Repo.RemoveUserFromEvent(r.Context(), param)
	if err != nil {
		if utlis.IsNotFound(err) {
			handler.RespondWithError(w, http.StatusNotFound, "Attendance record not found")
			return
		}
		if log != nil {
			log.Error(fmt.Sprintf("HandleDeleteEventAttendee: failed to remove user %s from event %s: %v", userId, event_id, err))
		}
		handler.RespondWithError(w, http.StatusInternalServerError, "Failed to remove user from event")
		return
	}
	if log != nil {
		log.Info(fmt.Sprintf("HandleDeleteEventAttendee: user %s removed from event %s", userId, event_id))
	}
	handler.RespondWithJson(w, http.StatusOK, "ok")
}
func (h *EventAttendeeHandler) HandleAllEventAttendee(w http.ResponseWriter, r *http.Request) {
	log, _ := logger.GetLogger()
	id := chi.URLParam(r, "id")
	event_id, err := uuid.Parse(id)
	if err != nil {
		if log != nil {
			log.Error(fmt.Sprintf("HandleAllEventAttendee: invalid event ID format '%s': %v", id, err))
		}
		handler.RespondWithError(w, 400, err.Error())
		return
	}
	res, err := h.Repo.ListEventAttendees(r.Context(), event_id)

	if err != nil {
		if log != nil {
			log.Error(fmt.Sprintf("HandleAllEventAttendee: failed to list attendees for event %s: %v", event_id, err))
		}
		handler.RespondWithError(w, http.StatusInternalServerError, "Failed to list event attendees")
		return
	}
	if log != nil {
		log.Info(fmt.Sprintf("HandleAllEventAttendee: retrieved attendees for event %s successfully", event_id))
	}
	handler.RespondWithJson(w, http.StatusOK, res)
}
