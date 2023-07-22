package main

import (
	"log"
	"logger/data"
)

type JSONPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func (app *Config) WriteLog(event data.LogEntry) {
	log.Println("Processing new log entry from the queue")

	err := app.Models.Entry.Insert(event)

	if err != nil {
		log.Println("Error processing log: ", err)
		return
	}

	log.Println("New log entry created by logger service")
}
