package mail

import (
	"bytes"
	"html/template"
	"os"

	"github.com/resend/resend-go/v2"
)

func SendMail(recepient []string, subject string, html string) error {
	ResendApiKey := os.Getenv("RESEND_API_KEY")
	client := resend.NewClient(ResendApiKey)

	params := &resend.SendEmailRequest{
		From:    os.Getenv("EMAIL_ADDRESS"),
		To:      recepient,
		Html:    html,
		Subject: subject,
	}

	_, err := client.Emails.Send(params)
	if err != nil {
		return err
	}
	return nil
}

func PrepareHTML(templatePath string, data interface{}) (string, error) {
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}
