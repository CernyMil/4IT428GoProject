package v1

import (
	"encoding/json"
	"io"
	"net/http"

	"subscriber-service/transport/api/v1/model"
	"subscriber-service/transport/util"
)

func (h *Handler) SendPublishedPost(w http.ResponseWriter, r *http.Request) {
	b, err := io.ReadAll(r.Body)
	if err != nil {
		util.WriteErrResponse(w, http.StatusBadRequest, err)
		return
	}

	var post model.Post

	if err := json.Unmarshal(b, &post); err != nil {
		util.WriteErrResponse(w, http.StatusBadRequest, err)
		return
	}

	if err := validate.Struct(post); err != nil {
		util.WriteErrResponse(w, http.StatusBadRequest, err)
		return
	}

	if err = h.service.SendPublishedPost(r.Context(), post); err != nil {
		util.WriteErrResponse(w, http.StatusInternalServerError, err)
		return
	}
	util.WriteResponse(w, http.StatusOK, post)
}
