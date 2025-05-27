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

/*
url := fmt.Sprintf("http://subscriber-service:8083/nginx/newsletters/{newsletterId}/delete", newsletterID)

    req, err := http.NewRequest(http.MethodDelete, url, nil)
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
