package data

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

type LogEntry struct {
	ID        string    `bson:"_id,omitempty" json:"id,omitempty"`
	Name      string    `bson:"name" json:"name"`
	Data      any       `bson:"data" json:"data"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}

type Models struct {
	LogEntry LogEntry
}

func New(m *mongo.Client) Models {
	client = m

	return Models{LogEntry: LogEntry{}}
}

func (l *LogEntry) InsertOne(entry LogEntry) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()
	collection := client.Database("logs").Collection("logs")

	_, err := collection.InsertOne(ctx, LogEntry{Name: entry.Name, Data: entry.Data, CreatedAt: time.Now(), UpdatedAt: time.Now()})
	if err != nil {
		log.Fatalf("InsertOne failed: %v", err)
	}
	return nil
}

func (l *LogEntry) All() []*LogEntry {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	collection := client.Database("logs").Collection("logs")

	opts := options.Find()
	opts.SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := collection.Find(ctx, bson.D{}, opts)
	if err != nil {
		log.Fatal(err)
	}

	defer cursor.Close(ctx)
	var results []*LogEntry

	for cursor.Next(ctx) {
		result := new(LogEntry)
		if err := cursor.Decode(&result); err != nil {
			log.Fatal(`err decoding cursor : `, err)
		}
		results = append(results, result)
	}
	// err = cursor. All(ctx, &results)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	return results
}

func (l *LogEntry) GetOne(id string) (*LogEntry, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()
	collection := client.Database("logs").Collection("logs")

	docId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	entry := new(LogEntry)

	err = collection.FindOne(ctx, bson.M{"_id": docId}).Decode(entry)
	if err != nil {
		return nil, err
	}

	return entry, nil
}

func (l *LogEntry) DropCollection() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	collections := client.Database("logs").Collection("logs")

	err := collections.Drop(ctx)
	if err != nil {
		log.Fatal(`there's an error dropping collections: `, err)
	}
	return nil
}

func (l *LogEntry) UpdateOne() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	collections := client.Database("logs").Collection("logs")
	id, err := primitive.ObjectIDFromHex(l.ID)
	if err != nil {
		return err
	}
	filter := bson.D{{Key: "_id", Value: id}}
	update := bson.D{{Key: "$set", Value: bson.D{
		{Key: "name", Value: l.Name},
		{Key: "data", Value: l.Data},
		{Key: "updated_at", Value: time.Now()},
	}}}

	_, err = collections.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	return nil
}
