package config

import (
	"crypto/tls"
	"fmt"
	"net/smtp"
	"os"
)

var SMTPClient *smtp.Client

func SMTPConnect() error {
	host := os.Getenv("SMTP_HOST")
	port := os.Getenv("SMTP_PORT")
	email := os.Getenv("SMTP_EMAIL")
	password := os.Getenv("SMTP_PASSWORD")

	fmt.Printf("Connecting to SMTP server at %s:%s with email %s\n", host, port, email)

	smtpAuth := smtp.PlainAuth("", email, password, host)

	// connect to smtp server
	client, err := smtp.Dial(host + ":" + port)
	if err != nil {
		fmt.Println("Failed to connect to SMTP server:", err)
		return err
	}

	SMTPClient = client
	client = nil

	// initiate TLS handshake
	if ok, _ := SMTPClient.Extension("STARTTLS"); ok {
		config := &tls.Config{ServerName: host}
		if err = SMTPClient.StartTLS(config); err != nil {
			fmt.Println("Failed to start TLS:", err)
			return err
		}
	}

	// authenticate
	err = SMTPClient.Auth(smtpAuth)
	if err != nil {
		fmt.Println("Failed to authenticate:", err)
		return err
	}

	return nil
}
