package main

import "math/rand"

type Account struct {
	ID        int    `json:"id"`
	LastName  string `json:"lastName"`
	FirstName string `json:"firstName"`
	Number    int64  `json:"number"`
	Balance   int64  `json:"balance"`
}

func NewAccount(firstName, lastName string) *Account {
	return &Account{
		ID:        rand.Intn(10000000),
		LastName:  lastName,
		FirstName: firstName,
		Number:    int64(rand.Uint64()),
		Balance:   int64(20),
	}
}
