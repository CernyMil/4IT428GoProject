package mail

import (
	"bytes"
	"fmt"
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

func PrepareHTMLFromBytes(templateContent []byte, data interface{}) (string, error) {
	// Create a new template and parse the content
	tmpl, err := template.New("emailTemplate").Parse(string(templateContent))
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	// Execute the template with the provided data
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return buf.String(), nil
}
