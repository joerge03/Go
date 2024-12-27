package main

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"tune/foobarTest/bingScrape/metadata"

	"github.com/PuerkitoBio/goquery"
)


func handler(i int, s *goquery.Selection){
	url, ok := s.Find("a").Attr("href")
	if !ok{
		return
	}

	fmt.Printf("%d: %s\n", i, url)

	res, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		return 
	}
	defer res.Body.Close()
	buf, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	r, err := zip.NewReader(bytes.NewReader(buf), int64(len(buf)))
	if err != nil {
		fmt.Println(err, "this")
		return 
	}
	core, app, err := metadata.NewProperties(r)
	if err != nil {
		fmt.Println(err)
		return
	}
	// fmt.Println("test")

	
	log.Printf("%25s %25s - %s %s\n", core.Creator, core.LastModifiedBy, app.Application,app.GetMajorVersion())
}

func main(){
	if len(os.Args) != 3{
		log.Fatal("Missing args, Usage: main.go domain ext")
	}

	

	domain := os.Args[1]
	fileType := os.Args[2]

	query := fmt.Sprintf("site:%s && filetype:%s", domain, fileType)

	search :=fmt.Sprintf("https://www.bing.com/search?q=%s",url.QueryEscape(query))
	

	
	res, err := http.Get(search)	
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	// d,_ := io.ReadAll(res.Body)
	doc, err := goquery.NewDocumentFromReader(res.Body)
	// doc, err := goquery.NewDocument(search)
	if err != nil {
		log.Fatal(err)
	}
	// fmt.Printf("%s\n", d)

	f := "html body div#b_content ol#b_results li.b_algo h2"	//just a css selector
	doc.Find(f).Each(handler)
}