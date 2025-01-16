package memory

type Memory[T any] interface {
	Add(input T) error
	Get() ([]T, error)
}
