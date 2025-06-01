package v1

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"subscriber-service/pkg/id"
	"subscriber-service/transport/util"
)

func (h *Handler) CreateNewsletter(w http.ResponseWriter, r *http.Request) {
	var newsletterIdStr string

	if err := json.NewDecoder(r.Body).Decode(&newsletterIdStr); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if newsletterIdStr == "" {
		http.Error(w, "newsletterId is empty", http.StatusBadRequest)
		return
	}

	var newsletterID id.Newsletter
	if err := newsletterID.FromString(newsletterIdStr); err != nil {
		http.Error(w, "invalid newsletter ID format", http.StatusBadRequest)
		return
	}

	err := h.service.CreateNewsletter(r.Context(), newsletterID)
	if err != nil {
		util.WriteErrResponse(w, http.StatusInternalServerError, err)
		return
	}

	util.WriteResponse(w, http.StatusOK, newsletterIdStr)
}

func (h *Handler) DeleteNewsletter(w http.ResponseWriter, r *http.Request) {
	var newsletterIdStr string

	if err := json.NewDecoder(r.Body).Decode(&newsletterIdStr); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if newsletterIdStr == "" {
		http.Error(w, "newsletterId is empty", http.StatusBadRequest)
		return
	}

	var newsletterID id.Newsletter
	if err := newsletterID.FromString(newsletterIdStr); err != nil {
		http.Error(w, "invalid newsletter ID format", http.StatusBadRequest)
		return
	}

	err := h.service.DeleteNewsletter(r.Context(), newsletterID)
	if err != nil {
		util.WriteErrResponse(w, http.StatusInternalServerError, err)
		return
	}

	util.WriteResponse(w, http.StatusNoContent, newsletterIdStr)
}

// move to newsletter service
func notifySubscriberDeleteNewsletter(newsletterID string) error {
	url := "http://subscriber-service:8083/internal/delete"

	// Prepare JSON body
	payload, err := json.Marshal(map[string]string{"id": newsletterID})
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodDelete, url, bytes.NewBuffer(payload))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Service-Token", os.Getenv("SERVICE_TOKEN"))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("subscriber-service returned status %d", resp.StatusCode)
	}
	return nil
}
