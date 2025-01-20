package main

import (
	"fmt"
)

type ReActAgentStage string

const (
	ReActAgentStageThinking    ReActAgentStage = "thinking"
	ReActAgentStageToolCall    ReActAgentStage = "tool_call"
	ReActAgentStageObservation ReActAgentStage = "observation"
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
	ReActAgentStageThinking:    colorCodes["yellow"],
	ReActAgentStageToolCall:    colorCodes["cyan"],
	ReActAgentStageObservation: colorCodes["green"],
}

func logReActStage(arg string, message string, stage ReActAgentStage) {
	fmt.Printf("%s%s%s", colorMap[stage], message, colorCodes["reset"])
}

func main() {
	logReActStage("test", "test", ReActAgentStageThinking)
}
