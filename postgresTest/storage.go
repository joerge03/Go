package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

type Storage interface {
	createAccount(*Account) error
	updateAccount(*Account) error
	getAccountByID(int) (*Account, error)
	deleteAccount(int) error
}

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore() (*PostgresStore, error) {
	connStr := "user=pqgotest dbname=pqgotest sslmode=verify-full"

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &PostgresStore{
		db,
	}, nil

	// rows, err := db.Query()
}

func (s *PostgresStore) init() error {
	return nil
}

func (s *PostgresStore) createAccountTable() error {
	query := `create table account if not exist (
	id serial primary key,
	firstName varcha(50),
	lastName varcha(50),
	number serial,
	balance,
	createdAt timestamp
	)`
	return nil
}

func (p *PostgresStore) createAccount(account *Account) error {
	// p.
	return nil
}

func (p *PostgresStore) updateAccount(account *Account) error {
	return nil
}

func (p *PostgresStore) getAccountByID(id int) (*Account, error) {
	return nil, nil
}

func (p *PostgresStore) deleteAccount(id int) error {
	return nil
}
