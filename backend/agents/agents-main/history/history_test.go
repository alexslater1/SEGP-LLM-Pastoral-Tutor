package history

import (
	"fmt"
	"os"
	"testing"

	"github.com/segp/agents-main/storage"
	"github.com/segp/agents-main/utils"
)

func TestHistory(t *testing.T) {
	var (
		storage = storage.NewSupabaseStorage(utils.Required(os.Getenv("SUPABASE_URL"), "SUPABASE_URL"), utils.Required(os.Getenv("SUPABASE_SERVICE_KEY"), "SUPABASE_SERVICE_KEY"))
		history = NewAgentEventHistory(storage)
	)

	h, err := history.GetMessageHistory("7fc63fda-c955-4dd6-bbb8-9166fa95e53a")
	if err != nil {
		t.Fatalf("error getting message history: %v", err)
	}

	fmt.Printf("%+v\n", h)
}
