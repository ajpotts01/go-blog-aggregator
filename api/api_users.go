package api

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/ajpotts01/go-blog-aggregator/internal/database"
	"github.com/google/uuid"
)

type createUserRequest struct {
	Name string `json:"name"`
}

type userResponse struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`
	ApiKey    string    `json:"api_key"`
}

func createUserParams(name string) (database.CreateUserParams, error) {
	newId, err := uuid.NewUUID()

	if err != nil {
		return database.CreateUserParams{}, err
	}

	createdAt := time.Now()

	// API key will be generated by the database
	params := database.CreateUserParams{
		ID:        newId,
		Name:      name,
		CreatedAt: createdAt,
		UpdatedAt: createdAt,
	}

	return params, nil
}

// POST /api/users
func (config *ApiConfig) CreateUser(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	requestParams := createUserRequest{}
	err := decoder.Decode(&requestParams)

	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	dbUsrParams, err := createUserParams(requestParams.Name)

	if err != nil {
		log.Printf("Error creating new user params: %v", err)
		errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	// what is ctx?
	// context.TODO() seems to be the convention?
	newUser, err := config.DbConn.CreateUser(context.TODO(), dbUsrParams)

	if err != nil {
		log.Printf("Error creating new user: %v", err)
		errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	validResponse(w, http.StatusCreated, userResponse{
		ID:        newUser.ID,
		CreatedAt: newUser.CreatedAt,
		UpdatedAt: newUser.UpdatedAt,
		Name:      newUser.Name,
		ApiKey:    newUser.ApiKey,
	})
	return
}

// GET /api/users
func (config *ApiConfig) GetUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// No body expected: just API key header
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

	validResponse(w, http.StatusOK, userResponse{
		ID:        usr.ID,
		CreatedAt: usr.CreatedAt,
		UpdatedAt: usr.UpdatedAt,
		Name:      usr.Name,
		ApiKey:    usr.ApiKey,
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
