package tools

import (
	"fmt"
	"os"
	"testing"
)

func TestGetCurrentDate(t *testing.T) {
	if os.Getenv("CICD") == "true" {
		t.Skip("Skipping test in CICD")
	}

	date := GetCurrentDate()
	fmt.Println(date)
}
