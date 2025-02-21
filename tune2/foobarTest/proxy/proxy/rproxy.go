package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/gorilla/mux"
)

var host string

func init() {
	flag.StringVar(&host, "h", "", "-h for Hosts")

	flag.Parse()
	if len(host) == 0 {
		fmt.Println("Please provide a host with '-h' ")
	}
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	parsedHost, err := url.Parse(host)
	if err != nil {
		log.Panic(err, "unable to parse err.")
	}
	r.Host = parsedHost.Host
	r.URL.Host = parsedHost.Host
	r.URL.Scheme = parsedHost.Scheme
	r.RequestURI = ""

	res, err := http.DefaultClient.Do(r)
	if err != nil {
		log.Panic(err, "something wrong with defaultClient")
	}
	io.Copy(os.Stdout, res.Body)
	fmt.Printf("reverse proxy %v", time.Now())
}

func main() {

	r := mux.NewRouter()

	r.HandleFunc("/", handleHome)

	log.Fatal(http.ListenAndServe(":8081", r))

}
