package search

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestVenueIndexMapping(t *testing.T) {
	m := VenueIndexMapping()

	// Marshal to JSON and back — ensures no un-serialisable values.
	b, err := json.Marshal(m)
	if err != nil {
		t.Fatalf("json.Marshal: %v", err)
	}

	var out map[string]any
	if err := json.Unmarshal(b, &out); err != nil {
		t.Fatalf("json.Unmarshal: %v", err)
	}

	// settings.number_of_shards == 3
	settings, ok := out["settings"].(map[string]any)
	if !ok {
		t.Fatal("settings is missing or wrong type")
	}
	shards, ok := settings["number_of_shards"]
	if !ok {
		t.Fatal("settings.number_of_shards is missing")
	}
	// JSON numbers unmarshal as float64.
	if shards.(float64) != 3 {
		t.Errorf("settings.number_of_shards: got %v, want 3", shards)
	}

	// mappings.properties.type is present
	mappings, ok := out["mappings"].(map[string]any)
	if !ok {
		t.Fatal("mappings is missing or wrong type")
	}
	properties, ok := mappings["properties"].(map[string]any)
	if !ok {
		t.Fatal("mappings.properties is missing or wrong type")
	}
	if _, ok := properties["type"]; !ok {
		t.Fatal("mappings.properties.type is missing")
	}

	// mappings.properties.categories.type == "nested"
	categories, ok := properties["categories"].(map[string]any)
	if !ok {
		t.Fatal("mappings.properties.categories is missing or wrong type")
	}
	catType, ok := categories["type"].(string)
	if !ok {
		t.Fatal("mappings.properties.categories.type is missing or wrong type")
	}
	if catType != "nested" {
		t.Errorf("mappings.properties.categories.type: got %q, want \"nested\"", catType)
	}
}

func TestBulkItemsError_ReturnsErrorWhenAnyItemFailed(t *testing.T) {
	body := []byte(`{
		"errors": true,
		"items": [
			{"index": {"_id": "ok", "status": 201}},
			{"index": {"_id": "bad", "status": 400, "error": {"type": "mapper_parsing_exception", "reason": "bad field"}}}
		]
	}`)

	err := bulkItemsError(body)
	if err == nil {
		t.Fatal("expected bulk item error")
	}
	if !strings.Contains(err.Error(), "bad") || !strings.Contains(err.Error(), "mapper_parsing_exception") {
		t.Fatalf("expected item id and error type, got %v", err)
	}
}

func TestBulkItemsError_IgnoresSuccessfulBulkResponse(t *testing.T) {
	body := []byte(`{"errors": false, "items": [{"index": {"_id": "ok", "status": 201}}]}`)

	if err := bulkItemsError(body); err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
}
