package main

import (
	"log"
	"net/http"
	"os"
	"strconv"
)

const webPort = "0.0.0.0:8080"

type Config struct {
	Mailer Mail
}

func main() {

	app := Config{
		Mailer: createMail(),
	}

	srv := http.Server{
		Addr:    webPort,
		Handler: app.routes(),
	}

	log.Println("Starting service")
	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}

func createMail() Mail {
	port, _ := strconv.Atoi(os.Getenv("MAIL_PORT"))
	m := Mail{
		Domain:      os.Getenv("MAIL_DOMAIN"),
		Host:        os.Getenv("MAIL_HOST"),
		Port:        port,
		Username:    os.Getenv("MAIL_USERNAME"),
		Password:    os.Getenv("MAIL_PASSWORD"),
		Encryption:  os.Getenv("MAIL_ENCRYPTION"),
		FromName:    os.Getenv("MAIL_FROM_NAME"),
		FromAddress: os.Getenv("MAIL_FROM_ADDRESS"),
	}

	return m
}
