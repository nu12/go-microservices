package main

import (
	"log"
	"logger/data"
	"net/http"
)

type JSONPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func (app *Config) WriteLog(w http.ResponseWriter, r *http.Request) {

	var requestPayload JSONPayload

	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	event := data.LogEntry{
		Name: requestPayload.Name,
		Data: requestPayload.Data,
	}

	err = app.Models.Entry.Insert(event)

	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	resp := jsonResponse{
		Error:   false,
		Message: "logged",
	}

	log.Println("New log entry created by logger service")

	app.writeJSON(w, http.StatusOK, resp)

}
