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

	data, err := metadata.NewBook(res.Body, wg)
	if err != nil {
		log.Fatal(err, "unable to create newbook")
	}
	*books = append(*books, *data)
	// fmt.Printf("%+v\n", data)
	// metadata.NewBook(data)
}

func main(){
	// if len(os.Args) < 1 {
	// 	fmt.Println("missing args. format eg. (main.go 'link' 'how many pages')")
	// }

	// link := os.Args[1]
	// pages := os.Args[2]
	// apiChannel := make()
	books := new(metadata.Books)	
	for i := 0; i <= 1 ;i++{
		wg.Add(1)
		url := fmt.Sprintf("https://books.toscrape.com/catalogue/page-%d.html", i)
		go proccess(url, &wg, books)
	}
	wg.Wait()
	fmt.Printf("%+v ------- \n\n", books)
}	