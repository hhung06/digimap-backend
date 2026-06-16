// Package localization provides helpers for reading the flat localization blob
// convention used throughout the codebase.
//
// The on-disk shape is a flat JSON object with keys of the form "field_lang",
// e.g. {"common_name_en": "Acme Corp", "common_name_ja": "アクメ社"}.
// Location blobs use the "common_" prefix; Product and Category blobs use bare
// field names ("name_en", "name_ja", etc.).
package localization

import (
	"encoding/json"
)

// DecodeBlob unmarshals a JSON blob into a flat map. It is tolerant of nil and
// invalid input — both return an empty (non-nil) map so callers can safely call
// Field without a nil check.
func DecodeBlob(raw json.RawMessage) map[string]any {
	if len(raw) == 0 {
		return map[string]any{}
	}
	var m map[string]any
	if err := json.Unmarshal(raw, &m); err != nil {
		return map[string]any{}
	}
	return m
}

// Field reads a flat localization value using the convention "field_lang".
// It returns fallback when the key is absent, empty, or not a string.
func Field(blob map[string]any, field, lang, fallback string) string {
	if blob == nil {
		return fallback
	}
	v, ok := blob[field+"_"+lang]
	if !ok {
		return fallback
	}
	s, ok := v.(string)
	if !ok || s == "" {
		return fallback
	}
	return s
}
