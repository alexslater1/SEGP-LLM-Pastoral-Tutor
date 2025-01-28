package ragger

type Ragger interface {
	ChunksFor(text string) ([]string, error)
	ContactsFrom(text string) ([]string, error)

	EmbeddingsFor(text string) ([]float32, error)
}
