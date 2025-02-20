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

type testEntity struct {
	id          string
	description string
}

func newTestEntity(id string, description string) entity.Entity {
	return &testEntity{id: id, description: description}
}

func (e *testEntity) Id() string {
	return e.id
}

func (e *testEntity) Description() string {
	return e.description
}

func TestCrewAgent(t *testing.T) {

	fastAgent := newFastAgent("test", "test", "test", tools.NewToolHandler([]tools.Tool{}), llm.NewMockLLM(), knowledge.NewLocalKnowledge(), clock.NewMockClock(), history.NewLocalHistory(), entity.UserEntity)

	crewAgent := NewCrewAgent(fastAgent)
	assert.Equal(t, []entity.Entity{entity.UserEntity}, crewAgent.canOffloadToEntities())

	crewAgent.AddCanOffloadToEntity(newTestEntity("entity_id_1", "test"), newTestEntity("entity_id_2", "test"))
	assert.Equal(t, []entity.Entity{newTestEntity("entity_id_1", "test"), newTestEntity("entity_id_2", "test"), entity.UserEntity}, crewAgent.canOffloadToEntities())
}
