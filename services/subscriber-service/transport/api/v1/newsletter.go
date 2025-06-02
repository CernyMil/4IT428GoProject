package v1

import (
	"encoding/json"
	"net/http"

	"subscriber-service/pkg/id"
	"subscriber-service/transport/util"
)

/*
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
*/

func (h *Handler) DeleteNewsletterSubscriptions(w http.ResponseWriter, r *http.Request) {
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

	err := h.service.DeleteNewsletterSubscriptions(r.Context(), newsletterID)
	if err != nil {
		util.WriteErrResponse(w, http.StatusInternalServerError, err)
		return
	}

	util.WriteResponse(w, http.StatusNoContent, newsletterIdStr)
}
