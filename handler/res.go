package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
)

func RespondWithJson(w http.ResponseWriter, code int, payload interface{}) {
	data, err := json.Marshal(payload)
	if err != nil {
		log.Println("Failed to parse Json")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-type", "Application/Json")
	w.WriteHeader(code)
	w.Write(data)

}
func RespondWithError(w http.ResponseWriter, code int, msg string) {
	if code > 499 {
		log.Println("Responding with 5xx error", msg)
	}
	type response struct {
		Msg string `json:"Error"`
	}
	RespondWithJson(w, code, response{Msg: msg})
}
func Health(w http.ResponseWriter, r *http.Request) {
	type res struct {
		Health string
	}
	response := res{
		Health: "ok",
	}
	RespondWithJson(w, http.StatusOK, response)
}
func RespondWithHTML(w http.ResponseWriter, code int, payload string) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(code)
	w.Write([]byte(payload))
}
func MainPage(w http.ResponseWriter, r *http.Request) {
	content, err := os.ReadFile("template/index.html")
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Could not read template")
		return
	}
	RespondWithHTML(w, http.StatusOK, string(content))
}
