package v1

import (
	"encoding/json"
	"net/http"

	"subscriber-service/pkg/id"
	"subscriber-service/transport/util"
)

func (h *Handler) CreateNewsletter(w http.ResponseWriter, r *http.Request) {
	var newsletterId string
	if err := json.NewDecoder(r.Body).Decode(&newsletterId); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if newsletterId == "" {
		http.Error(w, "newsletterId is empty", http.StatusBadRequest)
		return
	}

	var newsletterID id.Newsletter
	if err := newsletterID.FromString(newsletterId); err != nil {
		http.Error(w, "invalid newsletter ID format", http.StatusBadRequest)
		return
	}

	err := h.service.CreateNewsletter(r.Context(), newsletterID)
	if err != nil {
		util.WriteErrResponse(w, http.StatusInternalServerError, err)
		return
	}

	util.WriteResponse(w, http.StatusOK, newsletterId)
}

/*
   url := "http://subscriber-service:8083/nginx/newsletters"
   payload, err := json.Marshal(newsletterId)
   if err != nil {
       return err
   }

   resp, err := http.Post(url, "application/json", bytes.NewBuffer(payload))
   if err != nil {
       return err
   }
   defer resp.Body.Close()

   if resp.StatusCode != http.StatusOK {
       return fmt.Errorf("subscriber-service returned status %d", resp.StatusCode)
   }
*/

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
url := fmt.Sprintf("http://subscriber-service:8083/nginx/newsletters/%s/delete", newsletterID)

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
