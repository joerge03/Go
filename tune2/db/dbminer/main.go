package main

import (
	"context"
	"fmt"
	"log"
	"regexp"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DatabaseMiner interface {
	GetSchema() (*Schema, error)
}

type Schema struct {
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

func Search(m DatabaseMiner) error {
	s, err := m.GetSchema()
	if err != nil {
		return err
	}

	re := getRegex()

	for _, database := range s.Databases {
		// fmt.Println(database)
		for _, table := range database.Tables {
			// fmt.Println(table)
			for _, col := range table.Columns {
				fmt.Println(col)
				for _, test := range re {
					if test.MatchString(col) {

						fmt.Printf("%+v\n", database)
						fmt.Printf("[+] HIT: %v\n", test)
					}
				}
			}
		}
	}
	return nil
}

func getRegex() []*regexp.Regexp {
	return []*regexp.Regexp{
		regexp.MustCompile(`(?i)social`),
		regexp.MustCompile(`(?i)ssn`),
		regexp.MustCompile(`(?i)pass(word)?`),
		regexp.MustCompile(`(?i)hash`),
		regexp.MustCompile(`(?i)ccnum`),
		regexp.MustCompile(`(?i)card`),
		regexp.MustCompile(`(?i)security`),
		regexp.MustCompile(`(?i)key`),
	}
}

type MongoMiner struct {
	Host   string
	Client *mongo.Client
	ctx    context.Context
}

func New(ctx *context.Context, h string) (*MongoMiner, error) {
	m := &MongoMiner{Host: h, ctx: *ctx}
	client, err := mongo.Connect(*ctx, options.Client().ApplyURI(fmt.Sprintf("mongodb://%s", h)))
	if err != nil {
		log.Panic(err, "cannot connect")
	}
	m.Client = client
	return m, nil
}

func (m *MongoMiner) GetSchema() (*Schema, error) {
	var schema Schema
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	names, err := m.Client.ListDatabases(ctx, bson.M{})
	if err != nil {
		log.Panic(err, "unable to list database names")
	}
	for _, database := range names.Databases {
		db := Database{Name: database.Name, Tables: []Table{}}
		dbNames, err := m.Client.Database(database.Name).ListCollectionNames(m.ctx, bson.M{})
		if err != nil {
			log.Panic(err, "cannot get collection names")
		}
		for _, collections := range dbNames {
			table := Table{Name: collections}
			collection := m.Client.Database(database.Name).Collection(collections)
			// fmt.Printf("collections %v\n")
			cursor := collection.FindOne(ctx, bson.D{})
			// data := bson.D{}
			// cursor.All(ctx, &data)
			// fmt.Printf("data %+v\n ", data)
			// if err != nil {
			// 	log.Panic(err)
			// }

			var v bson.Raw
			err := cursor.Decode(&v)
			if err != nil {
				log.Panic(err)
			}
			// fmt.Println(v)
			// fmt.Printf("%v\n", v)

			// var doc bson.Raw
			// for cursor.Next(ctx) {
			// 	cursor.Decode(&doc)
			// 	fmt.Printf("-- %+v --\n", doc)

			// 	// table.Columns = append(table.Columns, doc)
			// }
			if err != nil {
				fmt.Println(err)
			}
			// for _, d := range doc {
			// 	fmt.Printf("%v -- %v -- %v\n", d.Key(), d.String(), d.Value())
			// 	fmt.Println(d.String(), "doc string")
			// 	table.Columns = append(table.Columns, d.String())
			// }
			docs, err := v.Elements()
			if err != nil {
				log.Panic(err)
			}
			for _, rawDoc := range docs {
				fmt.Println(rawDoc.Key())
				table.Columns = append(table.Columns, rawDoc.String())
			}

			db.Tables = append(db.Tables, table)
		}

		schema.Databases = append(schema.Databases, db)
	}

	return &schema, nil
}

func main() {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	m, err := New(&ctx, "localhost:27017")
	if err != nil {
		log.Panic(err)
	}
	err = Search(m)
	// _, err = m.GetSchema()
	if err != nil {
		log.Panic(err)
	}

}
