package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
)

type Config map[string]string

func main() {
	c := Config{
		"Address": "http://localhost:8081",
		"Commit":  os.Getenv("COMMIT"),
	}
	if remoteAddress, isSet := os.LookupEnv("ADDRESS"); isSet {
		c["Address"] = remoteAddress
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		render(w, "test.page.gohtml", &c)
	})

	fmt.Println("Starting front end service on port 8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Panic(err)
	}
}

func render(w http.ResponseWriter, t string, config *Config) {

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

	if err := tmpl.Execute(w, *config); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
