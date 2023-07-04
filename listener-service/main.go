package main

import (
	"listener/event"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	rabbitConn, err := connect()
	if err != nil {
		log.Panicln("Cannot connect to RabbitMQ", err)
	}

	defer rabbitConn.Close()
	log.Println("Connected to RabbitMQ")

	consumer, err := event.NewConsumer(rabbitConn)
	if err != nil {
		log.Panicln(err)
	}

	err = consumer.Listen([]string{"log.INFO", "log.WARNING", "log.ERROR"})
	if err != nil {
		log.Panicln(err)
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
