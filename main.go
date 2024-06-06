package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
)

func main() {
	fmt.Print()

	godotenv.Load()

	portStr := os.Getenv("PORT")
	if portStr == "" {
		log.Fatal("PORT not found")
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

	router.Mount("/v1", v1Router)

	server := &http.Server {
		Handler: router,
		Addr: ":" + portStr,
	}

	log.Printf("Server starting on %v port", portStr)
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}