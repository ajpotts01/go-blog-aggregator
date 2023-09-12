package api

import "net/http"

type readyResponse struct {
	Status string `json:"status"`
}

func Ready(w http.ResponseWriter, r *http.Request) {
	validResponse(w, http.StatusOK, readyResponse{
		Status: "ok",
	})
}

func Err(w http.ResponseWriter, r *http.Request) {
	errorResponse(w, http.StatusInternalServerError, "Internal Server Error")
}
