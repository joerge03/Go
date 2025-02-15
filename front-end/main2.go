package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
)

type Todo struct {
	Title       string
	Description string
	Done        bool
}

type TodoList []Todo

type Price float64

type FormData struct {
	email    string
	password string
}

var (
	todoList TodoList
	// err       error
	template1 *template.Template
	price     Price
)

func (p Price) PriceTest() string {
	remainder := int(p*100) % 5
	quotiant := int(p*100) / 5
	fmt.Println(quotiant, remainder)
	if remainder < 3 {
		s := float64(quotiant*5) / 100
		return fmt.Sprintf("%.2f", s)
	}
	s := float64((quotiant*5)+5) / 100
	return fmt.Sprintf("%.2f", s)
}

func init() {
	t := template.Must(template.New("name").Parse("page1"))
	// if err != nil {
	// 	fmt.Println(err)
	// }
	t1 := template.Must(t.New("page2").Parse("page 2"))
	// fmt.Printf("%+v\n", t.Templates())
	// for _, temp := range t.Templates() {
	// 	fmt.Println(temp)
	// }
	log.Fatal("Asdfsdaf")
	fmt.Printf("%+v\n", t1.Execute(os.Stdout, nil))
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	t, err := template.New("name").Parse("page1")
	if err != nil {
		fmt.Println(err)
	}
	_, err = t.New("page2").Parse("page 2")
	// fmt.Printf("%+v\n", t.Templates())
	// for _, temp := range t.Templates() {
	// 	fmt.Println(temp)
	// }
	log.Fatal(t.Execute(os.Stdout, price))
}

func handleSubmit(w http.ResponseWriter, r *http.Request) {
	var s FormData
	s.email = r.FormValue("email")
	s.password = r.FormValue("password")
	fmt.Println(s)
	fmt.Printf("email: %v ,password: %v ", s.email, s.password)
	template1.ExecuteTemplate(w, "success.html", nil)
}

func main() {
	price = 3.36
	// t, err := template.ParseGlob("test/pages/*.html")

	http.HandleFunc("/", handleHome)
	// log.Fatal()
	// http.HandleFunc("/submit", handleSubmit)
	http.ListenAndServe(":8080", nil)
}
