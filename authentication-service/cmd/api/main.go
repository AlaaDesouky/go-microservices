package main

import (
	"authentication/data"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

const webPort = "80"
var count int64

type Config struct{
	DB *sql.DB
	Models data.Models
}

func main() {
	dbConn := connectToDB()
	if dbConn == nil {
		log.Panic("Con't connect to Postgres!")
	}
	
	app := Config{
		DB: dbConn,
		Models: data.New(dbConn),
	}
	
	srv := &http.Server{
		Addr: fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}
	
	log.Printf("Starting authentication service on port %s\n", webPort)
	if err := srv.ListenAndServe(); err != nil {
		log.Panic(err)
	}
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func connectToDB() *sql.DB {
	dsn := os.Getenv("DSN")

	for {
		connection, err := openDB(dsn)
		if err != nil {
			log.Println("Postgres not yet ready...")
			count++
		} else {
			log.Println("Connected to Postgres!")
			return connection
		}

		if count > 10 {
			log.Println(err)
			return nil
		}

		log.Panicln("Backing off for two seconds...")
		time.Sleep(2 * time.Second)
		continue
	}
}