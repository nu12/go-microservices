package main

import (
	"authentication/data"
	authentication "authentication/grpc"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
	"google.golang.org/grpc"
)

type Config struct {
	Repo data.Repository
	Env  map[string]string
}

func main() {
	log.Println("Starting authentication service")

	app, err := setup()
	if err != nil {
		log.Panicln(err)
	}

	lis, err := net.Listen("tcp", "0.0.0.0:50001")
	if err != nil {
		log.Panicln(err)
	}

	log.Println("Starting gRPC server")
	grpcServer := grpc.NewServer()
	authentication.RegisterAuthenticationServer(grpcServer, &AuthenticationServer{Config: app})
	err = grpcServer.Serve(lis)
	if err != nil {
		log.Panicln(err)
	}

}

func setup() (*Config, error) {
	app := Config{
		Env: map[string]string{},
	}

	for _, v := range []string{"DSN"} {
		env, isSet := os.LookupEnv("DSN")
		if !isSet {
			return nil, errors.New(fmt.Sprintf("Env %s not found", v))
		}
		app.Env[v] = env
	}

	conn := app.connectToDB()
	if conn == nil {
		return nil, errors.New("Cannot connect to database")
	}
	app.setupRepo(conn)

	return &app, nil
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

func (app *Config) connectToDB() *sql.DB {
	dsn := app.Env["DSN"]

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

func (app *Config) setupRepo(conn *sql.DB) {
	repo := data.NewPostgresRepository(conn)
	app.Repo = repo
}
