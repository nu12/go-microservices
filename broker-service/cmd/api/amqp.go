package main

import (
	"context"
	"encoding/json"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func (app *Config) pushToQueue(topicName string, data interface{}) error {

	ch, err := app.Rabbit.Channel()
	if err != nil {
		log.Println("Error 1")
		return err
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		topicName, // name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		log.Println("Error 1")
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	j, err := json.MarshalIndent(&data, "", "\t")
	if err != nil {
		log.Println("Error 3")
		return err
	}
	err = ch.PublishWithContext(ctx,
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(j),
		})

	if err != nil {
		log.Println("Error 4")
		return err
	}
	log.Println("Message sent to queue. Topic: ", topicName)

	return nil

}
