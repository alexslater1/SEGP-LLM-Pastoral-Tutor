package utils

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

func Required[T any](value T, name string) T {
	if reflect.ValueOf(value).IsZero() {
		panic(fmt.Sprintf("%s is required", name))
	}
	return value
}

type Task[T any] struct {
	ch      chan T
	errorCh chan error
}

func DoAsync[T any](fn func() (T, error)) *Task[T] {
	ch := make(chan T)
	errorCh := make(chan error)

	go func() {
		result, err := fn()
		if err != nil {
			errorCh <- err
			return
		}
		ch <- result
	}()

	return &Task[T]{ch, errorCh}
}

func DoAsyncList[T any, U any](items []T, fn func(T) (U, error)) []*Task[U] {
	tasks := make([]*Task[U], len(items))

	for i, item := range items {
		tasks[i] = DoAsync(func() (U, error) {
			return fn(item)
		})
	}

	return tasks
}

func (task *Task[T]) Get() (T, error) {
	var zero T // This will initialize `zero` to the zero value for type T
	select {
	case result := <-task.ch:
		return result, nil
	case err := <-task.errorCh:
		return zero, err
	}
}

func GetAsyncList[T any](tasks []*Task[T]) ([]T, error) {
	results := make([]T, len(tasks))

	for i, task := range tasks {
		result, err := task.Get()
		if err != nil {
			return nil, err
		}
		results[i] = result
	}

	return results, nil
}

func CleanText(text string) string {
	// Split into lines, trim each line, and handle multiple newlines
	lines := strings.Split(text, "\n")
	var cleanedLines []string

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed != "" {
			cleanedLines = append(cleanedLines, trimmed)
		}
	}

	// Join with double newlines and clean up any remaining multiple newlines
	text = strings.Join(cleanedLines, "\n")
	re := regexp.MustCompile(`\n\s*\n`)
	text = re.ReplaceAllString(text, "\n\n")

	// Remove zero-width characters
	text = strings.ReplaceAll(text, "\u200c", "") // Remove zero-width non-joiner
	text = strings.ReplaceAll(text, "\u200b", "") // Remove zero-width space

	// Handle escape sequences
	text = strings.ReplaceAll(text, "\\n", "\n")
	text = strings.ReplaceAll(text, "\\\"", "\"")
	text = strings.ReplaceAll(text, "\\\\", "\\")

	// Ensure text doesn't end with a partial escape sequence
	text = strings.TrimSuffix(text, "\\")

	return text
}
