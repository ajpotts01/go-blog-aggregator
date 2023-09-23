package main

// Import Postgres driver w/ side effects
import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/ajpotts01/go-blog-aggregator/api"
	"github.com/ajpotts01/go-blog-aggregator/internal/database"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func getApiConfig(dbConnStr string) (*api.ApiConfig, error) {
	db, err := sql.Open("postgres", dbConnStr)

	if err != nil {
		return &api.ApiConfig{}, err
	}

	dbq := database.New(db)

	return &api.ApiConfig{
		DbConn: dbq,
	}, nil
}

func getApiRouterV1(config *api.ApiConfig) *chi.Mux {
	const errEndpoint = "/err"
	const readyEndpoint = "/readiness"
	const usersEndpoint = "/users"
	const feedsEndpoint = "/feeds"
	const followsEndpoint = "/follows"
	const singleFollowEndpoint = "/follows/{id}"

	apiRouter := chi.NewRouter()
	apiRouter.Get(readyEndpoint, api.Ready)
	apiRouter.Get(errEndpoint, api.Err)
	apiRouter.Post(usersEndpoint, config.CreateUser)
	apiRouter.Get(usersEndpoint, config.AuthMiddleware(config.GetUser))
	apiRouter.Post(feedsEndpoint, config.AuthMiddleware(config.CreateFeed))
	apiRouter.Get(feedsEndpoint, config.GetFeeds)
	apiRouter.Post(followsEndpoint, config.AuthMiddleware(config.FollowFeed))
	apiRouter.Get(followsEndpoint, config.AuthMiddleware(config.GetFollows))
	apiRouter.Delete(singleFollowEndpoint, config.AuthMiddleware(config.UnfollowFeed))

	return apiRouter
}

func main() {
	godotenv.Load()
	port := os.Getenv("PORT")
	dbConnStr := os.Getenv("PG_CONN")

	// Database
	apiConfig, err := getApiConfig(dbConnStr)

	fmt.Printf("API config: %v", apiConfig)

	if err != nil {
		log.Fatalf("Error setting up database: %v", err)
	}

	// App router
	appRouter := chi.NewRouter()

	// CORS
	corsOptions := cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET, POST, OPTIONS, PUT, DELETE"},
		AllowedHeaders: []string{"*"},
	}
	appRouter.Use(cors.Handler(corsOptions))

	appRouter.Mount("/v1", getApiRouterV1(apiConfig))

	server := &http.Server{
		Addr:    ":" + port,
		Handler: appRouter,
	}

	// url := "https://blog.boot.dev/index.xml"
	// altUrl := "https://wagslane.dev/index.xml"
	// log.Printf("URL: %v", url)
	// log.Printf("Alt URL: %v", url)
	// log.Printf("Using alt URL")
	//feed, err := apiConfig.FetchFeed(altUrl)

	// if err != nil {
	// 	log.Printf("Error getting feed: %v", err)
	// }

	// for _, c := range feed.Channels {
	// 	log.Printf("Channel Link: %v", c.Link)
	// 	//log.Printf("Atom: %v", c.AtomLink)
	// 	log.Printf("Items:")
	// 	for _, i := range c.Items {
	// 		log.Printf("Title: %v", i.Title)
	// 		log.Printf("Desc: %v", i.Description)
	// 		log.Printf("Publish Date: %v", i.PubDate)
	// 		log.Printf("Link: %v", i.Link)
	// 	}
	// }

	log.Printf("Now serving on port: %v", port)
	log.Fatal(server.ListenAndServe())
}
