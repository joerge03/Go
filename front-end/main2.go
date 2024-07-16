package main

import (
	"fmt"
	"html/template"
	"net/http"
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

func handleHome(w http.ResponseWriter, r *http.Request) {
	template1.ExecuteTemplate(w, "index.html", price)
}

func handleSubmit(w http.ResponseWriter, r *http.Request) {
	fmt.Println("running like a billin")
	var s FormData

	s.email = r.FormValue("email")
	s.password = r.FormValue("password")

	fmt.Println(s)

	fmt.Printf("email: %v ,password: %v ", s.email, s.password)

	template1.ExecuteTemplate(w, "success.html", nil)
}

func mains() {
	// r := gin.Default()
	price = 3.36
	todoList = TodoList{
		{Title: "train to busan", Description: "busan to train", Done: false},
		{Title: "do chore", Description: "fckit", Done: false},
		{Title: "play games", Description: "yep", Done: true},
	}

	template1, err = template.ParseGlob("test/pages/*.html")
	if err != nil {
		fmt.Println(err)
	}
	http.HandleFunc("/", handleHome)
	http.HandleFunc("/submit", handleSubmit)
	http.ListenAndServe(":8080", nil)
}
