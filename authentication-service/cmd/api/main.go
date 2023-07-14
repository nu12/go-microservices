package main

import (
	"authentication/data"
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

const webPort = "0.0.0.0:8080"

// TODO: Create a struct for the Environment variables
type Config struct {
	DB     *sql.DB
	Models data.Models
	Env    map[string]string
}

func main() {
	log.Println("Starting authentication service")

	conn := connectToDB()
	if conn == nil {
		panic("Cannot connect to the database")
	}

	app := Config{
		DB:     conn,
		Models: data.New(conn),
		Env:    map[string]string{"logger": "http://logger:8080/log"},
	}

	if logger, isSet := os.LookupEnv("LOGGER_URL"); isSet {
		app.Env["logger"] = logger
	}

	srv := http.Server{
		Addr:    webPort,
		Handler: app.routes(),
	}

	err := srv.ListenAndServe()
	if err != nil {
		panic(err)
	}

}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil

}

func connectToDB() *sql.DB {
	dsn := os.Getenv("DSN")

	// TODO: Create a stopping condition for the loop

	for {
		connection, err := openDB(dsn)
		if err != nil {
			println("Database is not ready...")
		} else {
			println("Database is ready...")
			return connection
		}

		time.Sleep(2 * time.Second)
	}
	//return nil
}
