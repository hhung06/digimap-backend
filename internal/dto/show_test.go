package dto

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestResponseShapes(t *testing.T) {
	b, _ := json.MarshalIndent(OK(map[string]string{"status": "ok"}), "", "  ")
	fmt.Println("Success:"); fmt.Println(string(b))

	items := []map[string]string{{"id": "abc", "name": "Venue 1"}}
	b, _ = json.MarshalIndent(Paginated(items, 42, 1, 20), "", "  ")
	fmt.Println("\nPaginated:"); fmt.Println(string(b))

	b, _ = json.MarshalIndent(Fail(CodeNotFound, "venue not found"), "", "  ")
	fmt.Println("\nNot Found:"); fmt.Println(string(b))

	b, _ = json.MarshalIndent(FailMessages(CodeValidationError, []string{"name: required", "email: invalid"}), "", "  ")
	fmt.Println("\nValidation:"); fmt.Println(string(b))
}
