package mailer

import (
	"fmt"
	"net/smtp"
	"Backend/config"
)

func SendMail(emailID []string, mailSubject string, mailBody string) error {

	cfg := config.GetConfig()

	emailAuth := smtp.PlainAuth(
		"",
		cfg.EmailID,
		cfg.MailerPassword,
		cfg.MailHost,
	)

	emailMessage := fmt.Sprintf("From: %s\r\n"+
		"To: %s\r\n"+
		"Subject: %s\r\n"+
		"\r\n"+
		"%s", cfg.EmailID, emailID[0], mailSubject, mailBody)

	err := smtp.SendMail(
		cfg.MailHostAddress,
		emailAuth,
		cfg.EmailID,
		emailID,
		[]byte(emailMessage),
	)

	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}
	
	return nil
}