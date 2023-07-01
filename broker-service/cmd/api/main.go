package main

import (
	"log"
	"net/http"
)

const addr = "0.0.0.0:8080"

type Config struct {
}

func main() {

	app := Config{}
	log.Printf("Starting Broker service on port %s\n", addr)

	server := http.Server{
		Addr:    addr,
		Handler: app.routes(),
	}
	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}

}
