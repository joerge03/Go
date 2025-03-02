package main

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Payment struct {
	CCnum  string  `bson:"ccnum"`
	Date   string  `bson:"date"`
	Amount float64 `bson:"amount"`
	Cvv    string  `bson:"cvv"`
	Exp    string  `bson:"exp"`
}

type Transaction struct {
	ID  primitive.ObjectID `bson:"_id"`
	Key Payment
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Panic(err, "mongo db err")
	}

	collection := client.Database("store").Collection("transactions")

	res, err := collection.Find(context.Background(), bson.M{})
	// _, err = collection.InsertOne(context.Background(), Transaction{ID: primitive.NewObjectID(), Key: Payment{CCnum: "asdfsad", Date: "12-25-23", Amount: 123.45, Cvv: "123213", Exp: "123123"}})
	if err != nil {
		log.Panic(err, "test")
	}
	defer res.Close(ctx)
	// var a bson.A
	var t []Transaction

	// res.All(ctx, &t)
	if err := res.All(ctx, &t); err != nil {
		log.Panic("there's something wrong decoding", err)
	}
	for _, transaction := range t {
		fmt.Printf("%v\n", transaction)
		// if tr, ok := transaction.(Transaction); ok {

		// 	t = append(t, tr)
		// 	fmt.Printf("the value %+v\n", t)
		// }
	}
}
