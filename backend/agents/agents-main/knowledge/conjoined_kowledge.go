package knowledge

import "fmt"

type ConjoinedKnowledge struct {
	Knowledges []Knowledge
}

func NewConjoinedKnowledge(knowledges ...Knowledge) *ConjoinedKnowledge {
	return &ConjoinedKnowledge{
		Knowledges: knowledges,
	}
}

func (c *ConjoinedKnowledge) Get(query string) (*string, error) {
	resultStr := ""
	for _, knowledge := range c.Knowledges {
		result, err := knowledge.Get(query)
		if err == nil {
			return result, nil
		}
		resultStr += fmt.Sprintf("Source: %s\nData:%s\n\n", knowledge.Name(), *result)
	}

	return &resultStr, nil
}

func (c *ConjoinedKnowledge) Name() string {
	return "conjoinedknowledge"
}
