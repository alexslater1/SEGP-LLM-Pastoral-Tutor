package tools

import "fmt"

type AvailableEntity struct {
	ID          string
	Description string
}

type OffloadTaskTool struct {
	AvailableEntities []AvailableEntity
}

func NewOffloadTaskTool(availableEntities []AvailableEntity) *OffloadTaskTool {
	return &OffloadTaskTool{AvailableEntities: availableEntities}
}

func (s *OffloadTaskTool) Definition() ToolDefinition {
	return ToolDefinition{
		Name:        "offload_task",
		Description: fmt.Sprintf("This tool is used to give a task to another entity, who you believe will be able to help solve the current problem. Here are the available entities you have to choose from: %s", s.availableEntitiesString()),
		Parameters: []Parameter{
			{Name: "entity_id", Description: "The ID of the entity to give the task to", Type: ParameterTypeString},
			{Name: "task", Description: "The exact query or task to give to the entity", Type: ParameterTypeString},
		},
	}
}

func (s *OffloadTaskTool) availableEntitiesString() string {
	availableEntitiesStr := ""
	for _, availableEntity := range s.AvailableEntities {
		availableEntitiesStr += fmt.Sprintf("EntityID: %s, Description: %s\n", availableEntity.ID, availableEntity.Description)
	}
	return availableEntitiesStr
}
