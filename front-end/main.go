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

func main2() {
	Template, err = template.ParseGlob("test/pages/*.html")
	if err != nil {
		fmt.Println(err)
	}
	http.HandleFunc("/", handleIndex)
	err := http.ListenAndServe(":8082", nil)
	log.Printf("running on port: %vasdfasdfasdf", 8081)
	if err != nil {
		fmt.Println(err)
	}
}
