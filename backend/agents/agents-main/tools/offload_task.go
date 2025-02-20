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
		Description: "This tool is used to give a task to another entity, who you believe will be able to help solve a part of the current problem. You may only give one task to the entity. Be extremely descriptive about the task, including any context you currently have which would be useful to the entity. If this is the only remaining task to do, you should set would_like_the_answer to false. If something else should be done (either by you or by another entity), you should set would_like_the_answer to true.",
		Parameters: []Parameter{
			{Name: "entity_id", Description: "The ID of the entity to give the task to. Must be from the supplied list of entities", Type: ParameterTypeString},
			{Name: "task", Description: "The query to give to the new entity. This should be a single, finegrained task which the new entity can complete. If you want to give multiple tasks, you should use multiple calls to this tool.", Type: ParameterTypeString},
			{Name: "something_else_should_be_done", Description: "Whether there exists another task which should be done (either by you or by another entity) other than that which you are giving to the new entity (true/false)", Type: ParameterTypeBoolean},
		},
	}
}
