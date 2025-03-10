package dbminer

import (
	"fmt"
	"log"
	"regexp"
)

type DBMiner interface {
	GetSchema() (*Schema, error)
}

type Schema struct {
	Name      string
	Databases []Database
}

type Database struct {
	Name   string
	Tables []Table
}

type Table struct {
	Name    string
	Columns []string
}

func Search(dbMiner DBMiner) error {
	schema, err := dbMiner.GetSchema()
	if err != nil {
		log.Panic(err)
	}
	for _, s := range schema.Databases {
		fmt.Printf("[DATABASE] : %s \n", s.Name)
		for _, t := range s.Tables {
			fmt.Printf("[TABLE] : %s\n", t.Name)
			for _, c := range t.Columns {
				getRegex(c)
				// fmt.Println(c)
			}
		}
	}
	return nil
}

func getRegex(s string) {
	re := []*regexp.Regexp{
		regexp.MustCompile(`(?i)social`),
		regexp.MustCompile(`(?i)ssn`),
		regexp.MustCompile(`(?i)pass(word)?`),
		regexp.MustCompile(`(?i)hash`),
		regexp.MustCompile(`(?i)ccnum`),
		regexp.MustCompile(`(?i)card`),
		regexp.MustCompile(`(?i)security`),
		regexp.MustCompile(`(?i)key`),
	}
	for _, regex := range re {
		b := regex.Find([]byte(s))
		if b != nil {
			fmt.Printf("[+] HIT: %s\n", string(b))
		}
	}
}
