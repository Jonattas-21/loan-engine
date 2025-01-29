package database

import (
	"context"
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DatabaseNosql struct {
	Logger *logrus.Logger
}

func (d *DatabaseNosql) NewDatabase() *mongo.Client {
	//Conection to the database
	host := os.Getenv("MONGO_HOST")
	clientOptions := options.Client().ApplyURI(host)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		d.Logger.Errorln("Error connecting to the database")
		panic(err)
	}

	return client
}

func (d *DatabaseNosql) CloseDatabase(client *mongo.Client) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := client.Disconnect(ctx)
	if err != nil {
		d.Logger.Errorln("Error disconnecting from the database")
		panic(err)
	}
}
