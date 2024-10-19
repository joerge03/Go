package main

import (
	"fmt"
	"log"
	"os"
)

func main() {

	// test := bytes.NewBufferString("")
	file, _ := os.Create("test.txt")
	log.SetOutput(file)
	defer file.Close()

	test := log.New(file, "1", log.Ldate|log.Lshortfile)
	log.Printf(`testasdfasdfsadf`)
	log.Printf(`testasdfasdfsadf3`)
	log.Printf(`testasdfasdfsadf2`)
	test.Println("test12121")

	fmt.Println(file)
}
