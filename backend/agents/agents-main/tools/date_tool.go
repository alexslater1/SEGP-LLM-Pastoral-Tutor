package tools

import "time"

type DateTool struct{}

func (d *DateTool) Definition() ToolDefinition {
	return ToolDefinition{
		Name:        "date",
		Description: "Get the current date in YYYY-MM-DD format",
	}
}

func (d *DateTool) GetCurrentDate() string {
	return time.Now().Format("2006-01-02")
}
