package main

import (
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
)

type Config struct {
	RequiredEnv []string
	Env         map[string]string
}

func main() {
	c := Config{
		RequiredEnv: []string{"BACKEND_ADDRESS", "COMMIT"},
		Env:         map[string]string{},
	}
	err := c.loadRequiredEnv()
	if err != nil {
		log.Panicln(err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		render(w, "test.page.gohtml", &c)
	})

	fmt.Println("Starting front-end service on port 8080")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Panic(err)
	}
}

func (c *Config) loadRequiredEnv() error {
	for _, v := range c.RequiredEnv {
		envVal, isSet := os.LookupEnv(v)
		if !isSet {
			return errors.New(fmt.Sprintf("Environment variable %s not found", v))
		}
		c.Env[v] = envVal
	}
	return nil
}

func render(w http.ResponseWriter, t string, config *Config) {
	log.Println("Rendering template")

	partials := []string{
		"./cmd/web/templates/base.layout.gohtml",
		"./cmd/web/templates/header.partial.gohtml",
		"./cmd/web/templates/footer.partial.gohtml",
	}

	var templateSlice []string
	templateSlice = append(templateSlice, fmt.Sprintf("./cmd/web/templates/%s", t))

	for _, x := range partials {
		templateSlice = append(templateSlice, x)
	}

	tmpl, err := template.ParseFiles(templateSlice...)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, config.Env); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
