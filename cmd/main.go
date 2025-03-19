package main

import (
	"database/sql"
	"fmt"
	"frappuccino/internal/server"
	"frappuccino/internal/utils"
	"log"

	_ "github.com/lib/pq"
)

func main() {
	config, err := getConfig()
	if err != nil {
		log.Fatal(err)
	}

	connString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.Postgres.Host,
		config.Postgres.Port,
		config.Postgres.User,
		config.Postgres.Password,
		config.Postgres.Database,
	)

	db, err := connectDB(connString)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	server := server.NewServer(":8080", db, utils.GetLogger())
	server.RunServer()
}

func connectDB(connectionString string) (*sql.DB, error) {
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, err
	}

	log.Println("Connected to Database successfully!")

	return db, nil
}
