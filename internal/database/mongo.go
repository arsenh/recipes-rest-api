package database

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type MongoDatabase struct {
	Client *mongo.Client
	DB     *mongo.Database
}

func ConnectMongo(url, dbName string) (*MongoDatabase, error) {
	ctx := context.Background()

	// Create client
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(url))
	if err != nil {
		return nil, fmt.Errorf("failed to create MongoDB client: %w", err)
	}

	// Ping to verify connection
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		_ = client.Disconnect(ctx) // cleanup on failure
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	log.Println("connected to MongoDB successfully")

	return &MongoDatabase{
		Client: client,
		DB:     client.Database(dbName), // this call never return nil
	}, nil
}

// Close gracefully disconnects from MongoDB (call this in main with defer)
func (m *MongoDatabase) Close(ctx context.Context) error {
	if m.Client == nil {
		return nil
	}
	return m.Client.Disconnect(ctx)
}
