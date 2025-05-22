package main

import (
	"context"
	"editor-service/repository"
	"editor-service/service"
	"editor-service/transport"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	dbURL := os.Getenv("DATABASE_URL")
	firebaseCred := os.Getenv("FIREBASE_CRED")

	dbpool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}

	auth, err := transport.NewFirebaseAuth(firebaseCred)
	if err != nil {
		log.Fatalf("Failed to initialize Firebase: %v", err)
	}

	repo := repository.NewEditorRepository(dbpool)
	svc := service.NewEditorService(repo, auth)
	handler := transport.NewEditorHandler(svc)

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Post("/signup", handler.SignUp)

	log.Println("Server running on :8080")
	http.ListenAndServe(":8080", r)
}
