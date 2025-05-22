package v1

import (
	"encoding/json"
	"io"
	"net/http"

	"subscriber-api/pkg/id"
	svcmodel "subscriber-api/service/model"
	"subscriber-api/transport/util"

	"github.com/go-chi/chi"
)

func (h *Handler) SendPublishedPost(w http.ResponseWriter, r *http.Request) {
	b, err := io.ReadAll(r.Body)
	if err != nil {
		util.WriteErrResponse(w, http.StatusBadRequest, err)
		return
	}

	var post svcmodel.Post

	post.ID, err = getPostId(w, r)
	if err != nil {
		util.WriteErrResponse(w, http.StatusBadRequest, err)
		return
	}

	if err := json.Unmarshal(b, &post); err != nil {
		util.WriteErrResponse(w, http.StatusBadRequest, err)
		return
	}

	if err := validate.Struct(post); err != nil {
		util.WriteErrResponse(w, http.StatusBadRequest, err)
		return
	}

	post, err = h.service.SendPublishedPost(r.Context(), post)

	if err != nil {
		util.WriteErrResponse(w, http.StatusInternalServerError, err)
		return
	}

	util.WriteResponse(w, http.StatusOK, post)
}

func getPostId(w http.ResponseWriter, r *http.Request) (id.Post, error) {
	var postID id.Post
	if err := postID.FromString(chi.URLParam(r, "id")); err != nil {
		return id.Post{}, err
	}
	return postID, nil
}
