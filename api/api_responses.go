package api

import (
	"encoding/json"
	"log"
	"net/http"
)

func validResponse(w http.ResponseWriter, code int, obj interface{}) {
	w.Header().Set("Content-Type", "application/json")

	resp, err := json.Marshal(obj)

	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}

	w.WriteHeader(code)
	w.Write(resp)
}

func errorResponse(w http.ResponseWriter, code int, errorMsg string) {
	type errorResponse struct {
		Err string `json:"error"`
	}

	errorObj := errorResponse{
		Err: errorMsg,
	}

	validResponse(w, code, errorObj)
}
