package tools

type AvailableEntity struct {
	ID          string
	Description string
}

type OffloadTaskTool struct {
}

func NewOffloadTaskTool() *OffloadTaskTool {
	return &OffloadTaskTool{}
}

func (s *OffloadTaskTool) Definition() ToolDefinition {
	return ToolDefinition{
		Name:        "offload_task",
		Description: "This tool is used to give a task to another entity, who you believe will be able to help solve the current problem. This will exit your agent loop, and pass your task onto the new entity. Note that the entity may be able to pass it back to you after.",
		Parameters: []Parameter{
			{Name: "entity_id", Description: "The ID of the entity to give the task to. Must be from the supplied list of entities", Type: ParameterTypeString},
			{Name: "task", Description: "The exact query or task to give to the entity", Type: ParameterTypeString},
		},
	}
}
