package metadata

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"sync"

	"github.com/PuerkitoBio/goquery"
)

type BookDetail struct {
	UPC string	
	Product_type string
	Price_excl_tax float64
	Price_incl_tax float64
	Tax float64
	Availability string
	Image string
	Number_of_reviews float64
	// Stocks string
}

type BookDetailInfoMap map[int]string 

type Book struct {
	Name string
	Price int64
	Rating int64
	Description string
	InStock bool
	BookDetail BookDetail
}

type Books []Book

// var BookDetailMap = BookDetailInfoMap{
// 	0:"UPC",
// 	1:"Product_type",
// 	2:"Price_excl_tax",
// 	3:"Price_incl_tax",
// 	4:"Tax",
// 	5:"Availability",
// 	6:"Number_of_reviews",
// }

func removePriceCharacters(s string) float64{
	regex := regexp.MustCompile(`[^0-9.]`)

	fmt.Println("test price", s)
	
	str := ""
	var result float64
	var err error
	if len(s) > 0 {
		str = regex.ReplaceAllString(s, "")
		fmt.Println("test result", str)
		result, err = strconv.ParseFloat(str, 64)
	}
	// fmt.Println("--------", s, "-------------")
	if err != nil {
		fmt.Println(err, "regex err")
		return 0
	}
	return result
}

func processBooks(s *goquery.Document, book *Book){
	// fmt.Printf("%v test1 \n", i)
	descriptionSelector := "html body div div.page_inner div.content div#content_inner article.product_page"
	// fmt.Println(s.Find("#product_description + p").First().Text(), "first p")
	book.Description = s.Find("#product_description + p").First().Text()
	// book.Description = s.Find(fmt.Sprintf("%s p",descriptionSelector)).First().Text()

	BookDetailMap := make(BookDetailInfoMap)

	// fmt.Println(fmt.Sprintf("%s table tbody",descriptionSelector),"test1")
	s.Find(fmt.Sprintf("%s table tbody tr",descriptionSelector)).Each(func(i int, qs *goquery.Selection){	
		value := qs.Find("td").Text()
		fmt.Println(value,"test1", i)
		BookDetailMap[i] = value
	})
		
	priceEx := removePriceCharacters(BookDetailMap[2])
	priceInc := removePriceCharacters(BookDetailMap[3])
	tax := removePriceCharacters(BookDetailMap[4])
	// numberOfReviews := removePriceCharacters()
	numberOfReviews := removePriceCharacters(BookDetailMap[6])
	
	
	book.BookDetail.UPC = BookDetailMap[0]
	book.BookDetail.Product_type = BookDetailMap[1]
	book.BookDetail.Price_excl_tax = priceEx
	book.BookDetail.Price_incl_tax = priceInc
	book.BookDetail.Tax = tax
	book.BookDetail.Availability = BookDetailMap[5]
	book.BookDetail.Number_of_reviews = numberOfReviews
}

func processDetail(s *goquery.Document, book *Book){
	// selectionString := "html body div div.page_inner div.content div#content_inner article.product_page div.row"	
	// s.Find(selectionString).Each(func(i int, qs *goquery.Selection){
	// 	processBooks(i,qs, book)
	// })
	processBooks(s,book)
}



func process(i int, s *goquery.Selection, books *Books) {
	book := new(Book)
	
	link, ok := s.Find("div.image_container a").Attr("href")
	fmt.Println(link, "link")
	mainLink := fmt.Sprintf("https://books.toscrape.com/catalogue/%s", link)
	if !ok {
		fmt.Println("unknown element")
		return
	}
	res, err := http.Get(mainLink)
	if err != nil {
		log.Panic(err, "could not get link")
	}
	defer res.Body.Close()
	
	selection, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil { 
		log.Panic(err, "something wrong creating new document from reader")
	}
	
	
	title, ok := s.Find("h3 a").Attr("title")	
	if !ok{
		fmt.Println("unable to get title")
		return 
	}
	book.Name = title
	
	isStock, ok := s.Find("product_price p.instock_availability i.icon-ok").Attr("class")
	book.InStock = false				
	if isStock != "icon-ok" || !ok {
	}else {
		book.InStock = true
	}	
	processDetail(selection,book)
	(*books) = append((*books), *book)
}
	
	
func NewBook(r io.Reader, wg *sync.WaitGroup) (*Books, error){
	defer wg.Done()
	books := new(Books)
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return nil,err
	}
	selector := "html body div.page_inner div.row section div ol.row li article.product_pod"

	doc.Find(selector).Each(func(i int, s *goquery.Selection){
		process(i,s,books)
	})
	return books, nil
}



