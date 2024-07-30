package main

import "log"

func main() {
	store, err := NewPostgresStore()
	if err != nil {
		log.Fatal(err)
	}
	err = store.init()
	if err != nil {
		log.Fatal(err)
	}

	newServer := NewAPIServer(":8082", store)
	newServer.Run()
}
