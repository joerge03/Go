package main

import "log"

func main() {
	store, err := NewPostgresStore()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("%+v \n ", store)

	newServer := NewAPIServer(":8082", store)
	newServer.Run()
}
