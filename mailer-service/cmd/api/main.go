package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Config struct {
	Mailer           Mail
	RabbitConnection *amqp.Connection
	RabbitChannel    *amqp.Channel
}

func main() {
	log.Println("Starting mailer service")

	for _, env := range []string{"RABBITMQ_URL", "MAIL_PORT", "MAIL_DOMAIN", "MAIL_HOST", "MAIL_USERNAME", "MAIL_PASSWORD", "MAIL_ENCRYPTION", "MAIL_FROM_NAME", "MAIL_FROM_ADDRESS"} {
		if _, isSet := os.LookupEnv(env); !isSet {
			log.Panicln(fmt.Sprintf("Variable %s not found", env))
		}
	}

	app := Config{
		Mailer: createMail(),
	}

	msgStream, err := app.setupRabbitMQ()
	if err != nil {
		log.Panic(err)
	}
	defer app.RabbitConnection.Close()
	defer app.RabbitChannel.Close()

	var forever chan struct{}
	go func() {
		for d := range msgStream {
			var entry MailMessage
			err = json.Unmarshal(d.Body, &entry)
			if err != nil {
				log.Println("Error: ", err)
				return
			}
			err = app.SendMail(entry)
			if err != nil {
				log.Println("Error: ", err)
			}
		}
	}()

	log.Println("Waiting for messages")
	<-forever

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

func (app *Config) setupRabbitMQ() (<-chan amqp.Delivery, error) {
	conn, err := app.connectRabbitMQ()
	if err != nil {
		log.Panic(err)
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Panic(err)
	}

	q, err := ch.QueueDeclare(
		"mails_topic", // name
		false,         // durable
		false,         // delete when unused
		false,         // exclusive
		false,         // no-wait
		nil,           // arguments
	)
	if err != nil {
		log.Panic(err)
	}

	return ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
}

func (app *Config) connectRabbitMQ() (*amqp.Connection, error) {
	var counts int64
	var backOff = 5 * time.Second
	var connection *amqp.Connection

	for {
		log.Println("Dial ", counts)
		c, err := amqp.Dial(os.Getenv("RABBITMQ_URL"))
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
