package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"newsletter-management-api/repository"

	"github.com/joho/godotenv" // For loading .env files
	_ "github.com/lib/pq"      // PostgreSQL driver
)

func main() {
	// Load environment variables from the .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Get the port from the environment variables
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT is not set in the environment variables")
	}

	// Get the database URL from the environment variables
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL is not set in the environment variables")
	}

	// Connect to the database
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	defer db.Close()

	// Test the database connection
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping the database: %v", err)
	}
	log.Println("Connected to the database successfully!")

	// Initialize the repository
	repo := repository.NewPostgresRepository(db)

	// Start the server
	addr := fmt.Sprintf(":%s", port)
	log.Printf("Starting server on %s...", addr)
	if err := startServer(addr, repo); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func startServer(addr string, repo repository.Repository) error {
	// Create and retrieve newsletters
	http.HandleFunc("/newsletters", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost: // Create a new newsletter
			var n repository.Newsletter
			if err := json.NewDecoder(r.Body).Decode(&n); err != nil {
				http.Error(w, "Invalid request payload", http.StatusBadRequest)
				return
			}
			n.CreatedAt = time.Now()
			if err := repo.Save(r.Context(), &n); err != nil {
				http.Error(w, "Failed to create newsletter", http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(n)

		case http.MethodGet: // Retrieve all newsletters
			newsletters, err := repo.FindAll(r.Context())
			if err != nil {
				http.Error(w, "Failed to retrieve newsletters", http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(newsletters)

		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Retrieve, update, and delete a specific newsletter by ID
	http.HandleFunc("/newsletters/", func(w http.ResponseWriter, r *http.Request) {
		id := strings.TrimPrefix(r.URL.Path, "/newsletters/")
		if id == "" {
			http.Error(w, "Newsletter ID is required", http.StatusBadRequest)
			return
		}

		switch r.Method {
		case http.MethodGet: // Retrieve a specific newsletter by ID
			newsletter, err := repo.FindByID(r.Context(), id)
			if err != nil {
				http.Error(w, "Newsletter not found", http.StatusNotFound)
				return
			}
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(newsletter)

		case http.MethodPut: // Update a specific newsletter
			var input repository.UpdateNewsletterInput
			if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
				http.Error(w, "Invalid request payload", http.StatusBadRequest)
				return
			}
			updatedNewsletter, err := repo.Update(r.Context(), id, input)
			if err != nil {
				http.Error(w, "Failed to update newsletter", http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(updatedNewsletter)

		case http.MethodDelete: // Delete a specific newsletter
			if err := repo.Delete(r.Context(), id); err != nil {
				http.Error(w, "Failed to delete newsletter", http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusNoContent)

		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Start the HTTP server
	log.Printf("Server is running on %s", addr)
	return http.ListenAndServe(addr, nil)
}
