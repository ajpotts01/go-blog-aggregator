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

type createFeedRequest struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type feedResponse struct {
	Id        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`
	Url       string    `json:"url"`
	UserId    uuid.UUID `json:"user_id"`
}

type feedList struct {
	Feeds []feedResponse `json:"feeds"`
}

func createFeedParams(name string, url string, userId uuid.UUID) (database.CreateFeedParams, error) {
	newId, err := uuid.NewUUID()

	if err != nil {
		return database.CreateFeedParams{}, err
	}

	createdAt := time.Now()

	params := database.CreateFeedParams{
		ID:        newId,
		Name:      name,
		CreatedAt: createdAt,
		UpdatedAt: createdAt,
		Url:       url,
		UserID:    userId,
	}

	return params, nil
}

// POST /api/feeds
func (config *ApiConfig) CreateFeed(w http.ResponseWriter, r *http.Request, user database.User) {
	decoder := json.NewDecoder(r.Body)
	requestParams := createFeedRequest{}
	err := decoder.Decode(&requestParams)

	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	dbFeedParams, err := createFeedParams(requestParams.Name, requestParams.Url, user.ID)

	if err != nil {
		log.Printf("Error creating new feed params: %v", err)
		errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	newFeed, err := config.DbConn.CreateFeed(context.TODO(), dbFeedParams)

	if err != nil {
		log.Printf("Error creating new feed: %v", err)
		errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	validResponse(w, http.StatusCreated, feedResponse{
		Id:        newFeed.ID,
		CreatedAt: newFeed.CreatedAt,
		UpdatedAt: newFeed.UpdatedAt,
		Name:      newFeed.Name,
		Url:       newFeed.Url,
		UserId:    newFeed.UserID,
	})
	return
}

// GET /api/feeds
func (config *ApiConfig) GetFeeds(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	feeds, err := config.DbConn.GetFeeds(context.TODO())
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, "Error retrieving feeds")
		return
	}

	var returnedFeeds []feedResponse

	for _, feed := range feeds {
		returnedFeeds = append(returnedFeeds, feedResponse{
			Id:        feed.ID,
			CreatedAt: feed.CreatedAt,
			UpdatedAt: feed.UpdatedAt,
			Name:      feed.Name,
			Url:       feed.Url,
			UserId:    feed.UserID,
		})
	}

	validResponse(w, http.StatusOK, returnedFeeds)
	return
}
