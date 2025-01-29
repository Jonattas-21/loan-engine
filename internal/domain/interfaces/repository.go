package interfaces

type Repository[T any] interface {
	SaveItemCollection(itemToSave T) error
	GetItemsCollection(collectionKey string) ([]T, error)
	DeleteItemCollection(collectionItemKey string) error
	UpdateItemCollection(collectionItemKey string, fields map[string]interface{}) error
	Ping() error
}
