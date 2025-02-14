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
		Description: "This tool is used to give a task to another entity, who you believe will be able to help solve the current problem. This will exit your agent loop immediately and you wont be able to perform any more tasks. Thus, you must pass all aditional context to the next entity via the task, so no information is lost.",
		Parameters: []Parameter{
			{Name: "entity_id", Description: "The ID of the entity to give the task to. Must be from the supplied list of entities", Type: ParameterTypeString},
			{Name: "task", Description: "The query or task to give to the new entity, including any context which is relevant to the task (which has yet to be solved)", Type: ParameterTypeString},
		},
	}
}
