package v1

import (
	"encoding/json"
	"io"
	"net/http"

	"subscriber-service/pkg/id"
	"subscriber-service/transport/api/v1/model"
	"subscriber-service/transport/util"

	"github.com/go-chi/chi"
)

func (h *Handler) SendPublishedPost(w http.ResponseWriter, r *http.Request) {
	b, err := io.ReadAll(r.Body)
	if err != nil {
		util.WriteErrResponse(w, http.StatusBadRequest, err)
		return
	}

	var post model.Post

	post.ID = getPostId(w, r)

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

func getPostId(w http.ResponseWriter, r *http.Request) id.Post {
	var postID id.Post
	if err := postID.FromString(chi.URLParam(r, "id")); err != nil {
		util.WriteErrResponse(w, http.StatusBadRequest, err)
		return id.Post{}
	}
	return postID
}

/*
url := fmt.Sprintf("http://subscriber-service:8083/nginx/newsletters/{newsletterId}/posts/{postId}/publish", newsletterID)

    req, err := http.NewRequest(http.MethodGet, url, nil)
    if err != nil {
        return err
    }

    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return fmt.Errorf("subscriber-service returned status %d", resp.StatusCode)
    }
*/
