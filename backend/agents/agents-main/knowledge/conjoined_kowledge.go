package knowledge

import (
	"fmt"
	"log"
	"log/slog"
)

type ConjoinedKnowledge struct {
	Knowledges []Knowledge
}

func NewConjoinedKnowledge(knowledges ...Knowledge) *ConjoinedKnowledge {
	return &ConjoinedKnowledge{
		Knowledges: knowledges,
	}
}

func (c *ConjoinedKnowledge) Get(query string) (*string, error) {
	slog.Info("Getting knowledge", "query", query, "knowledges", c.Knowledges)

	resultStr := ""
	for i, knowledge := range c.Knowledges {
		log.Println(i)
		result, err := knowledge.Get(query)
		if err != nil {
			return nil, err
		}
		resultStr += fmt.Sprintf("Source: %s\nData:%s\n\n", knowledge.Name(), *result)
	}

	return &resultStr, nil
}

func (c *ConjoinedKnowledge) Name() string {
	return "conjoinedknowledge"
}
