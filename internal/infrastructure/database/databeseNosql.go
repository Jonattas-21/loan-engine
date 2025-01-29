package database

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewDatabase() *mongo.Client {
	//Conection to the database
	host := os.Getenv("MONGO_HOST")
	clientOptions := options.Client().ApplyURI(host)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal("Error connecting to the database")
		panic(err)
	}

	return client
}

func CloseDatabase(client *mongo.Client) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := client.Disconnect(ctx)
	if err != nil {
		log.Fatal("Error disconnecting from the database")
		panic(err)
	}
}
