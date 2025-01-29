package repositories

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type DefaultRepository[T any] struct {
	Client         *mongo.Client
	DatabaseName   string
	CollectionName string
}

func (d *DefaultRepository[T]) SaveItemCollection(itemToSave T) error {
	collection := d.Client.Database(d.DatabaseName).Collection(d.CollectionName)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := collection.InsertOne(ctx, itemToSave)
	if err != nil {
		log.Println(fmt.Sprintf("Error during insert item in DB: %v", err.Error()))
		return err
	}

	return nil
}

func (d *DefaultRepository[T]) GetItemsCollection(collectionKey string) ([]T, error) {
	collection := d.Client.Database(d.DatabaseName).Collection(d.CollectionName)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var filter = bson.M{}
	//if there is no key, get all itens
	if collectionKey != "" {
		filter = bson.M{"id": collectionKey}
	}

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		log.Println(fmt.Sprintf("Error during get item in DB: %v", err.Error()))
		return nil, err
	}
	defer cursor.Close(ctx)

	var items []T
	for cursor.Next(ctx) {
		var item T
		err := cursor.Decode(&item)
		if err != nil {
			log.Println(fmt.Sprintf("Error during decode item in DB: %v", err.Error()))
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
		log.Println(fmt.Sprintf("Error during update item in DB: %v", err.Error()))
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
		log.Println(fmt.Sprintf("Error during delete item in DB: %v", err.Error()))
		return err
	}

	return nil
}

func (d *DefaultRepository[T]) Ping() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := d.Client.Ping(ctx, nil)
	if err != nil {
		log.Println(fmt.Sprintf("Error during ping in DB: %v", err.Error()))
		return err
	}

	return nil
}