package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

var (
	Template *template.Template
	err      error
)

func handleIndex(w http.ResponseWriter, r *http.Request) {
	Template.ExecuteTemplate(w, "index.html", nil)
}

func main() {
	// test := make(map[string][]string, 5)
	// test["Te"] = []string{"Test"}
	// test["Tesdf"] = []string{"Testasdf", "asdfasdfsaf"}
	// fmt.Println(test["Te"])
	// // fmt.Println(test)
	Template, err = template.ParseGlob("test/pages/*.html")
	if err != nil {
		fmt.Println(err)
	}
	// test := 2 / len([]string{"asdf", "23", "asdf", "23", "asdf", "23", "asdf", "23", "23", "23"})
	// log.Printf(`test sfadf %v`, test)
	http.HandleFunc("/", handleIndex)
	err := http.ListenAndServe(":8082", nil)
	log.Printf("running on port: %vasdfasdfasdf", 8081)
	if err != nil {
		fmt.Println(err, "errr")
		fmt.Println(err)
	}
}
