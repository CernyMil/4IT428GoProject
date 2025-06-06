package main

import (
	"context"
	"editor-service/repository"
	"editor-service/service"
	api "editor-service/transport/api"
	v1 "editor-service/transport/api/v1"
	"editor-service/transport/middleware"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	dbURL := os.Getenv("DATABASE_URL")
	firebaseCred := os.Getenv("FIREBASE_CRED")
	firebaseAPIKey := os.Getenv("FIREBASE_API_KEY")
	dbpool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}

	auth, err := middleware.NewFirebaseAuth(firebaseCred, firebaseAPIKey)
	if err != nil {
		log.Fatalf("Failed to initialize Firebase: %v", err)
	}

	repo := repository.NewPgxEditorRepository(dbpool)
	svc := service.NewEditorService(repo, auth)
	handler := v1.NewEditorHandler(svc)
	controller := api.NewController(handler)
	r := chi.NewRouter()

	r.Use(chiMiddleware.Logger)
	r.Use(chiMiddleware.Recoverer)

	r.Post("/signup", controller.SignUp)
	r.Post("/signin", controller.SignIn)
	r.Post("/change-password", controller.ChangePassword)

	log.Println("Server running on :8081")
	http.ListenAndServe(":8081", r)
}
