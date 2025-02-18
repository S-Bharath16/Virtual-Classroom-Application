package main

import (
	"bytes"
	"fmt"
	"html/template"
	"log"

	"gopkg.in/gomail.v2"
)

// Struct for email template data
type EmailData struct {
	Name string
}

// Function to send an email
func sendEmail(to string, name string) error {
	tmpl, err := template.New("email").Parse(emailTemplate)
	if err != nil {
		return err
	}

	var body bytes.Buffer
	tmpl.Execute(&body, EmailData{Name: name})

	m := gomail.NewMessage()
	m.SetHeader("From", "your-email@example.com")
	m.SetHeader("To", to)
	m.SetHeader("Subject", "Welcome to Virtual Classroom!")
	m.SetBody("text/html", body.String())

	d := gomail.NewDialer("smtp.gmail.com", 587, "your-email@example.com", "your-email-password")

	if err := d.DialAndSend(m); err != nil {
		log.Println("Failed to send email:", err)
		return err
	}

	fmt.Println("Email sent successfully to:", to)
	return nil
}
