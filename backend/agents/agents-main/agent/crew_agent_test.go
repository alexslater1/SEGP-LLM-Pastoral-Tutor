package agent

import (
	"testing"

	"github.com/segp/agents-main/clock"
	"github.com/segp/agents-main/entity"
	"github.com/segp/agents-main/history"
	"github.com/segp/agents-main/knowledge"
	"github.com/segp/agents-main/llm"
	"github.com/segp/agents-main/tools"
	"github.com/stretchr/testify/assert"
)

func TestCrewAgent(t *testing.T) {
	fastAgent := newFastAgent("test", "test", "test", []string{entity.UserEntityId}, tools.NewToolHandler([]tools.Tool{}), llm.NewMockLLM(), knowledge.NewLocalKnowledge(), clock.NewMockClock(), history.NewLocalHistory())
	crewAgent := NewCrewAgent(fastAgent, []string{"entity_id_1", "entity_id_2"})

	assert.Equal(t, []string{"entity_id_1", "entity_id_2", "user"}, crewAgent.canOffloadToEntities())
}
