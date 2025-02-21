package knowledge

type LocalKnowledge struct {
}

func NewLocalKnowledge() *LocalKnowledge {
	return &LocalKnowledge{}
}

func (k *LocalKnowledge) Get(query string) (*string, error) {
	knowledge := ""
	return &knowledge, nil
}

func (k *LocalKnowledge) Name() string {
	return "local"
}
