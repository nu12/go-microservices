package main

import (
	"log"
	"net/http"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

const addr = "0.0.0.0:8080"

type Config struct {
	Rabbit *amqp.Connection
	Env    map[string]string
}

func main() {

	app := Config{
		Env: map[string]string{
			"rabbitmq":     "amqp://guest:guest@rabbitmq",
			"authenticate": "http://authentication:8080",
			"logger":       "http://logger:8080",
			"mailer":       "http://mailer:8080",
			"loggerrpc":    "logger:5001",
			"loggergrpc":   "logger:50001",
		},
	}

	app.loadEnv()

	rabbitConn, err := app.connect()
	if err != nil {
		log.Panicln("Cannot connect to RabbitMQ", err)
	}

	app.Rabbit = rabbitConn

	defer rabbitConn.Close()
	log.Println("Connected to RabbitMQ")

	log.Printf("Starting Broker service on port %s\n", addr)

	server := http.Server{
		Addr:    addr,
		Handler: app.routes(),
	}
	err = server.ListenAndServe()
	if err != nil {
		panic(err)
	}

}

func (app *Config) loadEnv() {
	if s, isSet := os.LookupEnv("AUTHENTICATE_URL"); isSet {
		app.Env["authenticate"] = s
	}

	if s, isSet := os.LookupEnv("RABBITMQ_URL"); isSet {
		app.Env["rabbitmq"] = s
	}
	log.Println(app.Env["rabbitmq"])
	if s, isSet := os.LookupEnv("LOGGER_URL"); isSet {
		app.Env["logger"] = s
	}

	if s, isSet := os.LookupEnv("LOGGER_RPC_URL"); isSet {
		app.Env["loggerrpc"] = s
	}

	if s, isSet := os.LookupEnv("LOGGER_GRPC_URL"); isSet {
		app.Env["loggergrpc"] = s
	}

	if s, isSet := os.LookupEnv("MAILER_URL"); isSet {
		app.Env["mailer"] = s
	}

}

func (app *Config) connect() (*amqp.Connection, error) {
	var counts int64
	var backOff = 5 * time.Second
	var connection *amqp.Connection

	for {
		log.Println("Dial ", counts)
		c, err := amqp.Dial(app.Env["rabbitmq"])
		if err != nil {
			log.Println("Waiting RabbitMQ...")
			counts++
		} else {
			connection = c
			break
		}

		if counts > 10 {
			return nil, err
		}

		time.Sleep(backOff)
		continue
	}

	return connection, nil
}
