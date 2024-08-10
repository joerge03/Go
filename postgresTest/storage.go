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
	getAccountByNumber(int) (*Account, error)
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
	ID serial primary key,
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
	row, err := p.db.Query("SELECT * FROM account WHERE ID = $1", id)
	if err != nil {
		return nil, err
	}
	account, err := scanIntoAccount(row)
	if err != nil {
		return nil, err
	} else if len(account) == 0 {
		return nil, fmt.Errorf("account %v not found", id)
	}
	return account[0], nil
}

func (p *PostgresStore) getAccountByNumber(num int) (*Account, error) {
	row, err := p.db.Query("SELECT * FROM account WHERE number = $1", num)
	if err != nil {
		return nil, err
	}

	acc, err := scanIntoAccount(row)
	if err != nil {
		return nil, err
	} else if len(acc) == 0 {
		return nil, fmt.Errorf("account %v not found", num)
	}

	return acc[0], nil
}

func (p *PostgresStore) getAccounts() ([]*Account, error) {
	row, err := p.db.Query("select * from account")
	if err != nil {
		return nil, err
	}

	accounts, err := scanIntoAccount(row)
	if err != nil {
		return nil, err
	}

	return accounts, nil
}

func (p *PostgresStore) deleteAccount(id int) error {
	_, err := p.db.Query(`delete from account where ID = $1`, id)
	if err != nil {
		return err
	}
	return nil
}

func scanIntoAccount(rows *sql.Rows) ([]*Account, error) {
	account := []*Account{}

	for rows.Next() {
		acc := new(Account)
		if err := rows.Scan(&acc.ID, &acc.FirstName, &acc.LastName, &acc.Number, &acc.Balance, &acc.CreatedAt); err != nil {
			return nil, err
		}

		account = append(account, acc)
	}

	return account, nil
}
