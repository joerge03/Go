package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"logger-service/data"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	webPort  = "8089"
	rpcPort  = "5001"
	mongoURL = "mongodb://mongo:27017"
	gRpcPort = "50001"
)

// type client *mongo.Client

type Config struct {
	Models data.Models
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	client, err := connectToMongo(ctx, cancel)
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	fmt.Println("set app config")
	app := Config{
		Models: data.New(client),
	}

	app.serve()
}

func (app *Config) serve() {
	srv := http.Server{
		Addr:    fmt.Sprintf(`:%s`, webPort),
		Handler: app.routes(),
	}
	fmt.Println("running on port: ", webPort)

	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}

func connectToMongo(ctx context.Context, cancel context.CancelFunc) (*mongo.Client, error) {
	defer cancel()

	clientOptions := options.Client().ApplyURI(mongoURL)
	clientOptions.SetAuth(options.Credential{
		Username: "admin",
		Password: "password",
	})

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		fmt.Println("error connecting in client: ", client)
		return nil, err
	}

	return client, nil
}
