package v1

import (
	"net/http"

	"subscriber-service/transport/util"
)

func (h *Handler) DeleteNewsletter(w http.ResponseWriter, r *http.Request) {
	newsletterID := getNewsletterId(w, r)

	err := h.service.DeleteNewsletter(r.Context(), newsletterID)
	if err != nil {
		util.WriteErrResponse(w, http.StatusInternalServerError, err)
		return
	}

	util.WriteResponse(w, http.StatusOK, newsletterID)
}
