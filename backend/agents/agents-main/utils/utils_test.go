package utils

import (
	"errors"
	"slices"
	"strconv"
	"testing"
	"time"
)

func TestRequired(t *testing.T) {
	// Test valid string
	result := Required("test", "stringParam")
	if result != "test" {
		t.Errorf("Expected 'test', got '%s'", result)
	}

	// Test valid int
	intResult := Required(42, "intParam")
	if intResult != 42 {
		t.Errorf("Expected 42, got %d", intResult)
	}

	// Test panic with empty string
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for empty string, but nothing happened")
		}
	}()
	Required("", "emptyString")
}

func TestDoAsync(t *testing.T) {
	// Test successful async operation
	task := DoAsync(func() (string, error) {
		return "success", nil
	})

	result, err := task.Get()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result != "success" {
		t.Errorf("Expected 'success', got '%s'", result)
	}

	// Test error case
	expectedErr := errors.New("test error")
	task = DoAsync(func() (string, error) {
		return "", expectedErr
	})

	result, err = task.Get()
	if err != expectedErr {
		t.Errorf("Expected error %v, got %v", expectedErr, err)
	}
	if result != "" {
		t.Errorf("Expected empty string, got '%s'", result)
	}
}

func TestDoAsyncList(t *testing.T) {
	numbers := []int{1, 2, 3}

	// Test successful case
	tasks := DoAsyncList(numbers, func(n int) (int, error) {
		return n * 2, nil
	})

	results, err := GetAsyncList(tasks)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	expected := []int{2, 4, 6}
	for i, v := range results {
		if v != expected[i] {
			t.Errorf("Expected %d at index %d, got %d", expected[i], i, v)
		}
	}

	// Test error case
	expectedErr := errors.New("test error")
	tasks = DoAsyncList(numbers, func(n int) (int, error) {
		if n == 2 {
			return 0, expectedErr
		}
		return n * 2, nil
	})

	results, err = GetAsyncList(tasks)
	if err != expectedErr {
		t.Errorf("Expected error %v, got %v", expectedErr, err)
	}
	if results != nil {
		t.Errorf("Expected nil results, got %v", results)
	}
}

func TestConcurrency(t *testing.T) {
	start := time.Now()

	numbers := []int{1, 2, 3}
	tasks := DoAsyncList(numbers, func(n int) (int, error) {
		time.Sleep(100 * time.Millisecond)
		return n * 2, nil
	})

	_, err := GetAsyncList(tasks)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	duration := time.Since(start)
	if duration >= 300*time.Millisecond {
		t.Errorf("Expected concurrent execution under 300ms, took %v", duration)
	}
}

func TestRemoveDuplicates(t *testing.T) {
	slice := []int{1, 2, 2, 3, 4, 4, 5}
	result := RemoveDuplicates(slice)
	if !slices.Equal(result, []int{1, 2, 3, 4, 5}) {
		t.Errorf("Expected [1, 2, 3, 4, 5], got %v", result)
	}
}

func TestSorted(t *testing.T) {
	slice := []int{3, 1, 4, 1, 5, 9, 2, 6, 5, 3, 5}
	result := Sorted(slice, func(n int) string {
		return strconv.Itoa(n)
	})
	if !slices.Equal(result, []int{1, 1, 2, 3, 3, 4, 5, 5, 5, 6, 9}) {
		t.Errorf("Expected [1, 1, 2, 3, 3, 4, 5, 5, 5, 6, 9], got %v", result)
	}
}
