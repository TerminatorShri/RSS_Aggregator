package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/TerminatorShri/RSS_Aggregator/internal/database"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

type apiConfig struct {
	DB *database.Queries
}

func main() {
	fmt.Print()

	godotenv.Load()

	portStr := os.Getenv("PORT")
	if portStr == "" {
		log.Fatal("PORT not found")
	}

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("Database URL not found")
	}

	conn, err:=	sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Cannot connect to Database")
	}

	apiCfg := apiConfig {
		DB: database.New(conn),
	}
	

	router := chi.NewRouter()
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, 
	  }))

	v1Router := chi.NewRouter()
	v1Router.Get("/healthz", handelerReadiness)
	v1Router.Get("/err", handlerError)
	v1Router.Get("/users", apiCfg.middlewareAuth(apiCfg.handlerGetUser))
	v1Router.Get("/feeds", apiCfg.handlerGetFeeds)
	v1Router.Get("/feed-follows", apiCfg.middlewareAuth(apiCfg.handlerGetFeedFollow))
	v1Router.Post("/users", apiCfg.handlerCreateUser)
	v1Router.Post("/feeds", apiCfg.middlewareAuth(apiCfg.handlerCreateFeed))
	v1Router.Post("/feed-follows", apiCfg.middlewareAuth(apiCfg.handlerCreateFeedFollow))
	v1Router.Delete("/feed-follows/{feedFollowID}", apiCfg.middlewareAuth(apiCfg.handlerDeleteFeedFollow))

	router.Mount("/v1", v1Router)

	server := &http.Server {
		Handler: router,
		Addr: ":" + portStr,
	}

	log.Printf("Server starting on %v port", portStr)
	err = server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}