package api

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/ajpotts01/go-blog-aggregator/internal/database"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type createFeedRequest struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type followFeedRequest struct {
	FeedId uuid.UUID `json:"feed_id"`
}

type feedResponse struct {
	Id            uuid.UUID  `json:"id"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	LastFetchedAt *time.Time `json:"last_fetched_at"`
	Name          string     `json:"name"`
	Url           string     `json:"url"`
	UserId        uuid.UUID  `json:"user_id"`
}

type followResponse struct {
	Id        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	UserId    uuid.UUID `json:"user_id"`
	FeedId    uuid.UUID `json:"feed_id"`
}

type newFeedResponse struct {
	Feed   feedResponse   `json:"feed"`
	Follow followResponse `json:"follow"`
}

type feedList struct {
	Feeds []feedResponse `json:"feeds"`
}

type postResponse struct {
	Id          uuid.UUID `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Title       string    `json:"title"`
	Url         string    `json:"url"`
	Description string    `json:"description"`
	PublishedAt time.Time `json:"published_at"`
	FeedID      uuid.UUID `json:"feed_id"`
	FeedName    string    `json:"feed_name"`
	FeedUrl     string    `json:"feed_url"`
}

func mapFeedResponse(feed database.Feed) feedResponse {
	return feedResponse{
		Id:            feed.ID,
		UserId:        feed.UserID,
		CreatedAt:     feed.CreatedAt,
		UpdatedAt:     feed.UpdatedAt,
		LastFetchedAt: &feed.LastFetchedAt.Time,
		Name:          feed.Name,
		Url:           feed.Url,
	}
}

func mapFollowResponse(follow database.Follow) followResponse {
	return followResponse{
		Id:        follow.ID,
		UserId:    follow.UserID,
		FeedId:    follow.FeedID,
		CreatedAt: follow.CreatedAt,
		UpdatedAt: follow.UpdatedAt,
	}
}

func createNewFeedResponse(feed database.Feed, follow database.Follow) newFeedResponse {
	return newFeedResponse{
		Feed:   mapFeedResponse(feed),
		Follow: mapFollowResponse(follow),
	}
}

func getPostByUserParams(userId uuid.UUID) database.GetPostsByUserParams {
	params := database.GetPostsByUserParams{
		UserID: userId,
		Limit:  5,
	}

	return params
}

func createFollowParams(userId uuid.UUID, feedId uuid.UUID) (database.CreateFollowParams, error) {
	newId, err := uuid.NewUUID()

	if err != nil {
		return database.CreateFollowParams{}, err
	}

	createdAt := time.Now()

	params := database.CreateFollowParams{
		ID:        newId,
		CreatedAt: createdAt,
		UpdatedAt: createdAt,
		UserID:    userId,
		FeedID:    feedId,
	}

	return params, nil
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

func unfollowParams(followId uuid.UUID, userId uuid.UUID) database.DeleteFollowParams {
	params := database.DeleteFollowParams{
		ID:     followId,
		UserID: userId,
	}

	return params
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

	dbFollowParams, err := createFollowParams(user.ID, newFeed.ID)

	if err != nil {
		log.Printf("Error creating new follow params: %v", err)
		errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	newFollow, err := config.DbConn.CreateFollow(context.TODO(), dbFollowParams)

	if err != nil {
		log.Printf("Error creating new follow: %v", err)
		errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	responses := createNewFeedResponse(newFeed, newFollow)

	validResponse(w, http.StatusCreated, responses)
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
		returnedFeeds = append(returnedFeeds, mapFeedResponse(feed))
	}

	validResponse(w, http.StatusOK, returnedFeeds)
	return
}

// POST /api/follows
func (config *ApiConfig) FollowFeed(w http.ResponseWriter, r *http.Request, user database.User) {
	decoder := json.NewDecoder(r.Body)
	requestParams := followFeedRequest{}
	err := decoder.Decode(&requestParams)

	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	dbFollowParams, err := createFollowParams(user.ID, requestParams.FeedId)

	if err != nil {
		log.Printf("Error creating new feed params: %v", err)
		errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	newFollow, err := config.DbConn.CreateFollow(context.TODO(), dbFollowParams)

	if err != nil {
		log.Printf("Error creating new follow: %v", err)
		errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	validResponse(w, http.StatusCreated, mapFollowResponse(newFollow))
	return
}

// DELETE /api/follows/{id}
func (config *ApiConfig) UnfollowFeed(w http.ResponseWriter, r *http.Request, user database.User) {
	providedId := chi.URLParam(r, "id")
	followId, err := uuid.Parse(providedId)

	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	dbUnfollowParams := unfollowParams(followId, user.ID)

	log.Printf("Deleting: Follow %v owned by %v", followId, user.ID)
	err = config.DbConn.DeleteFollow(context.TODO(), dbUnfollowParams)

	if err != nil {
		log.Printf("Error deleting existing follow: %v", err)
		errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	return
}

// GET /api/follows
func (config *ApiConfig) GetFollows(w http.ResponseWriter, r *http.Request, user database.User) {
	w.Header().Set("Content-Type", "application/json")

	follows, err := config.DbConn.GetFollows(context.TODO(), user.ID)
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, "Error retrieving follows")
		return
	}

	var returnedFollows []followResponse

	for _, follow := range follows {
		returnedFollows = append(returnedFollows, mapFollowResponse(follow))
	}

	validResponse(w, http.StatusOK, returnedFollows)
	return
}

// GET /api/posts
func (config *ApiConfig) GetPostsForUser(w http.ResponseWriter, r *http.Request, user database.User) {
	var returnedPosts []postResponse

	w.Header().Set("Content-Type", "application/json")

	params := getPostByUserParams(user.ID)
	posts, err := config.DbConn.GetPostsByUser(context.TODO(), params)

	if err != nil {
		errorResponse(w, http.StatusInternalServerError, "Error retrieving posts")
		return
	}

	for _, post := range posts {
		returnedPosts = append(returnedPosts, postResponse{
			Id:          post.ID,
			CreatedAt:   post.CreatedAt,
			UpdatedAt:   post.UpdatedAt,
			Title:       post.Title,
			Url:         post.Url,
			Description: post.Description.String,
			PublishedAt: post.PublishedAt.Time,
			FeedID:      post.FeedID,
			FeedName:    post.FeedName,
			FeedUrl:     post.FeedUrl,
		})
	}

	validResponse(w, http.StatusOK, returnedPosts)
	return
}
