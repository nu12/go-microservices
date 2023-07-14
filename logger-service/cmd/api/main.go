package main

import (
	"context"
	"log"
	"logger/data"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	webPort  = "0.0.0.0:8080"
	rpcPort  = "0.0.0.0:5001"
	gRpcPort = "0.0.0.0:50001"
)

var client *mongo.Client

type Config struct {
	Models data.Models
	Env    map[string]string
}

func main() {
	app := Config{
		Env: map[string]string{
			"mongo": os.Getenv("MONGO_URL"),
		},
	}
	mongoClient, err := app.connectToMongo()
	if err != nil {
		log.Panic(err)
	}
	client = mongoClient
	app.Models = data.New(client)

	// Close mongo connection
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			log.Panic(err)
		}
	}()

	err = rpc.Register(new(RPCServer))
	if err != nil {
		log.Panic(err)
	}
	log.Println("Starting RPC service")
	go app.rpcListen()
	log.Println("Starting gRPC service")
	go app.gRPCListen()

	log.Println("Starting logging service")
	srv := http.Server{
		Addr:    webPort,
		Handler: app.routes(),
	}

	err = srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}

}

func (app *Config) rpcListen() error {
	log.Println("Starting RPC server")
	listen, err := net.Listen("tcp", rpcPort)
	if err != nil {
		return err
	}

	defer listen.Close()

	for {
		rpcConn, err := listen.Accept()
		if err != nil {
			return err
		}

		go rpc.ServeConn(rpcConn)
	}
}

func (app *Config) connectToMongo() (*mongo.Client, error) {
	clientOptions := options.Client().ApplyURI(app.Env["mongo"])
	clientOptions.SetAuth(options.Credential{
		Username: os.Getenv("MONGO_USER"),
		Password: os.Getenv("MONGO_PASSWORD"),
	})

	c, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Println("Error connection to mongo DB:")
		return nil, err
	}
	return c, nil
}
