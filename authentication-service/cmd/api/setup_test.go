package main

import (
	"authentication/data"
	"net/http"
	"os"
	"testing"
)

var testApp Config

// To mock http client
type RoundTripFunc func(req *http.Request) *http.Response

func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

func NewTestClient(fn RoundTripFunc) *http.Client {
	return &http.Client{
		Transport: fn,
	}
}

func TestMain(m *testing.M) {

	testApp.Repo = data.NewPostgresTestRepository()
	testApp.Client = *NewTestClient(func(req *http.Request) *http.Response {
		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     make(http.Header),
		}
	})

	os.Exit(m.Run())
}
