package user

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Blue-Onion/ArtmeisterBackend/handler"
	"github.com/Blue-Onion/ArtmeisterBackend/handler/logger"
	"github.com/Blue-Onion/ArtmeisterBackend/internal/database"
	"github.com/Blue-Onion/ArtmeisterBackend/middleware"
	"github.com/Blue-Onion/ArtmeisterBackend/model"
	"github.com/Blue-Onion/ArtmeisterBackend/utlis"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/sqlc-dev/pqtype"
)

type Handler struct {
	Repo database.UserRepository
}

func (h *Handler) HandleUpdateImg(w http.ResponseWriter, r *http.Request) {
	log, _ := logger.GetLogger()
	user, ok := middleware.GetUser(r.Context())
	if !ok {
		if log != nil {
			log.Error("HandleUpdateImg: unauthenticated request")
		}
		handler.RespondWithError(w, http.StatusUnauthorized, "Authentication required")
		return
	}
	err := r.ParseMultipartForm(20 << 20)
	if err != nil {
		if log != nil {
			log.Error(fmt.Sprintf("HandleUpdateImg: failed to parse multipart form: %v", err))
		}
		handler.RespondWithError(w, http.StatusBadRequest, "Failed to parse form data")
		return
	}
	file, _, err := r.FormFile("image")
	if err != nil {
		if log != nil {
			log.Error(fmt.Sprintf("HandleUpdateImg: image file form retrieval failed: %v", err))
		}
		handler.RespondWithError(w, http.StatusBadRequest, "Image file is required")
		return
	}
	defer file.Close()
	path := fmt.Sprintf("uploads/%s", user.ID.String())
	filepath, err := utlis.SaveLocal(file, "userPhoto", path)
	if err != nil {
		if log != nil {
			log.Error(fmt.Sprintf("HandleUpdateImg: failed to save image for user %s: %v", user.ID, err))
		}
		handler.RespondWithError(w, http.StatusInternalServerError, "Failed to save image")
		return
	}
	if log != nil {
		log.Info(fmt.Sprintf("HandleUpdateImg: image updated for user %s", user.ID))
	}
	handler.RespondWithJson(w, http.StatusOK, filepath)
}

func (h *Handler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	log, _ := logger.GetLogger()
	params := model.AuthenticateUser{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {
		if log != nil {
			log.Error("HandleLogin: invalid request body")
		}
		handler.RespondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	user, err := h.Repo.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		if utlis.IsNotFound(err) {
			if log != nil {
				log.Info(fmt.Sprintf("HandleLogin: no account found for email %s", params.Email))
			}
			handler.RespondWithError(w, http.StatusNotFound, "No account found with this email")
			return
		}
		if log != nil {
			log.Error(fmt.Sprintf("HandleLogin: db error looking up email %s: %v", params.Email, err))
		}
		handler.RespondWithError(w, http.StatusInternalServerError, "Failed to look up user")
		return
	}
	isValid := utlis.CheckPassword(user.Password, params.Password)
	if !isValid {
		if log != nil {
			log.Info(fmt.Sprintf("HandleLogin: incorrect password for user %s", user.ID))
		}
		handler.RespondWithError(w, http.StatusUnauthorized, "Incorrect password")
		return
	}
	token, err := utlis.GenerateJwt(user.ID)
	if err != nil {
		if log != nil {
			log.Error(fmt.Sprintf("HandleLogin: failed to generate JWT for user %s: %v", user.ID, err))
		}
		handler.RespondWithError(w, http.StatusInternalServerError, "Failed to generate authentication token")
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "authToken",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		MaxAge:   3600 * 24,
		SameSite: http.SameSiteLaxMode,
	})
	if log != nil {
		log.Info(fmt.Sprintf("HandleLogin: user %s logged in successfully", user.ID))
	}
	handler.RespondWithJson(w, http.StatusOK, map[string]string{
		"id":    user.ID.String(),
		"name":  user.Name,
		"email": user.Email,
		"image": user.Image.String,
	})
}

func (h *Handler) HandleMe(w http.ResponseWriter, r *http.Request) {
	log, _ := logger.GetLogger()
	user, ok := middleware.GetUser(r.Context())
	if !ok {
		if log != nil {
			log.Error("HandleMe: unauthenticated request")
		}
		handler.RespondWithError(w, http.StatusUnauthorized, "Not Authenticated User")
		return
	}
	if log != nil {
		log.Info(fmt.Sprintf("HandleMe: profile retrieved for user %s", user.ID))
	}
	handler.RespondWithJson(w, 200, user)

}
func (h *Handler) HandleLogOut(w http.ResponseWriter, r *http.Request) {
	log, _ := logger.GetLogger()
	http.SetCookie(w, &http.Cookie{
		Name:     "authToken",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	})
	if log != nil {
		log.Info("HandleLogOut: user logged out")
	}
	handler.RespondWithJson(w, http.StatusOK, map[string]string{
		"message": "Logged out successfully",
	})
}

func (h *Handler) HandleCreateUser(w http.ResponseWriter, r *http.Request) {
	log, _ := logger.GetLogger()
	param := model.CreateUser{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&param)
	if err != nil {
		if log != nil {
			log.Error("HandleCreateUser: invalid request body")
		}
		handler.RespondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	hashPass, err := utlis.HashPassword(param.Password)
	if err != nil {
		if log != nil {
			log.Error(fmt.Sprintf("HandleCreateUser: failed to hash password for %s: %v", param.Email, err))
		}
		handler.RespondWithError(w, http.StatusInternalServerError, "Failed to process password")
		return
	}
	var socialLinks json.RawMessage
	if param.Social != nil {
		socialLinks = *param.Social
	}
	userParam := database.CreateUserParams{
		Name:        param.Name,
		Email:       param.Email,
		Password:    hashPass,
		Batch:       param.Batch,
		Description: utlis.ToNilStr(&param.Description),
		Status:      database.AccountStatusPending,
		Role:        database.UserRoleUser,
		Column10:    socialLinks,
	}
	fmt.Println(userParam)
	user, err := h.Repo.CreateUser(r.Context(), userParam)
	if err != nil {
		if log != nil {
			log.Error(fmt.Sprintf("HandleCreateUser: failed to create user %s: %v", param.Email, err))
		}
		handler.RespondWithError(w, http.StatusInternalServerError, "Failed to create user")
		return
	}
	if log != nil {
		log.Info(fmt.Sprintf("HandleCreateUser: user created with email %s", param.Email))
	}
	handler.RespondWithJson(w, http.StatusCreated, user)
}

func (h *Handler) HandleUpdateUserProfile(w http.ResponseWriter, r *http.Request) {
	log, _ := logger.GetLogger()
	user, ok := middleware.GetUser(r.Context())
	if !ok {
		if log != nil {
			log.Error("HandleUpdateUserProfile: unauthenticated request")
		}
		handler.RespondWithError(w, http.StatusUnauthorized, "Authentication required")
		return
	}
	id := chi.URLParam(r, "id")
	userId, err := uuid.Parse(id)
	if err != nil {
		if log != nil {
			log.Error(fmt.Sprintf("HandleUpdateUserProfile: invalid user ID format '%s': %v", id, err))
		}
		handler.RespondWithError(w, http.StatusBadRequest, "Invalid user ID format")
		return
	}
	if user.Role != database.UserRoleAdmin && user.ID != userId {
		if log != nil {
			log.Error(fmt.Sprintf("HandleUpdateUserProfile: user %s unauthorized to update profile %s", user.ID, userId))
		}
		handler.RespondWithError(w, http.StatusForbidden, "You are not authorized to update this profile")
		return
	}
	req := model.PatchUserProfileRequest{}
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&req)
	if err != nil {
		if log != nil {
			log.Error(fmt.Sprintf("HandleUpdateUserProfile: failed to decode request body: %v", err))
		}
		handler.RespondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	params := database.PatchUserProfileParams{
		ID:          userId,
		Name:        utlis.ToNilStr(req.Name),
		Email:       utlis.ToNilStr(req.Email),
		Batch:       utlis.ToNilStr(req.Batch),
		Description: utlis.ToNilStr(req.Desc),
	}
	if req.Social != nil {
		params.SocialLinks = pqtype.NullRawMessage{
			RawMessage: *req.Social,
			Valid:      true,
		}
	}
	updatedUser, err := h.Repo.PatchUserProfile(r.Context(), params)

	if err != nil {
		if utlis.IsNotFound(err) {
			handler.RespondWithError(w, http.StatusNotFound, "User profile not found")
			return
		}
		if log != nil {
			log.Error(fmt.Sprintf("HandleUpdateUserProfile: failed to update profile for user %s: %v", userId, err))
		}
		handler.RespondWithError(w, http.StatusInternalServerError, "Failed to update user profile")
		return
	}
	if log != nil {
		log.Info(fmt.Sprintf("HandleUpdateUserProfile: profile updated for user %s", userId))
	}
	handler.RespondWithJson(w, http.StatusOK, updatedUser)
}

func (h *Handler) HandleImageChange(w http.ResponseWriter, r *http.Request) {
	log, _ := logger.GetLogger()
	user, ok := middleware.GetUser(r.Context())
	if !ok {
		if log != nil {
			log.Error("HandleImageChange: unauthenticated request")
		}
		handler.RespondWithError(w, http.StatusUnauthorized, "Authentication required")
		return
	}
	err := r.ParseMultipartForm(20 << 20)
	if err != nil {
		if log != nil {
			log.Error(fmt.Sprintf("HandleImageChange: failed to parse multipart form for user %s: %v", user.ID, err))
		}
		handler.RespondWithError(w, http.StatusBadRequest, "Failed to parse form data")
		return
	}
	path := fmt.Sprintf("uploads/%s", user.ID.String())
	params := database.PatchUserImagesParams{
		ID: user.ID,
	}
	hasUpdate := false

	userfile, _, err := r.FormFile("user_image")
	if err == nil && userfile != nil {
		defer userfile.Close()
		userImageFilePath, saveErr := utlis.SaveLocal(userfile, "user", path)
		if saveErr != nil {
			if log != nil {
				log.Error(fmt.Sprintf("HandleImageChange: failed to save user image for user %s: %v", user.ID, saveErr))
			}
			handler.RespondWithError(w, http.StatusInternalServerError, "Failed to save profile image")
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
				log.Error(fmt.Sprintf("HandleImageChange: failed to save banner image for user %s: %v", user.ID, saveErr))
			}
			handler.RespondWithError(w, http.StatusInternalServerError, "Failed to save banner image")
			return
		}
		params.BannerImage = utlis.ToNilStr(&bannerImageFilePath)
		hasUpdate = true
	}

	if !hasUpdate {
		if log != nil {
			log.Error(fmt.Sprintf("HandleImageChange: no image provided in request for user %s", user.ID))
		}
		handler.RespondWithError(w, http.StatusBadRequest, "At least one image (user_image or banner_image) is required")
		return
	}
	res, err := h.Repo.PatchUserImages(r.Context(), params)
	if err != nil {
		if utlis.IsNotFound(err) {
			handler.RespondWithError(w, http.StatusNotFound, "User profile not found")
			return
		}
		if log != nil {
			log.Error(fmt.Sprintf("HandleImageChange: failed to update images for user %s: %v", user.ID, err))
		}
		handler.RespondWithError(w, http.StatusInternalServerError, "Failed to update images")
		return
	}
	if log != nil {
		log.Info(fmt.Sprintf("HandleImageChange: images updated for user %s", user.ID))
	}
	handler.RespondWithJson(w, http.StatusOK, res)
}

func (h *Handler) HandlePasswordChange(w http.ResponseWriter, r *http.Request) {
	log, _ := logger.GetLogger()
	user, ok := middleware.GetUser(r.Context())
	if !ok {
		if log != nil {
			log.Error("HandlePasswordChange: unauthenticated request")
		}
		handler.RespondWithError(w, http.StatusUnauthorized, "Authentication required")
		return
	}
	req := model.PatchUserPassword{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&req)
	if err != nil {
		if log != nil {
			log.Error(fmt.Sprintf("HandlePasswordChange: failed to decode request body for user %s: %v", user.ID, err))
		}
		handler.RespondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	dbUser, err := h.Repo.GetUserByEmail(r.Context(), user.Email)
	if err != nil {
		if utlis.IsNotFound(err) {
			if log != nil {
				log.Info(fmt.Sprintf("HandlePasswordChange: user email %s not found in DB: %v", user.Email, err))
			}
			handler.RespondWithError(w, http.StatusNotFound, "User not found")
			return
		}
		if log != nil {
			log.Error(fmt.Sprintf("HandlePasswordChange: DB error looking up email %s: %v", user.Email, err))
		}
		handler.RespondWithError(w, http.StatusInternalServerError, "Failed to verify user")
		return
	}
	if !utlis.CheckPassword(dbUser.Password, req.OldPassword) {
		if log != nil {
			log.Info(fmt.Sprintf("HandlePasswordChange: incorrect old password for user %s", user.ID))
		}
		handler.RespondWithError(w, http.StatusUnauthorized, "Current password is incorrect")
		return
	}

	password, err := utlis.HashPassword(req.Password)
	if err != nil {
		if log != nil {
			log.Error(fmt.Sprintf("HandlePasswordChange: failed to hash new password for user %s: %v", user.ID, err))
		}
		handler.RespondWithError(w, http.StatusInternalServerError, "Failed to process new password")
		return
	}
	params := database.PatchUserPasswordParams{
		ID:       user.ID,
		Password: password,
	}
	res, err := h.Repo.PatchUserPassword(r.Context(), params)
	if err != nil {
		if utlis.IsNotFound(err) {
			handler.RespondWithError(w, http.StatusNotFound, "User not found")
			return
		}
		if log != nil {
			log.Error(fmt.Sprintf("HandlePasswordChange: failed to update password for user %s: %v", user.ID, err))
		}
		handler.RespondWithError(w, http.StatusInternalServerError, "Failed to update password")
		return
	}
	if log != nil {
		log.Info(fmt.Sprintf("HandlePasswordChange: password updated for user %s", user.ID))
	}
	handler.RespondWithJson(w, http.StatusOK, res)
}
