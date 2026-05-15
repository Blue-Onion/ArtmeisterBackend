package user

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/Blue-Onion/ArtmeisterBackend/handler"
	"github.com/Blue-Onion/ArtmeisterBackend/internal/database"
	"github.com/Blue-Onion/ArtmeisterBackend/middleware"
	"github.com/Blue-Onion/ArtmeisterBackend/model"
	"github.com/Blue-Onion/ArtmeisterBackend/utlis"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
)

type Handler struct {
	Repo database.UserRepository
}

func (h *Handler) HandleUpdateImg(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value("user").(middleware.User)
	if !ok {
		handler.RespondWithError(w, 400, "Not AutheticateUser")
		return
	}
	err := r.ParseMultipartForm(20 << 20)
	if err != nil {
		handler.RespondWithError(w, 400, err.Error())
		return
	}
	file, fileHeader, err := r.FormFile("image")

	if err != nil {
		handler.RespondWithError(w, 400, err.Error())
		return
	}
	defer file.Close()
	path := fmt.Sprintf("uploads/%s", user.ID.String())
	filepath, err := utlis.SaveLocal(file, fileHeader, path)
	if err != nil {
		handler.RespondWithError(w, 400, err.Error())
		return
	}
	handler.RespondWithJson(w, 200, filepath)
}
func (h *Handler) HandleUpdateUserProfile(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(middleware.User)
	id := chi.URLParam(r, "id")
	userId, err := uuid.Parse(id)
	if err != nil {
		handler.RespondWithError(w, 400, err.Error())
		return
	}
	if user.Role != database.UserRoleAdmin || user.ID != userId {
		handler.RespondWithError(w, 403, "Not Authorized to Change User Data")
		return
	}
	req := model.PatchUserProfileRequest{}
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&req)
	if err != nil {
		handler.RespondWithError(w, 400, err.Error())
		return
	}
	params := database.PatchUserProfileParams{
		ID:    userId,
		Name:  toNilStr(req.Name),
		Email: toNilStr(req.Email),
		Batch: toNilStr(req.Batch),
	}
	updatedUser, err := h.Repo.PatchUserProfile(r.Context(), params)
	if err != nil {
		handler.RespondWithError(w, 400, err.Error())
		return
	}
	handler.RespondWithJson(w, 200, updatedUser)
}
func toNilStr(str *string) sql.NullString {
	if str == nil {
		return sql.NullString{
			Valid: false,
		}
	}
	return sql.NullString{
		String: *str,
		Valid:  true,
	}
}
func (h *Handler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	params := model.AutheticateUser{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {
		handler.RespondWithError(w, 400, err.Error())
		return
	}
	user, err := h.Repo.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {

			handler.RespondWithError(w, 404, "User not found")
			return
		}
		handler.RespondWithError(w, 400, err.Error())
		return

	}
	isValid := utlis.CheckPassword(user.Password, params.Password)
	if !isValid {
		handler.RespondWithError(w, 400, "Incorrect Password")
		return

	}
	token, err := utlis.GenerateJwt(user.ID)
	if err != nil {
		handler.RespondWithError(w, 400, err.Error())
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
	handler.RespondWithJson(w, 200, map[string]string{
		"Message": "Login Successfull",
	})
}
func (h *Handler) HandleLogOut(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "authToken",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	})
	handler.RespondWithJson(w, 200, map[string]string{
		"Message": "LogOut Successfully",
	})
}

func (h *Handler) HandleCreateUser(w http.ResponseWriter, r *http.Request) {
	param := model.CreateUser{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&param)
	if err != nil {
		handler.RespondWithError(w, 400, "Error in Parsing Json")
		return
	}
	hashPass, err := utlis.HashPassword(param.Password)
	if err != nil {
		handler.RespondWithError(w, 400, "Error in Hashing Password")
		return
	}
	user, err := h.Repo.CreateUser(r.Context(), database.CreateUserParams{
		Name:     param.Name,
		Email:    param.Email,
		Password: hashPass,
		Batch:    param.Batch,
		Status:   database.AccountStatusPending,
		Role:     database.UserRoleUser,
	})
	if err != nil {
		handler.RespondWithError(w, 500, err.Error())
		return
	}

	handler.RespondWithJson(w, 201, user)
}
