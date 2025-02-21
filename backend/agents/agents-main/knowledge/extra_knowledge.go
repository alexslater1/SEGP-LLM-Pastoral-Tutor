package knowledge

import (
	"fmt"
)

type ExtraKnowledge struct {
	Knowledge []interface{}
}

func NewExtraKnowledge(knowledge ...interface{}) *ExtraKnowledge {
	return &ExtraKnowledge{
		Knowledge: knowledge,
	}
}

func (e *ExtraKnowledge) Get(query string) (*string, error) {
	knowledgeStr := ""
	for _, knowledge := range e.Knowledge {
		knowledgeStr += fmt.Sprintf("%+v\n", knowledge)
	}

	return &knowledgeStr, nil
}

func (e *ExtraKnowledge) Name() string {
	return "extraknowledge"
}
