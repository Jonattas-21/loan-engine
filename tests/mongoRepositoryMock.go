package tests

import (
	"github.com/stretchr/testify/mock"
)

type MockRepository[T any] struct {
	mock.Mock
}

func (m *MockRepository[T]) SaveItemCollection(itemToSave T) error {
	args := m.Called(itemToSave)
	return args.Error(0)
}

func (m *MockRepository[T]) GetItemsCollection(collection string) ([]T, error) {
	args := m.Called(collection)
	return args.Get(0).([]T), args.Error(1)
}

func (m *MockRepository[T]) UpdateItemCollection(name string, fields map[string]interface{}) error {
	args := m.Called(name, fields)
	return args.Error(0)
}

func (m *MockRepository[T]) DeleteItemCollection(collectionItemKey string) error {
	args := m.Called(collectionItemKey)
	return args.Error(0)
}

func (m *MockRepository[T]) Ping() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockRepository[T]) TrunkCollection() error {
	args := m.Called()
	return args.Error(0)
}
