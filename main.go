package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/ajpotts01/go-blog-aggregator/api"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
)

func getApiRouterV1() *chi.Mux {
	const errEndpoint = "/err"
	const readyEndpoint = "/readiness"

	apiRouter := chi.NewRouter()
	apiRouter.Get(readyEndpoint, api.Ready)
	apiRouter.Get(errEndpoint, api.Err)

	return apiRouter
}

func main() {
	godotenv.Load()
	port := os.Getenv("PORT")
	fmt.Println(port)

	// App router
	appRouter := chi.NewRouter()

	// CORS
	corsOptions := cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET, POST, OPTIONS, PUT, DELETE"},
		AllowedHeaders: []string{"*"},
	}
	appRouter.Use(cors.Handler(corsOptions))

	appRouter.Mount("/v1", getApiRouterV1())

	server := &http.Server{
		Addr:    ":" + port,
		Handler: appRouter,
	}

	log.Printf("Now serving on port: %v", port)
	log.Fatal(server.ListenAndServe())
}
