package v1

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-playground/validator/v10"

	"subscriber-service/pkg/id"
	token "subscriber-service/pkg/token"
	svcmodel "subscriber-service/service/model"
	"subscriber-service/transport/util"
)

var validate = validator.New()

func (h *Handler) SubscribeToNewsletter(w http.ResponseWriter, r *http.Request) {
	var subReq svcmodel.SubscribeRequest

	subReq.NewsletterID = getNewsletterId(w, r)
	subReq.Email = getEmail(w, r)

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
	tokenString := r.URL.Query().Get("token")
	if tokenString == "" {
		util.WriteErrResponse(w, http.StatusBadRequest, errors.New("missing token"))
		return
	}

	claims, err := token.ParseJWT(tokenString)
	if err != nil {
		util.WriteErrResponse(w, http.StatusBadRequest, errors.New("invalid or expired token"))
		return
	}

	email, ok := claims["email"].(string)
	if !ok {
		util.WriteErrResponse(w, http.StatusBadRequest, errors.New("invalid email in token claims"))
		return
	}

	newsletterId, ok := claims["newsletterId"].(string)
	if !ok {
		util.WriteErrResponse(w, http.StatusBadRequest, errors.New("invalid newsletterId in token claims"))
		return
	}
	var newsletterID id.Newsletter
	newsletterID.FromString(newsletterId)

	subReq := svcmodel.SubscribeRequest{
		NewsletterID: newsletterID,
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

	claims, err := token.ParseJWT(tokenString)
	if err != nil {
		util.WriteErrResponse(w, http.StatusBadRequest, errors.New("invalid or expired token"))
		return
	}

	subcscriptionId, ok := claims["subscriptionId"].(string)
	if !ok {
		util.WriteErrResponse(w, http.StatusBadRequest, errors.New("invalid subscriptionId in token claims"))
		return
	}
	var subscriptionID id.Subscription
	if err := subscriptionID.FromString(subcscriptionId); err != nil {
		util.WriteErrResponse(w, http.StatusBadRequest, errors.New("invalid subscriptionId"))
		return
	}
	unsubReq.SubscriptionID = subscriptionID

	if err := validate.Struct(unsubReq); err != nil {
		util.WriteErrResponse(w, http.StatusBadRequest, err)
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
	newsletterIdStr := chi.URLParam(r, "newsletterId")
	if newsletterIdStr == "" {
		http.Error(w, "missing newsletter ID", http.StatusBadRequest)
		return id.Newsletter{}
	}

	var newsletterID id.Newsletter
	if err := newsletterID.FromString(newsletterIdStr); err != nil {
		util.WriteErrResponse(w, http.StatusBadRequest, errors.New("invalid newsletter ID"))
		return id.Newsletter{}
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
