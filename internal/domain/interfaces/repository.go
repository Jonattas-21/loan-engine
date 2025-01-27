package interfaces

type Repository interface {
    SaveCollection[T any](itemToSave T) error
    GetItemCollection[T any](itemToGet T) ([]T, error)
}