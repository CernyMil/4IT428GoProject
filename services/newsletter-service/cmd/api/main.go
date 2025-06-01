package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"

	"newsletter-service/repository"
	"newsletter-service/service"
	v1 "newsletter-service/transport/api/v1"
	"newsletter-service/transport/middleware"
	"newsletter-service/transport/util"

	firebase "firebase.google.com/go"
	"google.golang.org/api/option"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq" // PostgreSQL driver
)

func initializeFirebase(credPath string) (*firebase.App, error) {
	opt := option.WithCredentialsFile(credPath)
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		return nil, err
	}
	return app, nil
}

func main() {
	// Load configuration
	ctx := context.Background()
	cfg := MustLoadConfig()
	util.SetServerLogLevel(slog.LevelInfo)

	// Connect to the database
	database, err := setupDatabase(ctx, cfg)
	if err != nil {
		slog.Error("initializing database", slog.Any("error", err))
	}

	// Initialize Firebase
	firebaseApp, err := initializeFirebase(cfg.FirebaseCred)
	if err != nil {
		log.Fatalf("Failed to initialize Firebase: %v", err)
	}

	// Initialize repository
	repo, err := repository.NewPostgresRepository(database)
	if err != nil {
		slog.Error("initializing repository", slog.Any("error", err))
	}

	// Initialize authenticator
	authenticator := middleware.NewFirebaseAuthenticator(firebaseApp)

	// Initialize service
	newsletterService, err := service.NewService(repo)
	if err != nil {
		slog.Error("initializing service", slog.Any("error", err))
		return // Exit the function or handle the error appropriately
	}

	// Initialize handler
	handler := v1.NewHandler(*authenticator, newsletterService)

	// Set up routes
	r := chi.NewRouter()
	r.Mount("/", handler.Routes())

	// Start the server
	addr := fmt.Sprintf(":%d", cfg.Port)
	log.Printf("Starting server on %s...", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func setupDatabase(ctx context.Context, cfg Config) (*pgxpool.Pool, error) {
	// Initialize the database connection pool.
	pool, err := pgxpool.New(
		ctx,
		cfg.DatabaseURL,
	)
	if err != nil {
		return nil, err
	}
	return pool, nil
}
