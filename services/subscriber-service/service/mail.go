package service

import (
	"bytes"
	"html/template"
	"os"
)

func SendConfirmationRequestMail(email string, newsletterId string) error {
	baseUrl := os.Getenv("BASE_URL")

	// Send confirmation email
	templateData := struct {
		ConfirmLink     string
		UnsubscribeLink string
	}{
		ConfirmLink:     baseUrl + "/api/v1/newsletters/" + newsletterId + "/confirm?email=" + email,
		UnsubscribeLink: baseUrl + "/api/v1/newsletters/" + newsletterId + "/unsubscribe?email=" + email,
	}

	t, err := template.ParseFiles("templates/confirm_request.html")

	if err != nil {
		return err
	}

	buf := new(bytes.Buffer)
	if err = t.Execute(buf, templateData); err != nil {
		return err
	}

	html := buf.String()

	errEmail := sendMail(email, "Confirm subscription | VŠE Newsletter", html)

	if errEmail != nil {
		return errEmail
	}

	return nil
}

func SendConfirmationMail(email string, newsletterId string) error {
	baseUrl := os.Getenv("BASE_URL")

	// Send confirmation email
	templateData := struct {
		SubscriberEmail string
		UnsubscribeLink string
	}{
		SubscriberEmail: email,
		UnsubscribeLink: baseUrl + "/api/v1/newsletters/" + newsletterId + "/unsubscribe?email=" + email,
	}

	t, err := template.ParseFiles("templates/confirmed.html")

	if err != nil {
		return err
	}

	buf := new(bytes.Buffer)
	if err = t.Execute(buf, templateData); err != nil {
		return err
	}

	html := buf.String()

	errEmail := sendMail(email, "Subscription confirmed | VŠE Newsletter", html)

	if errEmail != nil {
		return errEmail
	}

	return nil
}
