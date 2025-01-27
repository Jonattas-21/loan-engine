package repositories

import (
	"go.mongodb.org/mongo-driver/mongo"
)

type LoanRepository struct {
	Database *mongo.Client
}

func SaveCollection[T any](itemToSave T) error {
	//todo

	return nil
}

func GetItemCollection[T any](itemToGet T) []T {
	//todo

    return []T{itemToGet}
}