package localization_test

import (
	"encoding/json"
	"testing"

	"github.com/hhung06/digimap-backend/internal/localization"
)

func TestDecodeBlob(t *testing.T) {
	t.Run("flat blob", func(t *testing.T) {
		raw := json.RawMessage(`{"common_name_en":"Acme","common_name_ja":"アクメ"}`)
		m := localization.DecodeBlob(raw)
		if m["common_name_en"] != "Acme" {
			t.Fatalf("expected Acme, got %v", m["common_name_en"])
		}
	})

	t.Run("nil input returns empty map", func(t *testing.T) {
		m := localization.DecodeBlob(nil)
		if m == nil || len(m) != 0 {
			t.Fatalf("expected empty map, got %v", m)
		}
	})

	t.Run("empty input returns empty map", func(t *testing.T) {
		m := localization.DecodeBlob(json.RawMessage{})
		if m == nil || len(m) != 0 {
			t.Fatalf("expected empty map, got %v", m)
		}
	})

	t.Run("invalid JSON returns empty map", func(t *testing.T) {
		m := localization.DecodeBlob(json.RawMessage(`not-json`))
		if m == nil || len(m) != 0 {
			t.Fatalf("expected empty map, got %v", m)
		}
	})
}

func TestField(t *testing.T) {
	blob := map[string]any{
		"common_name_en": "Acme Corp",
		"common_name_ja": "アクメ社",
		"name_en":        "Widget",
	}

	t.Run("hit", func(t *testing.T) {
		got := localization.Field(blob, "common_name", "en", "fallback")
		if got != "Acme Corp" {
			t.Fatalf("expected Acme Corp, got %q", got)
		}
	})

	t.Run("miss returns fallback", func(t *testing.T) {
		got := localization.Field(blob, "common_name", "zh", "fallback")
		if got != "fallback" {
			t.Fatalf("expected fallback, got %q", got)
		}
	})

	t.Run("empty string returns fallback", func(t *testing.T) {
		b := map[string]any{"name_en": ""}
		got := localization.Field(b, "name", "en", "fb")
		if got != "fb" {
			t.Fatalf("expected fb, got %q", got)
		}
	})

	t.Run("non-string value returns fallback", func(t *testing.T) {
		b := map[string]any{"name_en": 42}
		got := localization.Field(b, "name", "en", "fb")
		if got != "fb" {
			t.Fatalf("expected fb, got %q", got)
		}
	})

	t.Run("nil blob returns fallback", func(t *testing.T) {
		got := localization.Field(nil, "name", "en", "fb")
		if got != "fb" {
			t.Fatalf("expected fb, got %q", got)
		}
	})
}
