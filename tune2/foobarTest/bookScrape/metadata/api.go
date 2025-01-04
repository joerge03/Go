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


func removePriceCharacters(s string) float64{
	regex := regexp.MustCompile(`[^0-9.]`)
	
	str := ""
	var result float64
	var err error
	if len(s) > 0 {
		str = regex.ReplaceAllString(s, "")
		result, err = strconv.ParseFloat(str, 64)
	}
	if err != nil {
		fmt.Println(err, "regex err")
		return 0
	}
	return result
}

func (b *BookDetailInfoMap) Populate(book *Book ) {
	priceEx := removePriceCharacters((*b)[2])
	priceInc := removePriceCharacters((*b)[3])
	tax := removePriceCharacters((*b)[4])
	numberOfReviews := removePriceCharacters((*b)[6])
	
	
	book.BookDetail.UPC = (*b)[0]
	book.BookDetail.Product_type = (*b)[1]
	book.BookDetail.Price_excl_tax = priceEx
	book.BookDetail.Price_incl_tax = priceInc
	book.BookDetail.Tax = tax
	book.BookDetail.Availability = (*b)[5]
	book.BookDetail.Number_of_reviews = numberOfReviews
}

func processBooks(s *goquery.Document, book *Book){
	descriptionSelector := "html body div div.page_inner div.content div#content_inner article.product_page"
	book.Description = s.Find("#product_description + p").First().Text()

	BookDetailMap := make(BookDetailInfoMap)
	s.Find(fmt.Sprintf("%s table tbody tr",descriptionSelector)).Each(func(i int, qs *goquery.Selection){
		value := qs.Find("td").Text()
		BookDetailMap[i] = value
	})
	
	BookDetailMap.Populate(book)
}




func process(i int, s *goquery.Selection, books *Books) {
	book := new(Book)
	
	link, ok := s.Find("div.image_container a").Attr("href")
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
	
	isStock, ok := s.Find("produc1t_price p.instock_availability i.icon-ok").Attr("class")
	book.InStock = false	
	if isStock != "icon-ok" || !ok {
	}else {
		book.InStock = true
	}
	processBooks(selection,book)
	(*books) = append((*books), *book)
}
	
func NewBooks(r io.Reader, wg *sync.WaitGroup) (*Books, error){
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



