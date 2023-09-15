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

	apiRouter := chi.NewRouter()
	apiRouter.Get(readyEndpoint, api.Ready)
	apiRouter.Get(errEndpoint, api.Err)
	apiRouter.Post(usersEndpoint, config.CreateUser)
	apiRouter.Get(usersEndpoint, config.AuthMiddleware(config.GetUser))
	apiRouter.Post(feedsEndpoint, config.AuthMiddleware(config.CreateFeed))

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

	log.Printf("Now serving on port: %v", port)
	log.Fatal(server.ListenAndServe())
}
