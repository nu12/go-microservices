package main

import (
	"fmt"
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
	log.Printf("Starting Broker service")
	app := Config{
		Env: map[string]string{},
	}
	app.loadEnv()
	app.setupRabbitMQ()
	defer app.Rabbit.Close()

	log.Printf("Starting Broker server on port %s\n", addr)

	server := http.Server{
		Addr:    addr,
		Handler: app.routes(),
	}
	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}

func (app *Config) loadEnv() {
	for _, v := range []string{"RABBITMQ_URL", "AUTHENTICATION_GRPC_SERVER"} {
		env, isSet := os.LookupEnv(v)
		if !isSet {
			log.Panicln(fmt.Sprintf("Cannot load environment variable %s", v))
		}
		app.Env[v] = env
	}
}

func (app *Config) setupRabbitMQ() {
	rabbitConn, err := app.connectRabbitMQ()
	if err != nil {
		log.Panicln("Cannot connect to RabbitMQ", err)
	}

	app.Rabbit = rabbitConn
	log.Println("Connected to RabbitMQ")
}

func (app *Config) connectRabbitMQ() (*amqp.Connection, error) {
	var counts int64
	var backOff = 5 * time.Second
	var connection *amqp.Connection

	for {
		log.Println("Dial ", counts)
		c, err := amqp.Dial(app.Env["RABBITMQ_URL"])
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
