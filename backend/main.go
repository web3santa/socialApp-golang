package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"social/database"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	DB *database.Queries
}

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dbURL := os.Getenv("DB_URL")
	conn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	apiCfg := apiConfig{
		DB: database.New(conn),
	}

	r := chi.NewRouter()
	r.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	V1Router := chi.NewRouter()
	V1Router.Get("/healthz", handleReadiness)
	V1Router.Get("/err", handleError)

	V1Router.Post("/users", apiCfg.handlerCreateUser)
	V1Router.Get("/users", apiCfg.handleGetUsers)
	V1Router.Get("/user/{id}", apiCfg.handleGetUser)
	V1Router.Patch("/user/{id}", apiCfg.handleUpdateUser)
	V1Router.Delete("/user/{id}", apiCfg.handleDeleteUser)

	r.Mount("/v1", V1Router)

	http.ListenAndServe(":3000", r)

}
