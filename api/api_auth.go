package api

import (
	"context"
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/ajpotts01/go-blog-aggregator/internal/database"
)

type authorisedMethod func(http.ResponseWriter, *http.Request, database.User)

func (config *ApiConfig) AuthMiddleware(method authorisedMethod) http.HandlerFunc {
	// No body expected: just API key header
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		log.Printf("Now in auth middleware")

		key, err := config.getAuthFromHeader(r, "ApiKey")
		if err != nil {
			errorResponse(w, http.StatusUnauthorized, "Bad authorization header")
			return
		}

		log.Printf("Received API key: %v", key)
		usr, err := config.DbConn.GetUserByApiKey(context.TODO(), key)

		if err != nil {
			errorResponse(w, http.StatusUnauthorized, "Bad API key")
			return
		}

		method(w, r, usr)
	})
}

func (config *ApiConfig) getAuthFromHeader(r *http.Request, tokenType string) (string, error) {
	authHeader := r.Header.Get("Authorization")

	if authHeader == "" {
		return "", errors.New("must supply authorization header")
	}

	suppliedItem := strings.Replace(authHeader, tokenType, "", 1)
	suppliedItem = strings.Trim(suppliedItem, " ")
	return suppliedItem, nil
}
