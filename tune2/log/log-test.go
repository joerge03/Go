package main

import (
	"fmt"
	"io"
	"log"
	"os"
)

func main() {

	// test := bytes.NewBufferString("")
	file, _ := os.Create("test.txt")
	testW := io.MultiWriter(file, os.Stdout)
	log.SetOutput(testW)

	defer file.Close()

	test := log.New(testW, "1", log.Ldate|log.Lshortfile)
	log.Printf(`testasdfasdfsadf`)
	log.Printf(`testasdfasdfsadf3`)
	log.Printf(`testasdfasdfsadf2`)
	test.Println("test12121")

	fmt.Println(file)
}
