package v1

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-playground/validator/v10"

	"subscriber-api/pkg/id"
	token "subscriber-api/pkg/token"
	svcmodel "subscriber-api/service/model"
	"subscriber-api/transport/util"
)

var validate = validator.New()

func (h *Handler) SubscribeToNewsletter(w http.ResponseWriter, r *http.Request) {
	var subReq svcmodel.SubscribeRequest
	var err error

	subReq.NewsletterID = getNewsletterId(w, r)

	subReq.Email = getEmail(w, r)

	if err := validate.Struct(subReq); err != nil {
		util.WriteErrResponse(w, http.StatusBadRequest, err)
		return
	}

	err = h.service.SubscribeToNewsletter(r.Context(), subReq)
	if err != nil {
		util.WriteErrResponse(w, http.StatusInternalServerError, err)
		return
	}

	util.WriteResponse(w, http.StatusOK, subReq)
}

func (h *Handler) ConfirmSubscription(w http.ResponseWriter, r *http.Request) {
	tokenString := r.URL.Query().Get("token")
	if tokenString == "" {
		util.WriteErrResponse(w, http.StatusBadRequest, errors.New("missing token"))
		return
	}

	email, err := token.DecryptToken(tokenString)
	if err != nil {
		util.WriteErrResponse(w, http.StatusBadRequest, errors.New("invalid or expired token"))
		return
	}

	subReq := svcmodel.SubscribeRequest{
		NewsletterID: getNewsletterId(w, r),
		Email:        email,
	}

	if err := validate.Struct(subReq); err != nil {
		util.WriteErrResponse(w, http.StatusBadRequest, err)
		return
	}

	if _, err = h.service.ConfirmSubscription(r.Context(), subReq); err != nil {
		util.WriteErrResponse(w, http.StatusInternalServerError, err)
		return
	}
}

func (h *Handler) UnsubscribeFromNewsletter(w http.ResponseWriter, r *http.Request) {
	var unsubReq svcmodel.UnsubscribeRequest
	unsubReq.NewsletterID = getNewsletterId(w, r)
	tokenString := getToken(w, r)

	subscriptionStr, err := token.DecryptToken(tokenString)
	if err != nil {
		util.WriteErrResponse(w, http.StatusBadRequest, errors.New("invalid or expired token"))
		return
	}
	if err := unsubReq.SubscriptionID.FromString(subscriptionStr); err != nil {
		util.WriteErrResponse(w, http.StatusBadRequest, errors.New("invalid subscription ID in token"))
		return
	}

	if err := h.service.UnsubscribeFromNewsletter(r.Context(), unsubReq); err != nil {
		util.WriteErrResponse(w, http.StatusInternalServerError, err)
		return
	}
}

func getEmail(w http.ResponseWriter, r *http.Request) string {
	var email string = r.URL.Query().Get("email")

	if email == "" {
		util.WriteErrResponse(w, http.StatusBadRequest, errors.New("email is required"))
		return ""
	}
	return email
}

func getNewsletterId(w http.ResponseWriter, r *http.Request) id.Newsletter {
	var newsletterID id.Newsletter
	if err := newsletterID.FromString(chi.URLParam(r, "id")); err != nil {
		http.Error(w, "invalid newsletter ID", http.StatusBadRequest)
	}
	return newsletterID
}

func getToken(w http.ResponseWriter, r *http.Request) string {
	token := r.URL.Query().Get("token")
	if token == "" {
		util.WriteErrResponse(w, http.StatusBadRequest, errors.New("missing token"))
		return ""
	}
	return token
}
