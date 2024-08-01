package main

import (
	"math/rand"
	"time"
)

type CreateAccountRequest struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

type TransferRequest struct {
	ToAccount int `json:"toAccount"`
	Amount    int `json:"amount"`
}

type Account struct {
	ID        string    `json:"ID"`
	LastName  string    `json:"lastName"`
	FirstName string    `json:"firstName"`
	Number    int64     `json:"number"`
	Balance   int64     `json:"balance"`
	CreatedAt time.Time `json:"createdAt"`
}

func NewAccount(firstName, lastName string) *Account {
	return &Account{
		LastName:  lastName,
		FirstName: firstName,
		Number:    rand.Int63n(10000),
		Balance:   int64(20),
		CreatedAt: time.Now().Local(),
	}
}
