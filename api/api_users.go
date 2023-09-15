package api

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/ajpotts01/go-blog-aggregator/internal/database"
	"github.com/google/uuid"
)

type createUserRequest struct {
	Name string `json:"name"`
}

type createUserResponse struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`
}

func createUserParams(name string) (database.CreateUserParams, error) {
	newId, err := uuid.NewUUID()

	if err != nil {
		return database.CreateUserParams{}, err
	}

	createdAt := time.Now()

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

	validResponse(w, http.StatusCreated, createUserResponse{
		ID:        newUser.ID,
		CreatedAt: newUser.CreatedAt,
		UpdatedAt: newUser.UpdatedAt,
		Name:      newUser.Name,
	})
	return
}
