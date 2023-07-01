package main

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

type jsonResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"` // interface{}
}

func (app *Config) readJSON(w http.ResponseWriter, r *http.Request, data any) error {
	maxBytes := 1048576 // 1Mb

	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	dec := json.NewDecoder(r.Body)
	err := dec.Decode(data)
	if err != nil {
		return err
	}

	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("Must contain only one JSON value")
	}

	return nil

}

func (app *Config) writeJSON(w http.ResponseWriter, status int, data any, header ...http.Header) error {
	j, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(j)

	return nil
}

func (app *Config) errorError(w http.ResponseWriter, err error, status ...int) error {
	statusCode := http.StatusBadRequest
	if len(status) > 0 {
		statusCode = status[0]
	}

	payload := jsonResponse{
		Error:   true,
		Message: err.Error(),
	}

	return app.writeJSON(w, statusCode, payload)
}
