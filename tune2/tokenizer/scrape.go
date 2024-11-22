package main

import (
	"fmt"
	"log"
	"net/http"

	"golang.org/x/net/html"
)

func main() {
	url := "https://webscraper.io/test-sites/e-commerce"

	res, err := http.Get(url)
	if err != nil {
		log.Panicf("there's something wrong with get %v\n", err)
	}

	defer res.Body.Close()

	tokenizer := html.NewTokenizer(res.Body)

	for {
		tokenType := tokenizer.Next()

		switch tokenType {
		case html.ErrorToken:
			return
		case html.StartTagToken, html.SelfClosingTagToken:
			token := tokenizer.Token()
			fmt.Println(token.Data)
		}

	}

}
