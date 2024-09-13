package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"authentication/data"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/joho/godotenv"
)

const port = "80"

var counts int64

type Config struct {
	DB     *sql.DB
	Models data.Models
}

func main() {
	godotenv.Load()

	fmt.Println(os.Getenv("DSN"))

	db := connectToDb()

	if db == nil {
		log.Fatal("Cannot connect right now hiii hiii.")
	}

	app := Config{
		DB:     db,
		Models: data.New(db),
	}

	server := http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: app.routes(),
	}

	err := server.ListenAndServe()
	if err != nil {
		log.Fatal("Something wrong with the server, err:", err)
	}
	log.Fatal("Running on localhosts:", port)
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}
	// conn, err := pgxpool.Connect(context.Background(), "test")
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}

func connectToDb() *sql.DB {
	dsnToken := os.Getenv("DSN")

	for {
		db, err := openDB(dsnToken)
		if err != nil {
			log.Fatal("there is a err", err)
			counts++
		} else {
			return db
		}

		if counts > 10 {
			log.Fatal("It took to many attempts of reconnecting")
			return nil
		}
		time.Sleep(2 * time.Second)
		continue
	}
}
