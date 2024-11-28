package smtp

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"log/slog"
	"net/smtp"
	"path/filepath"

	"github.com/MartynyukAlexey/gymshark/internal/config"
)

//go:embed templates
var templates embed.FS

type SMTPMailer struct {
	config *config.MailerConfig
	logger *slog.Logger
}

func NewSMTPMailer(cfg *config.MailerConfig, logger *slog.Logger) *SMTPMailer {
	return &SMTPMailer{
		config: cfg,
		logger: logger,
	}
}

func (m *SMTPMailer) sendEmail(to, subject, body string) error {
	addr := fmt.Sprintf("%s:%d", m.config.RelayHost, m.config.RelayPort)
	auth := smtp.PlainAuth("", m.config.SenderEmail, m.config.SenderPassword, m.config.RelayHost)

	message := fmt.Sprintf(
		"From: %s\r\n"+
			"To: %s\r\n"+
			"Subject: %s\r\n"+
			"MIME-version: 1.0\r\n"+
			"Content-Type: text/html; charset=\"UTF-8\"\r\n\r\n"+
			"%s",
		m.config.SenderEmail, to, subject, body)

	return smtp.SendMail(addr, auth, m.config.SenderEmail, []string{to}, []byte(message))
}

func (m *SMTPMailer) renderTemplate(templateName string, data map[string]string) (string, error) {
	t, err := template.ParseFS(templates, filepath.Join("templates", templateName), filepath.Join("templates", "base.html"))

	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func (m *SMTPMailer) SendActivationEmail(to string, activationCode string) error {
	subject := "Activate Your Account"

	body, err := m.renderTemplate("user_activation.html", map[string]string{
		"ActivationCode": activationCode,
	})

	if err != nil {
		return err
	}

	return m.sendEmail(to, subject, body)
}

func (m *SMTPMailer) SendPasswordResetEmail(to string, resetCode string) error {
	subject := "Reset Your Password"

	body, err := m.renderTemplate("password_reset.html", map[string]string{
		"ResetCode": resetCode,
	})

	if err != nil {
		return err
	}

	return m.sendEmail(to, subject, body)
}
