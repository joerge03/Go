package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

type Storage interface {
	createAccount(*Account) error
	updateAccount(*Account) error
	getAccountByID(int) (*Account, error)
	deleteAccount(int) error
	getAccounts() ([]*Account, error)
}

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore() (*PostgresStore, error) {
	connStr := "user=postgres dbname=postgres password=golang sslmode=disable"

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		log.Println(err, "printed")
		return nil, err
	}

	return &PostgresStore{
		db,
	}, nil

	// rows, err := db.Query()
}

func (s *PostgresStore) init() error {
	return s.createAccountTable()
}

func (s *PostgresStore) createAccountTable() error {
	query := `create table if not exists account (
	id serial primary key,
	firstName varchar(50),
	lastName varchar(50),
	number serial,
	balance serial,
	createdAt timestamp
	)`

	_, err := s.db.Exec(query)
	if err != nil {
		log.Fatal(err)
	}
	return nil
}

func (p *PostgresStore) createAccount(account *Account) error {
	query := `
	insert into account
	(firstName, lastName, number, balance, createdAt)
	values
	($1, $2, $3, $4, $5)
	`

	log.Println(account.FirstName, account.LastName, account.Number, account.Balance, account.CreatedAt)

	response, err := p.db.Query(query, account.FirstName, account.LastName, account.Number, account.Balance, account.CreatedAt)
	if err != nil {
		return err
	}

	log.Printf("%+v \n", response)

	return nil
}

func (p *PostgresStore) updateAccount(account *Account) error {
	return nil
}

func (p *PostgresStore) getAccountByID(id int) (*Account, error) {
	return nil, nil
}

func (p *PostgresStore) getAccounts() ([]*Account, error) {
	accounts := []*Account{}

	row, err := p.db.Query("select * from account")
	if err != nil {
		return nil, err
	}

	account := new(Account)
	for row.Next() {
		err := row.Scan(&account.ID, &account.FirstName, &account.LastName, &account.Number, &account.Balance, &account.CreatedAt)
		if err != nil {
			return nil, err
		}
	}

	fmt.Printf("%+v test \n ", account)

	return accounts, nil
}

func (p *PostgresStore) deleteAccount(id int) error {
	return nil
}
