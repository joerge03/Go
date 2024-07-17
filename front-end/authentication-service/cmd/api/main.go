package main

import "database/sql"

const port = "8080"

type Config struct {
	DB     *sql.DB
	Models data.Models
}

func main() {
}
