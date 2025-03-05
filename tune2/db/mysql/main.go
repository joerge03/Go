package main

import (
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Transactions struct {
	CCNum  string
	Date   string
	Amount float32
	Cvv    string
	Exp    string
}

func main() {
	dsn := "root:admin@tcp(127.0.0.1:3306)"
	db, err := gorm.Open(mysql.Open(dsn))
	if err != nil {
		log.Panic(err, "Failed to open db")
	}

	db.AutoMigrate(&Transactions{})

	// db.Create(Transactions{CCNum: "12345", Amount: 5123.3, Date: "7-4-25", Cvv: "3114", Exp: "4-2-26"})

	var transactions []Transactions
	db.Select("CCNum", "Amount", "Date", "Cvv").Find(&transactions)

	fmt.Println("done!")

	for _, t := range transactions {
		fmt.Printf("%+v\n", t)
	}

}
