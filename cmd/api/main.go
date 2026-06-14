package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/sakshipatel29/launchguard/internal/handlers"
)

func main() {
	r := chi.NewRouter()

	r.Get("/health", handlers.HealthCheck)

	log.Println("LaunchGuard API running on port 8080")
	err := http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatal("server failed to start:", err)
	}
}