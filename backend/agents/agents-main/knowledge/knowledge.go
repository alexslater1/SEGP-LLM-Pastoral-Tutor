package knowledge

type Knowledge interface {
	Get(query string) (*string, error)
	Name() string
}
