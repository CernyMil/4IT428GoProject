package v1

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-playground/validator/v10"

	svcmodel "subscriber-service/service/model"
	"subscriber-service/transport/util"
)

var validate = validator.New()

func (h *Handler) SubscribeToNewsletter(w http.ResponseWriter, r *http.Request) {
	var subReq svcmodel.SubscribeRequest
	if err := json.NewDecoder(r.Body).Decode(&subReq); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if err := validate.Struct(subReq); err != nil {
		util.WriteErrResponse(w, http.StatusBadRequest, err)
		return
	}

	err := h.service.SubscribeToNewsletter(r.Context(), subReq)
	if err != nil {
		util.WriteErrResponse(w, http.StatusInternalServerError, err)
		return
	}

	util.WriteResponse(w, http.StatusOK, subReq)
}

func (h *Handler) ConfirmSubscription(w http.ResponseWriter, r *http.Request) {
	token := getToken(w, r)

	subscription, err := h.service.ConfirmSubscription(r.Context(), token)
	if err != nil {
		util.WriteErrResponse(w, http.StatusInternalServerError, err)
		return
	}

	util.WriteResponse(w, http.StatusOK, subscription)
}

func (h *Handler) UnsubscribeFromNewsletter(w http.ResponseWriter, r *http.Request) {
	token := getToken(w, r)

	if err := h.service.UnsubscribeFromNewsletter(r.Context(), token); err != nil {
		util.WriteErrResponse(w, http.StatusInternalServerError, err)
		return
	}
}

func getToken(w http.ResponseWriter, r *http.Request) string {
	token := r.URL.Query().Get("token")
	if token == "" {
		util.WriteErrResponse(w, http.StatusBadRequest, errors.New("missing token"))
		return ""
	}
	return token
}
