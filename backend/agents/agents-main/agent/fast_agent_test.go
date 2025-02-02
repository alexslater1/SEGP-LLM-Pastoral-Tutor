package agent

import (
	"testing"

	googleSearch "github.com/segp/agents-main/google_search"
	"github.com/segp/agents-main/knowledge"
	"github.com/segp/agents-main/tools"
	"github.com/stretchr/testify/assert"
)

func TestToolCallChoiceString(t *testing.T) {
	var (
		knowledge   = knowledge.NewLocalKnowledge()
		toolHandler = tools.NewDefaultToolHandler(googleSearch.NewMockGoogleSearchClient(), knowledge)
	)

}
