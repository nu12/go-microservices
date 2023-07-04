package main

import (
	"log"
	"net/http"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

const addr = "0.0.0.0:8080"

type Config struct {
	Rabbit *amqp.Connection
}

func main() {

	rabbitConn, err := connect()
	if err != nil {
		log.Panicln("Cannot connect to RabbitMQ", err)
	}

	defer rabbitConn.Close()
	log.Println("Connected to RabbitMQ")

	app := Config{
		Rabbit: rabbitConn,
	}
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

func connect() (*amqp.Connection, error) {
	var counts int64
	var backOff = 5 * time.Second
	var connection *amqp.Connection

	for {
		log.Println("Dial ", counts)
		c, err := amqp.Dial("amqp://guest:guest@rabbitmq")
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
