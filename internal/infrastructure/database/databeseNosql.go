package database

import (
	"context"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewDatabase() *mongo.Client {
	//Conection to the database
	host := os.Getenv("MONGO_HOST")
	clientOptions := options.Client().ApplyURI(host)
	client, err	 := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal("Error connecting to the database")
		panic(err)
	}

	return client
}