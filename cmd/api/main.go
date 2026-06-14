package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/sakshipatel29/launchguard/internal/cache"
	"github.com/sakshipatel29/launchguard/internal/db"
	"github.com/sakshipatel29/launchguard/internal/events"
	"github.com/sakshipatel29/launchguard/internal/handlers"
	"github.com/sakshipatel29/launchguard/internal/store"
)

func main() {
	ctx := context.Background()

	database, err := db.ConnectPostgres()
	if err != nil {
		log.Fatal("failed to connect to PostgreSQL:", err)
	}
	defer database.Close()

	if err := db.RunMigrations(database); err != nil {
		log.Fatal("failed to run database migrations:", err)
	}

	redisClient, err := cache.ConnectRedis(ctx)
	if err != nil {
		log.Fatal("failed to connect to Redis:", err)
	}
	defer redisClient.Close()

	kafkaBrokers := os.Getenv("KAFKA_BROKERS")
	if kafkaBrokers == "" {
		kafkaBrokers = "localhost:9092"
	}

	eventPublisher := events.NewKafkaPublisher(
		strings.Split(kafkaBrokers, ","),
		"feature_flag_evaluations",
	)
	defer eventPublisher.Close()

	r := chi.NewRouter()

	postgresStore := store.NewPostgresFeatureFlagStore(database)
	cachedStore := store.NewCachedFeatureFlagStore(postgresStore, redisClient, 5*time.Minute)
	flagHandler := handlers.NewFeatureFlagHandler(cachedStore, eventPublisher)

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

	err = http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatal("server failed to start:", err)
	}
}
