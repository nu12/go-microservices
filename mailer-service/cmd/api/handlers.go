package main

import "log"

type MailMessage struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Subject string `json:"subject"`
	Message string `json:"message"`
}

func (app *Config) SendMail(data MailMessage) error {
	log.Println("Sending email")

	msg := Message{
		From:    data.From,
		To:      data.To,
		Subject: data.Subject,
		Data:    data.Message,
	}

	err := app.Mailer.SendSMTPMessage(msg)
	if err != nil {
		return err
	}

	log.Println("Email sent successfully")
	return nil
}
