package main

import (
	"broker/event"
	"broker/logs"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/rpc"
	"time"

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

// Marked for deletion
func (app *Config) Broker(w http.ResponseWriter, r *http.Request) {
	log.Printf("Processing information")

	payload := jsonResponse{
		Error:   false,
		Message: "Hit the Broker!",
	}

	_ = app.writeJSON(w, http.StatusAccepted, payload)

}

// Marked for deletion
func (app *Config) HandleSubmission(w http.ResponseWriter, r *http.Request) {
	var requestPayload RequestPayload
	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	switch requestPayload.Action {
	case "auth":
		app.authenticate(w, requestPayload.Auth)
	case "log":
		// app.logItem(w, requestPayload.Log)
		app.logEventViaRabbit(w, requestPayload.Log)
		//app.logItemViaRPC(w, requestPayload.Log)
	case "mail":
		app.sendMail(w, requestPayload.Mail)
	default:
		app.errorJSON(w, errors.New("Unknown action"))
	}

}

// Refactor to use gRPC
func (app *Config) authenticate(w http.ResponseWriter, a AuthPayload) {
	jsonData, err := json.MarshalIndent(a, "", "\t")
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	request, err := http.NewRequest("POST", app.Env["authenticate"]+"/authenticate", bytes.NewBuffer(jsonData))
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	client := http.Client{}
	response, err := client.Do(request)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	defer response.Body.Close()

	if response.StatusCode == http.StatusUnauthorized {
		println("Unauthorised")
		app.errorJSON(w, errors.New("Invalid credentials"), http.StatusUnauthorized)
		return
	}

	if response.StatusCode == http.StatusAccepted {
		println("Not accepted")
		app.errorJSON(w, errors.New("Error calling auth service"), http.StatusUnauthorized)
		return
	}

	var jsonFromService jsonResponse

	err = json.NewDecoder(response.Body).Decode(&jsonFromService)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	if jsonFromService.Error { // Error from the authentication service
		println(jsonFromService.Message)
		app.errorJSON(w, errors.New("Invalid credentials"))
		return
	}

	var payload jsonResponse
	payload.Error = false
	payload.Message = "Authenticated!"
	payload.Data = jsonFromService.Data

	app.writeJSON(w, http.StatusOK, payload)
}

// Unused since listener implementation, kept here for reference
func (app *Config) logItem(w http.ResponseWriter, entry LogPayload) {
	jsonData, err := json.MarshalIndent(entry, "", "\t")
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	logServiceURL := app.Env["logger"] + "/log"

	request, err := http.NewRequest("POST", logServiceURL, bytes.NewBuffer(jsonData))
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	request.Header.Set("Content-Type", "application/json")

	client := http.Client{}
	response, err := client.Do(request)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		app.errorJSON(w, errors.New("Logger service didn't accept request"))
		return
	}

	var payloadResponse jsonResponse
	payloadResponse.Error = false
	payloadResponse.Message = "logged"

	app.writeJSON(w, http.StatusOK, payloadResponse)

}

// Refactor to use gRPC
func (app *Config) sendMail(w http.ResponseWriter, m MailPayload) {
	jsonData, err := json.MarshalIndent(m, "", "\t")
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	request, err := http.NewRequest("POST", app.Env["mailer"]+"/send", bytes.NewBuffer(jsonData))
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	request.Header.Set("Content-Type", "application/json")

	client := http.Client{}
	response, err := client.Do(request)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		app.errorJSON(w, errors.New("Error calling mail server"))
		return
	}

	var payloadResponse jsonResponse
	payloadResponse.Error = false
	payloadResponse.Message = "Message sent to " + m.To

	app.writeJSON(w, http.StatusOK, payloadResponse)
}

func (app *Config) logEventViaRabbit(w http.ResponseWriter, entry LogPayload) {
	err := app.pushToQueue(entry.Name, entry.Data)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	var payloadResponse jsonResponse
	payloadResponse.Error = false
	payloadResponse.Message = "Logged via RabbitMQ"

	app.writeJSON(w, http.StatusOK, payloadResponse)
}

func (app *Config) pushToQueue(name, message string) error {
	emitter, err := event.NewEventEmitter(app.Rabbit)
	if err != nil {
		return err
	}

	payload := LogPayload{
		Name: name,
		Data: message,
	}

	j, _ := json.MarshalIndent(&payload, "", "\t")
	err = emitter.Push(string(j), "log.INFO")
	if err != nil {
		return err
	}
	return nil

}

// Marked for deletion
type RPCPayload struct {
	Name string
	Data string
}

// Marked for deletion
func (app *Config) logItemViaRPC(w http.ResponseWriter, payload LogPayload) {
	client, err := rpc.Dial("tcp", app.Env["loggerrpc"])
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	rpcPayload := RPCPayload{
		Name: payload.Name,
		Data: payload.Data,
	}

	var result string
	err = client.Call("RPCServer.LogInfo", rpcPayload, &result)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	var payloadResponse jsonResponse
	payloadResponse.Error = false
	payloadResponse.Message = "Logged via RPC"

	app.writeJSON(w, http.StatusOK, payloadResponse)
}

// Marked for deletion
func (app *Config) logItemViaGRPC(w http.ResponseWriter, r *http.Request) {

	var requestPayload RequestPayload

	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	conn, err := grpc.Dial(app.Env["loggergrpc"], grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	defer conn.Close()

	client := logs.NewLogServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	_, err = client.WriteLog(ctx, &logs.LogRequest{
		LogEntry: &logs.Log{
			Name: requestPayload.Log.Name,
			Data: requestPayload.Log.Data,
		},
	})
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	var payloadResponse jsonResponse
	payloadResponse.Error = false
	payloadResponse.Message = "Logged via gRPC"

	app.writeJSON(w, http.StatusOK, payloadResponse)
}
