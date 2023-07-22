package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"logger/data"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	amqp "github.com/rabbitmq/amqp091-go"
)

var client *mongo.Client

type Config struct {
	Models           data.Models
	MongoClient      *mongo.Client
	RabbitConnection *amqp.Connection
	RabbitChannel    *amqp.Channel
	Env              map[string]string
}

func main() {
	app := Config{
		Env: map[string]string{},
	}
	app.setupEnv()

	err := app.setupMongoDB()
	if err != nil {
		log.Panic(err)
	}
	// Close mongo connection
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	defer func() {
		if err = app.MongoClient.Disconnect(ctx); err != nil {
			log.Panic(err)
		}
	}()

	msgStream, err := app.setupRabbitMQ()
	if err != nil {
		log.Panic(err)
	}
	defer app.RabbitConnection.Close()
	defer app.RabbitChannel.Close()

	var forever chan struct{}
	go func() {
		for d := range msgStream {
			var entry data.LogEntry
			err = json.Unmarshal(d.Body, &entry)
			if err != nil {
				log.Println("Error: ", err)
				return
			}
			app.WriteLog(entry)
		}
	}()

	log.Println("Waiting for messages")
	os.WriteFile("/tmp/readyz", []byte(""), 0644)
	<-forever
}

func (app *Config) setupEnv() {
	for _, env := range []string{"RABBITMQ_URL", "MONGO_URL", "MONGO_USER", "MONGO_PASSWORD"} {
		val, isSet := os.LookupEnv(env)
		if !isSet {
			log.Panicln(fmt.Sprintf("Variable %s not found", env))
		}
		app.Env[env] = val
	}
}

func (app *Config) setupMongoDB() error {
	mongoClient, err := app.connectToMongo()
	if err != nil {
		return err
	}
	client = mongoClient
	app.Models = data.New(client)
	app.MongoClient = client

	return nil
}

func (app *Config) connectToMongo() (*mongo.Client, error) {
	clientOptions := options.Client().ApplyURI(app.Env["MONGO_URL"])
	clientOptions.SetAuth(options.Credential{
		Username: app.Env["MONGO_USER"],
		Password: app.Env["MONGO_PASSWORD"],
	})

	c, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Println("Error connection to mongo DB:")
		return nil, err
	}
	return c, nil
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
		"logs_topic", // name
		false,        // durable
		false,        // delete when unused
		false,        // exclusive
		false,        // no-wait
		nil,          // arguments
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
