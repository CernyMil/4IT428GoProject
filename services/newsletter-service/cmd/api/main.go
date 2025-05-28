package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"newsletter-service/repository"
	v1 "newsletter-service/transport/api/v1"
	"newsletter-service/transport/middleware"

	firebase "firebase.google.com/go"
	"google.golang.org/api/option"

	_ "github.com/lib/pq" // PostgreSQL driver
)

func initializeFirebase(credPath string) (*firebase.App, error) {
	if credPath == "" {
		log.Fatal("FIREBASE_CRED is not set in the configuration")
	}

	opt := option.WithCredentialsFile(credPath)
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		return nil, err
	}
	return app, nil
}

func connectDatabase(databaseURL string) (*sql.DB, error) {
	if databaseURL == "" {
		log.Fatal("DATABASE_URL is empty in the configuration")
	}

	var db *sql.DB
	var err error

	for i := 0; i < 10; i++ {
		db, err = sql.Open("postgres", databaseURL)
		if err == nil {
			err = db.Ping()
			if err == nil {
				break
			}
		}
		log.Println("Waiting for database to be ready...")
		time.Sleep(3 * time.Second)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to connect to the database: %v", err)
	}

	log.Println("Connected to the database successfully!")
	return db, nil
}

func main() {

	// Load configuration
	cfg := MustLoadConfig()

	// Connect to the database
	db, err := connectDatabase(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Database connection error: %v", err)
	}
	defer db.Close()

	// Initialize Firebase
	firebaseApp, err := initializeFirebase(cfg.FirebaseCred)
	if err != nil {
		log.Fatalf("Failed to initialize Firebase: %v", err)
	}

	// Initialize repository
	repo := repository.NewPostgresRepository(db)

	// Initialize handler
	handler := v1.NewNewsletterHandler(repo)

	// Set up routes
	r := handler.Routes()
	r.Use(middleware.FirebaseAuthMiddleware(firebaseApp))

	// Start the server
	addr := fmt.Sprintf(":%d", cfg.Port)
	log.Printf("Starting server on %s...", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
