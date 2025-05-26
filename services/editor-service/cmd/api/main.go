package main

import (
	"context"
	"editor-service/repository"
	"editor-service/service"
	transport "editor-service/transport/api/v1"
	"editor-service/transport/middleware"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	dbURL := os.Getenv("DATABASE_URL")
	firebaseCred := "firebase-cred.json"

	dbpool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}

	auth, err := middleware.NewFirebaseAuth(firebaseCred)
	if err != nil {
		log.Fatalf("Failed to initialize Firebase: %v", err)
	}

	repo := repository.NewPgxEditorRepository(dbpool)
	svc := service.NewEditorService(repo, auth)
	handler := transport.NewEditorHandler(svc)

	r := chi.NewRouter()
	//r.Use(middleware.Logger)

	r.Post("/signup", handler.SignUp)
	r.Post("/signin", handler.SignIn)
	r.Post("/change-password", handler.ChangePassword)

	log.Println("Server running on :8081")
	http.ListenAndServe(":8081", r)
}
