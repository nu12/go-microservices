package main

import (
	"log"
	"net/http"
)

func (app *Config) Broker(w http.ResponseWriter, r *http.Request) {
	log.Printf("Processing information")

	payload := jsonResponse{
		Error:   false,
		Message: "Hit the Broker!",
	}

	_ = app.writeJSON(w, http.StatusAccepted, payload)

}
