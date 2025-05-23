package service

import (
	"context"

	"github.com/google/uuid"

	"subscriber-api/pkg/id"
	svcmodel "subscriber-api/service/model"
)

func (s Service) SubscribeToNewsletter(ctx context.Context, subReq svcmodel.SubscribeRequest) error {
	/*
		subscriptionId := id.Subscription(uuid.New())
			errEmail := mail.SendConfirmationRequestMail(email, newsletterId.String())

			if errEmail != nil {
				return nil, errEmail
			}

			return err
	*/
	return nil // temporary TBD
}

func (s Service) ConfirmSubscription(ctx context.Context, subReq svcmodel.SubscribeRequest) (svcmodel.Subscription, error) {
	subscriptionId := id.Subscription(uuid.New())
	subscription, err := s.repository.AddSubscription(ctx, subReq.NewsletterID, subscriptionId, subReq.Email)
	if err != nil {
		return svcmodel.Subscription{}, err
	}
	/*  TBD
	errEmail := mail.SendConfirmationMail(email, newsletterId.String())

	if errEmail != nil {
		return nil, errEmail
	}
	*/
	return *subscription, err
}

func (s Service) UnsubscribeFromNewsletter(ctx context.Context, newsletterId id.Newsletter, subscriptionId id.Subscription) error {
	if err := s.repository.DeleteSubscription(ctx, newsletterId, subscriptionId); err != nil {
		return err
	}
	return nil
}

func (s Service) DeleteNewsletter(ctx context.Context, newsletterId id.Newsletter) error {
	err := s.repository.DeleteNewsletter(ctx, newsletterId)
	if err != nil {
		return err
	}
	return err
}

/*
func loadEmailTemplate(data map[string]string) (string, error) {
	fileData, err := fs.ReadFile(SubscribedMailTemplate, "static/mail_template.html")
	if err != nil {
		return "", err
	}
;
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


func (hd *CustomHandler) Subscribe(w http.ResponseWriter, r *http.Request) {
	token := utils.GetBearerToken(r)
	userId, _ := utils.ExtractSubFromToken(token)

	newsletterId := chi.URLParam(r, "id")
	if newsletterId == "" {
		handleError(w, "ID is required", nil, http.StatusBadRequest)
		return
	}

	subscription, subscriber, newsletter, err := hd.Repository.Subscribe(r.Context(), newsletterId, userId)
	if err != nil {
		if errors.Is(err, utils.ErrSubscriptionExists) {
			handleError(w, err.Error(), nil, http.StatusConflict) // 409 Conflict
		} else {
			handleError(w, "failed to subscribe", err, http.StatusInternalServerError)
		}
		return
	}

	cfg := utils.LoadConfig(".env")
	apiKey := cfg.ResendApiKey
	client := resend.NewClient(apiKey)
	unsubscribeUrl := fmt.Sprintf("%s/api/unsubscribe", cfg.ServerUrl)

	data := map[string]string{
		"unsubscribeUrl": unsubscribeUrl,
		"newsletterId":   newsletterId,
		"userId":         userId,
	}

	mail, loadErr := loadEmailTemplate(data)

	if loadErr != nil {
		log.Printf("Failed load mail template: %v", err)
	}

	params := &resend.SendEmailRequest{
		From:    "newsletter@tapeer.cz",
		To:      []string{subscriber.Email},
		Subject: fmt.Sprintf("You have been subscribed to %s.", newsletter.Title),
		Html:    mail,
		Headers: map[string]string{
			"List-Unsubscribe":      fmt.Sprintf("<%s?newsletterId=%s&userId=%s>", unsubscribeUrl, newsletterId, userId),
			"List-Unsubscribe-Post": "List-Unsubscribe=One-Click",
		},
	}

	_, emailErr := client.Emails.Send(params)
	if emailErr != nil {
		log.Printf("Failed to send email to %s: %v", subscriber.Email, err)
	}

	sendJSON(w, subscription, http.StatusOK)
}

*/
