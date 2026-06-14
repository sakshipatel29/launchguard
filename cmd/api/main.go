package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/sakshipatel29/launchguard/internal/handlers"
	"github.com/sakshipatel29/launchguard/internal/store"
)

func main() {
	r := chi.NewRouter()

	flagStore := store.NewFeatureFlagStore()
	flagHandler := handlers.NewFeatureFlagHandler(flagStore)

	r.Get("/health", handlers.HealthCheck)
	r.Post("/evaluate", flagHandler.EvaluateFlag)

	r.Route("/flags", func(r chi.Router) {
		r.Post("/", flagHandler.CreateFlag)
		r.Get("/", flagHandler.ListFlags)
		r.Get("/{id}", flagHandler.GetFlag)
		r.Put("/{id}", flagHandler.UpdateFlag)
		r.Delete("/{id}", flagHandler.DeleteFlag)
	})

	log.Println("LaunchGuard API running on port 8080")

	err := http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatal("server failed to start:", err)
	}
}
