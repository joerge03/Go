package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
)

type Config struct {
	Mailer Mail
}

const port = "80"

func main() {
	app := Config{
		Mailer: createMail(),
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf(`:%s`, port),
		Handler: app.routes(),
	}

	err := srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(`listening on port: `, port)
}

func createMail() Mail {
	port, _ := strconv.Atoi(os.Getenv("MAIL_PORT"))
	return Mail{
		Domain:      os.Getenv("MAIL_DOMAIN"),
		Host:        os.Getenv("MAIL_HOST"),
		Port:        port,
		Username:    os.Getenv("MAIL_USERNAME"),
		Password:    os.Getenv("MAIL_PASSWORD"),
		Encryption:  os.Getenv("MAIL_ENCRYPTION"),
		FromAddress: os.Getenv("MAIL_FROM_ADDRESS"),
		FromName:    os.Getenv("MAIL_FROM_NAME"),
	}
}
