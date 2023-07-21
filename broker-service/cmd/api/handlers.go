package main

import (
	"broker/grpc/authentication"
	"context"
	"errors"
	"log"
	"net/http"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type RequestPayload struct {
	Action string      `json:"action"`
	Auth   AuthPayload `json:"auth,omitempty"`
	Log    LogPayload  `json:"log,omitempty"`
	Mail   MailPayload `json:"mail,omitempty"`
}

type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LogPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

type MailPayload struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Subject string `json:"subject"`
	Message string `json:"message"`
}

func (app *Config) Authenticate(w http.ResponseWriter, r *http.Request) {
	log.Println("Processing authentication request")

	var requestPayload RequestPayload

	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		log.Println("Error: ", err)
		app.errorJSON(w, errors.New("Error during authetication request"), http.StatusBadRequest)
		return
	}

	conn, err := grpc.Dial(app.Env["AUTHENTICATION_GRPC_SERVER"], grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		log.Println("Error connecting to gRPC server: ", err)
		app.errorJSON(w, errors.New("Error during authetication request"), http.StatusBadRequest)
		return
	}
	defer conn.Close()

	client := authentication.NewAuthenticationClient(conn)

	authResponse, err := client.AuthenticateWithEmailAndPassword(context.TODO(), &authentication.AuthRequest{
		Email:    requestPayload.Auth.Email,
		Password: requestPayload.Auth.Password,
	})
	if err != nil {
		log.Println("Error processing gRPC call: ", err)
		app.errorJSON(w, errors.New("Error during authetication request"), http.StatusBadRequest)
		return
	}

	if !authResponse.Success {
		log.Println("Authentication failed")
		app.errorJSON(w, errors.New("Authentication failed"), http.StatusUnauthorized)
		return
	}

	log.Println("Authentication success")
	var payload jsonResponse
	payload.Error = false
	payload.Message = "Login successful!"
	//payload.Data = jsonFromService.Data

	_ = app.pushToQueue("logs_topic", LogPayload{Name: "auth", Data: "Authentication success"})
	app.writeJSON(w, http.StatusOK, payload)
}

// Refactor to use amqp
func (app *Config) SendMail(w http.ResponseWriter, r *http.Request) {
	//TODO: get entry from request, marchal and pass to pushToQueue as string

	var payloadResponse jsonResponse
	payloadResponse.Error = false
	payloadResponse.Message = "To be implemented"

	app.writeJSON(w, http.StatusOK, payloadResponse)
}

func (app *Config) Log(w http.ResponseWriter, r *http.Request) {
	log.Println("Processing log request")

	var requestPayload RequestPayload
	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		log.Println("Error: ", err)
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	entry := requestPayload.Log
	err = app.pushToQueue("logs_topic", entry)
	if err != nil {
		log.Println("Error: ", err)
		app.errorJSON(w, err)
		return
	}

	var payloadResponse jsonResponse
	payloadResponse.Error = false
	payloadResponse.Message = "Logged via RabbitMQ"

	app.writeJSON(w, http.StatusOK, payloadResponse)
}
