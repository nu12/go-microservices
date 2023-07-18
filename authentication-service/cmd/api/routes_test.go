package main

import (
	"net/http"
	"testing"

	"github.com/go-chi/chi/v5"
)

func TestRoutes(t *testing.T) {
	testRoutes := testApp.routes()
	chiRoutes := testRoutes.(chi.Router)

	mustExist := []string{"/authenticate"}

	for _, route := range mustExist {
		found := false
		_ = chi.Walk(chiRoutes, func(method, currentRoute string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
			if route == currentRoute {
				found = true
			}
			return nil
		})

		if !found {
			t.Errorf("Route %s not found", route)
		}
	}

}
