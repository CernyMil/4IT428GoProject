package service

/*
import (
	"context"

	"subscriber-api/service/errors"
	"subscriber-api/service/model"

	"github.com/resend/resend-go/v2"
)

var post = map[string]model.Subscriber{}

func (Service) SendPublishedPost(_ context.Context, post model.Post) error {
	if _, exists := subcribers[subscriber.Email]; exists {
		return errors.ErrEmailAlreadySubscribed
	}

	subcribers[subscriber.Email] = subscriber

	return nil

	client := resend.NewClient(apiKey)

	params := &resend.SendEmailRequest{
		From:    "onboarding@resend.dev",
		To:      []string{subscriber[Subscriber.Email]},
		Subject: post[Post.Title],
		Html:    "<p>Congrats on sending your <strong>first email</strong>!</p>",
	}

	sent, err := client.Emails.Send(params)

}

package handlers

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"publishing-service/handlers/model"
	"publishing-service/handlers/sql"
	"publishing-service/utils"
	"strings"

	"github.com/go-chi/chi/v5"
	supa "github.com/nedpals/supabase-go"
	"github.com/resend/resend-go/v2"
)

func loadEmailTemplate(data map[string]string) (string, error) {
    fileData, err := fs.ReadFile(PostTemplate, "static/post_template.html")
    if err != nil {
        return "", err
    }

    var builder strings.Builder
    scanner := bufio.NewScanner(strings.NewReader(string(fileData)))
    for scanner.Scan() {
        line := scanner.Text()
        for key, value := range data {
            placeholder := fmt.Sprintf("{{%s}}", key)
            line = strings.ReplaceAll(line, placeholder, value)
        }
        builder.WriteString(line + "\n")
    }

    if err := scanner.Err(); err != nil {
        return "", err
    }

    return builder.String(), nil
}

func (hd *CustomHandler) GetPosts(w http.ResponseWriter, r *http.Request) {
	newsletterId := chi.URLParam(r, "id")

	if newsletterId == "" {
		handleError(w, "ID is required", nil, http.StatusBadRequest)
		return
	}

	posts, err := hd.Repository.ListPosts(r.Context(), newsletterId)
	if err != nil {
		handleError(w, "Failed to fetch posts", err, http.StatusInternalServerError)
		return
	}

	sendJSON(w, posts, http.StatusOK)
}

	cfg := utils.LoadConfig(".env")
	apiKey := cfg.ResendApiKey
	client := resend.NewClient(apiKey)
	unsubscribeUrl := fmt.Sprintf("%s/api/unsubscribe", cfg.ServerUrl)

	for _, subscriber := range subscribers {
		data := map[string]string{
			"title": post.Title,
			"content": post.Content,
			"unsubscribeUrl": unsubscribeUrl,
			"newsletterId":   newsletter.ID,
			"userId":         subscriber.ID,
		}

		mail, loadErr := loadEmailTemplate(data)

		if loadErr != nil {
			log.Printf("Failed load mail template: %v", err)
		}

        params := &resend.SendEmailRequest{
            From:    "newsletter@tapeer.cz",
            To:      []string{subscriber.Email},
            Subject: fmt.Sprintf("A new post in %s!", newsletter.Title),
            Html:    mail,
			Headers: map[string]string{
				"List-Unsubscribe": fmt.Sprintf("<%s?newsletterId=%s&userId=%s>", unsubscribeUrl, newsletter.ID, subscriber.ID),
				"List-Unsubscribe-Post": "List-Unsubscribe=One-Click",
			},
        }

        _, err := client.Emails.Send(params)
        if err != nil {
            log.Printf("Failed to send email to %s: %v", subscriber.Email, err)
        }
    }

	sendJSON(w, post, http.StatusOK)

*/
