package mailer

import (
	"Backend/config"
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/smtp"
)

func SendMail(emailID []string, mailSubject string, mailBody string, category string) error {
	cfg := config.GetConfig()

	var templateFile string
	switch category {
	case "announcement":
		templateFile = "C:/Users/Dell/OneDrive/Desktop/6th Semester/Virtual-Classroom-Application/Backend/utilities/templates/announcementMailTemplate.html"
	}

	// Load email template
	tmpl, err := template.ParseFiles(templateFile)
	if err != nil {
		log.Println("Failed to load email template:", err)
		return fmt.Errorf("failed to load email template: %w", err)
	}

	// Generate email content
	var bodyBuffer bytes.Buffer
	err = tmpl.Execute(&bodyBuffer, struct {
		Subject string
		Body    string
	}{
		Subject: mailSubject,
		Body:    mailBody,
	})
	if err != nil {
		log.Println("Failed to execute email template:", err)
		return fmt.Errorf("failed to generate email content: %w", err)
	}

	emailAuth := smtp.PlainAuth("", cfg.EmailID, cfg.MailerPassword, cfg.MailHost)

	emailMessage := fmt.Sprintf("From: %s\r\n"+
		"To: %s\r\n"+
		"Subject: %s\r\n"+
		"MIME-version: 1.0;\r\n"+
		"Content-Type: text/html; charset=\"UTF-8\";\r\n"+
		"\r\n"+
		"%s", cfg.EmailID, emailID[0], mailSubject, bodyBuffer.String())

	err = smtp.SendMail(cfg.MailHostAddress, emailAuth, cfg.EmailID, emailID, []byte(emailMessage))
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}