package v1

import (
	"net/http"

	"subscriber-api/transport/util"
)

func (h *Handler) DeleteNewsletter(w http.ResponseWriter, r *http.Request) {
	NewsletterID := getNewsletterId(w, r)

	err := h.service.DeleteNewsletter(r.Context(), NewsletterID)
	if err != nil {
		util.WriteErrResponse(w, http.StatusInternalServerError, err)
		return
	}

	util.WriteResponse(w, http.StatusOK, NewsletterID)
}
