package util

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
)

// SetServerLogLevel sets the log level for the server.
func SetServerLogLevel(level slog.Level) {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: level})))
}

// NewServerLogger creates a new logger for the server.
func NewServerLogger(name string) *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{AddSource: true}))
}

// WriteResponse writes a JSON response to the http.ResponseWriter.
func WriteErrResponse(w http.ResponseWriter, statusCode int, err error) {
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
}

// WriteResponse writes a JSON response with the given status code and data.
func WriteResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
