package main

import (
	"flag"
	"fmt"
	"log"
)

func seedAccount(s Storage, firstName string, lastName string, password string) *Account {
	acc, err := NewAccount(firstName, lastName, password)
	if err != nil {
		log.Fatal(err)
	}

	err = s.createAccount(acc)
	if err != nil {
		log.Println("errrr", err)
	}
	return acc
}

func seedAccounts(s Storage) {
	seedAccount(s, "seedos", "amigos", "passwords")
}

func main() {
	seed := flag.Bool("seed", false, "seed the db")
	flag.Parse()

	store, err := NewPostgresStore()
	if err != nil {
		log.Fatal(err)
	}
	err = store.init()
	if err != nil {
		log.Fatal(err)
	}

	if *seed {
		fmt.Println("Teststs")
		seedAccounts(store)
	}

	newServer := NewAPIServer(":8082", store)
	newServer.Run()
}
