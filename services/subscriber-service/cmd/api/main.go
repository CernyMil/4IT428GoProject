package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"

	"subscriber-service/repository"
	"subscriber-service/service"
	"subscriber-service/transport/api"
	"subscriber-service/transport/util"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	httpx "go.strv.io/net/http"
	"google.golang.org/api/option"
)

var version = "v0.0.0"

func main() {
	ctx := context.Background()
	cfg := MustLoadConfig()
	util.SetServerLogLevel(slog.LevelInfo)

	firestoreClient, err := initializeFirebase(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to the firestore: %v", err)
	}
	defer firestoreClient.Close()

	firestoreRepo, err := repository.NewFirestoreRepository(firestoreClient)
	if err != nil {
		slog.Error("initializing repository", slog.Any("error", err))
	}

	controller, err := setupController(
		cfg,
		firestoreRepo,
	)
	if err != nil {
		slog.Error("initializing controller", slog.Any("error", err))
	}

	addr := fmt.Sprintf(":%d", cfg.Port)
	// Initialize the server config.
	serverConfig := httpx.ServerConfig{
		Addr:    addr,
		Handler: controller,
		Hooks:   httpx.ServerHooks{
			// BeforeShutdown: []httpx.ServerHookFunc{
			// 	func(_ context.Context) {
			// 		database.Close()
			// 	},
			// },
		},
		Limits: nil,
		Logger: util.NewServerLogger("httpx.Server"),
	}
	server := httpx.NewServer(&serverConfig)

	slog.Info("starting server", slog.Int("port", cfg.Port))
	if err := server.Run(ctx); err != nil {
		slog.Error("server failed", slog.Any("error", err))
	}
}

func initializeFirebase(cfg Config) (*firestore.Client, error) {
	// Use a service account

	ctx := context.Background()
	sa := option.WithCredentialsFile(cfg.ServiceAccount)
	app, err := firebase.NewApp(ctx, nil, sa)
	if err != nil {
		return nil, err
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalf("Failed to initialize Firestore: %v", err)
	}

	return client, nil
	/*
		ctx := context.Background()
		ProjectID := os.Getenv("FIREBASE_PROJECT_ID")
		conf := &firebase.Config{ProjectID: ProjectID}
		app, err := firebase.NewApp(ctx, conf)
		if err != nil {log.Fatalln(err)}
		client, err := app.Firestore(ctx)
		if err != nil {log.Fatalln(err)}
		return client, nil
	*/
}

func setupController(cfg Config, repository service.Repository) (*api.Controller, error) {
	// Initialize the service.
	svc, err := service.NewService(repository)
	if err != nil {
		return nil, fmt.Errorf("initializing user service: %w", err)
	}

	// Initialize the controller.
	controller, err := api.NewController(svc, version)
	if err != nil {
		return nil, fmt.Errorf("initializing controller: %w", err)
	}
	return controller, nil
}
