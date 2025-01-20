package tools

import (
	"os"
	"testing"

	googleSearch "github.com/segp/agents-main/google_search"
)

func TestCallGoogleSearchUrlTool(t *testing.T) {
	if os.Getenv("CICD") == "True" {
		t.Skip("Skipping test due to CICDs")
	}

	tool := NewGoogleSearchUrlTool(googleSearch.NewRodClient())
	result, err := tool.PageContentFor("https://www.google.com/search?q=hello&oq=hello&gs_lcrp=EgZjaHJvbWUqDggAEEUYJxg7GIAEGIoFMg4IABBFGCcYOxiABBiKBTITCAEQLhiDARjHARixAxjRAxiABDIKCAIQLhixAxiABDIKCAMQABixAxiABDINCAQQABiDARixAxiABDIKCAUQABixAxiABDIGOAYQRRg8MgYIBxBFGDzSAQgxNDg0ajBqNKgCALACAA&sourceid=chrome&ie=UTF-8")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(*result)
}
