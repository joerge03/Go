package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"tune/foobarTest/bookScrape/metadata"
)

var wg sync.WaitGroup


func proccess(url string, wg *sync.WaitGroup, books *metadata.Books) {
	res, err := http.Get(url)
	if err != nil {
		log.Panic(err, "unabled to get the give link")
	}
	defer res.Body.Close()

	data, err := metadata.NewBooks(res.Body, wg)
	if err != nil {
		log.Fatal(err, "unable to create newbook")
	}
	*books = append(*books, *data...)
}

func main(){
	books := new(metadata.Books)
	for i := 0; i <= 50 ;i++{
		wg.Add(1)
		url := fmt.Sprintf("https://books.toscrape.com/catalogue/page-%d.html", i)
		go proccess(url, &wg, books)
	}
	wg.Wait()
	fmt.Printf("\n%+v --------\n", books)
}	