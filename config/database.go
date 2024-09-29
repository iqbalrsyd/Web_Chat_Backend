package config

import (
    "context"
    "log"
    "time"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

type Database struct {
    Client   *mongo.Client
    Database *mongo.Database
}

// NewDatabase initializes the MongoDB connection and returns the Database instance
func NewDatabase(mongoURI, dbName string) (*Database, error) {
    clientOptions := options.Client().ApplyURI(mongoURI)
    client, err := mongo.NewClient(clientOptions)
    if err != nil {
        return nil, err
    }

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    err = client.Connect(ctx)
    if err != nil {
        return nil, err
    }

    db := client.Database(dbName)
    log.Println("Connected to MongoDB")

    return &Database{
        Client:   client,
        Database: db,
    }, nil
}

// Disconnect closes the MongoDB connection
func (d *Database) Disconnect() error {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    return d.Client.Disconnect(ctx)
}
