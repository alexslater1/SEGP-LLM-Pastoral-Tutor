package ragger

type Ragger interface {
	ChunksFrom(text string) ([]string, error)
	ContactsFrom(text string) ([]Contact, error)

	EmbeddingsFor(text string) ([]float32, error)
	EmbeddingsForAll(texts []string) ([][]float32, error)
}

type ContactType string

const (
	ContactTypeEmail   ContactType = "email"
	ContactTypePhone   ContactType = "phone"
	ContactTypeWebsite ContactType = "website"
)

type Contact struct {
	Value   string
	Context string
	Type    ContactType
}
