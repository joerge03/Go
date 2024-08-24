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
)

const port = "8085"

var counts int64

type Config struct {
	DB     *sql.DB
	Models data.Models
}

func main() {
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

	fmt.Printf("Running on localhosts:%v\n", port)
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

func connectToDb() *sql.DB {
	dsnToken := os.Getenv("DSN")

	for {
		db, err := openDB(dsnToken)
		if err != nil {
			fmt.Printf("there is a err %v", err)
			counts++
		} else {
			return db
		}

		if counts > 10 {
			fmt.Printf("It took to many attempts of reconnecting")
			return nil
		}

		fmt.Println("sleeping hehe XD HIIII HIIII. OWWW")
		time.Sleep(2 * time.Second)
		continue
	}
}
