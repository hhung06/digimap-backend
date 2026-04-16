package enricher

import "encoding/json"

// MergeInto merges extras on top of the base DTO and returns the result as any.
// Returns base unchanged (no JSON round-trip cost) if extras is nil or empty.
func MergeInto(base any, extras map[string]any) any {
	if len(extras) == 0 {
		return base
	}
	b, err := json.Marshal(base)
	if err != nil {
		return base
	}
	var m map[string]any
	if err := json.Unmarshal(b, &m); err != nil {
		return base
	}
	for k, v := range extras {
		m[k] = v
	}
	return m
}
