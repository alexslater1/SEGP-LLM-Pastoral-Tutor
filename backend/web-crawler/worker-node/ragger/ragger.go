package ragger

type Ragger interface {
	ChunksFrom(text string) ([]string, error)
	ContactsFrom(text string) ([]Contact, error)

	EmbeddingsFor(text string) ([]float32, error)
	EmbeddingsForAll(texts []string) ([][]float32, error)
}
