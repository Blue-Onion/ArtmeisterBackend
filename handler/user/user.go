package user

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/Blue-Onion/RestApi-Go/handler"
	"github.com/Blue-Onion/RestApi-Go/internal/database"
	"github.com/Blue-Onion/RestApi-Go/model"
	"github.com/Blue-Onion/RestApi-Go/utlis"
	"github.com/google/uuid"
	"net/http"
	"time"
)

type Handler struct {
	Repo database.UserRepository
}

func (h *Handler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	params := model.AutheticateUser{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {
		handler.RespondWithError(w, 400, "Error in Parsing Json")
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
		Name:     "auth_token",
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
		ID:        uuid.New(),
		Name:      param.Name,
		Email:     param.Email,
		Password:  hashPass,
		Createdat: time.Now(),
		Updatedat: time.Now(),
	})
	if err != nil {
		handler.RespondWithError(w, 500, "Couldn't create user")
		return
	}

	handler.RespondWithJson(w, 201, user)
}
