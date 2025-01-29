package repositories

import (
	"context"
	"time"

	"fmt"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type DefaultRepository[T any] struct {
	Client         *mongo.Client
	DatabaseName   string
	CollectionName string
	Logger         *logrus.Logger
}

func (d *DefaultRepository[T]) SaveItemCollection(itemToSave T) error {
	collection := d.Client.Database(d.DatabaseName).Collection(d.CollectionName)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	//todo insert ttl
	_, err := collection.InsertOne(ctx, itemToSave)
	if err != nil {
		d.Logger.Error(fmt.Printf("Error during insert item in DB: %v", err.Error()))
		return err
	}

	return nil
}

func (d *DefaultRepository[T]) GetItemsCollection(itemId string) ([]T, error) {
	collection := d.Client.Database(d.DatabaseName).Collection(d.CollectionName)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var filter = bson.D{}
	//if there is no key, get all itens
	// if itemId != "" {
	// 	filter = bson.D{{"_id", "itemId"}}
	// }

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		d.Logger.Error(fmt.Printf("Error during get item in DB: %v", err.Error()))
		return nil, err
	}
	defer cursor.Close(ctx)

	var items []T
	for cursor.Next(ctx) {
		var item T
		err := cursor.Decode(&item)
		if err != nil {
			d.Logger.Error(fmt.Printf("Error during decode item in DB: %v", err.Error()))
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}

func (d *DefaultRepository[T]) UpdateItemCollection(collectionItemKey string, fields map[string]interface{}) error {
	collection := d.Client.Database(d.DatabaseName).Collection(d.CollectionName)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"name": collectionItemKey}
	update := bson.M{}

	for key, value := range fields {
		update["$set"].(bson.M)[key] = value
	}

	_, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		d.Logger.Error(fmt.Printf("Error during update item in DB: %v", err.Error()))
		return err
	}
	return nil
}

func (d *DefaultRepository[T]) DeleteItemCollection(collectionItemKey string) error {
	collection := d.Client.Database(d.DatabaseName).Collection(d.CollectionName)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := collection.DeleteOne(ctx, collectionItemKey)
	if err != nil {
		d.Logger.Error(fmt.Printf("Error during delete item in DB: %v", err.Error()))
		return err
	}

	return nil
}

func (d *DefaultRepository[T]) Ping() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := d.Client.Ping(ctx, nil)
	if err != nil {
		d.Logger.Error(fmt.Printf("Error during ping in DB: %v", err.Error()))
		return err
	}

	return nil
}

func (d *DefaultRepository[T]) TrunkCollection() error {
	collection := d.Client.Database(d.DatabaseName).Collection(d.CollectionName)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := collection.DeleteMany(ctx, bson.M{})
	if err != nil {
		d.Logger.Error(fmt.Printf("Error during trunk collection in DB: %v", err.Error()))
		return err
	}

	return nil
}
