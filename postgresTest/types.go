package main

import (
	"math/rand"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type CreateAccountRequest struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Password  string `json:"password"`
}

type TransferRequest struct {
	ToAccount int `json:"toAccount"`
	Amount    int `json:"amount"`
}

type Account struct {
	ID                string    `json:"ID"`
	LastName          string    `json:"lastName"`
	FirstName         string    `json:"firstName"`
	EncryptedPassword string    `json:"-"`
	Number            int64     `json:"number"`
	Balance           int64     `json:"balance"`
	CreatedAt         time.Time `json:"createdAt"`
}

func NewAccount(firstName, lastName string, password string) (*Account, error) {
	encryptedPw, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	return &Account{
		LastName:          lastName,
		FirstName:         firstName,
		Number:            rand.Int63n(10000),
		Balance:           int64(20),
		EncryptedPassword: string(encryptedPw),
		CreatedAt:         time.Now().Local(),
	}, nil
}
