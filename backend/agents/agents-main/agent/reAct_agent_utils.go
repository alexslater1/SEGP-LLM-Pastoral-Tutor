package agent

import (
	"fmt"
	"log/slog"
)

type ReActAgentStage string

const (
	ReActAgentStageThinking       ReActAgentStage = "thinking"
	ReActAgentStageToolCall       ReActAgentStage = "tool_call"
	ReActAgentStageToolCallResult ReActAgentStage = "tool_call_result"
	ReActAgentStageObservation    ReActAgentStage = "observation"
)

var colorCodes = map[string]string{
	"red":     "\033[31m",
	"green":   "\033[32m",
	"yellow":  "\033[33m",
	"blue":    "\033[34m",
	"magenta": "\033[35m",
	"cyan":    "\033[36m",
	"white":   "\033[37m",
	"reset":   "\033[0m",
}

var colorMap = map[ReActAgentStage]string{
	ReActAgentStageThinking:       colorCodes["yellow"],
	ReActAgentStageToolCall:       colorCodes["cyan"],
	ReActAgentStageToolCallResult: colorCodes["green"],
	ReActAgentStageObservation:    colorCodes["green"],
}

func logReActStage(value string, stage ReActAgentStage) {
	slog.Info(string(stage), "value", value)
	fmt.Println()
}
