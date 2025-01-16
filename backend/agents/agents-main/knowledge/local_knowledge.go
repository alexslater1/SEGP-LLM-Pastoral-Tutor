package knowledge

type LocalKnowledge struct {
}

func (k *LocalKnowledge) Get(query string) (*string, error) {
	knowledge := ""
	return &knowledge, nil
}
