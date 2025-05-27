package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"newsletter-management-api/repository"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq" // PostgreSQL driver
)

func main() {
	// Load environment variables from the .env file.
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Get the port from the environment variables.
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT is not set in the environment variables")
	}

	// Get the database URL from the environment variables.
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL is not set in the environment variables")
	}

	// Connect to the database.
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	defer db.Close()

	// Test the database connection.
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping the database: %v", err)
	}
	log.Println("Connected to the database successfully!")

	// Initialize the repository.
	repo := repository.NewPostgresRepository(db)

	// Example usage of the repository.
	ctx := context.Background()
	newsletter := &repository.Newsletter{
		//ID:        "12",
		Subject:   "Vinko je fajne",
		Body:      "Daj si vínko, budeš ho potrebovať!",
		CreatedAt: time.Now(),
	}
	if err := repo.Save(ctx, newsletter); err != nil {
		log.Fatalf("Failed to save newsletter: %v", err)
	}

	log.Println("Newsletter saved successfully!")

	// Start the server.
	addr := fmt.Sprintf(":%s", port)
	log.Printf("Starting server on %s...", addr)
	if err := startServer(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func startServer(addr string) error {
	// Define a simple HTTP handler for testing.
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Server is running!"))
	})

	// Start the HTTP server.
	log.Printf("Server is running on %s", addr)
	return http.ListenAndServe(addr, nil)
}
