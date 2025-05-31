package v1

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

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

// move to newsletter service
func notifySubscriberSendPublishedPost(newsletterID, postID string, post model.Post) error {
	url := fmt.Sprintf("http://subscriber-service:8083/newsletters/%s/posts/%s/publish", newsletterID, postID)

	payload, err := json.Marshal(post)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(payload))
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
