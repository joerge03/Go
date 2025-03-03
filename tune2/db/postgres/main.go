package main

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Transactions struct {
	Ccnum  string
	Date   string
	Amount float32
	Cvv    string
	Exp    string
}

func main() {
	dsn := "host=localhost user=postgres password=admin dbname=store port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Panic(err, " cannot open gorm")
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Panic(err, "db error")
	}
	defer sqlDB.Close()
	if err := sqlDB.Ping(); err != nil {
		log.Panic(err, "ping error")
	}
	fmt.Println(" DB connected")
	db.AutoMigrate(&Transactions{})

	// db.Create(&Transactions{Ccnum: "Asdfasdf", Date: "10-2-13", Amount: 123, Cvv: "4521", Exp: "10-23-30"})
	var t Transactions

	g := db.First(&t, Transactions{Ccnum: "Asdfasdf"})
	if g.Error != nil {
		log.Panic(g.Error, "errr")
	}

	fmt.Printf("%+v asdf\n", t)

}
