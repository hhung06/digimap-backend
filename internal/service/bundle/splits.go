package bundle

import (
	"encoding/json"

	"github.com/hhung06/digimap-backend/internal/domain"
)

// SplitBundles holds the four viewer-fetched artifacts produced from the core-SDK draft.
// Each field is a map ready to be gzip-encrypted and uploaded as {env}/{type}/public/{venue}.digimap.
type SplitBundles struct {
	Base           map[string]any
	Overview       map[string]any
	LocationSimple map[string]any
	Metadata       map[string]any
}

// AssembleSplitBundles transforms a core-SDK draft (keys: information, levels, mapGroups,
// paths, meshes, locations, connections, beacons, categories, images, theme, connectionsList)
// into the four viewer-fetched bundles. theme is the resolved venue theme (DB or draft fallback).
//
// Viewer loaders (fromStructure/fromPaths/fromLocationSimple/fromMetadata) read exactly the
// keys placed in each bundle; missing draft keys default to []any{} so .map() calls don't throw.
func AssembleSplitBundles(draft map[string]any, theme any) SplitBundles {
	return SplitBundles{
		Base: map[string]any{
			"meshes": orAnyEmpty(draft["meshes"]),
		},
		Overview: map[string]any{
			"paths":     orAnyEmpty(draft["paths"]),
			"theme":     theme,
			"mapGroups": orAnyEmpty(draft["mapGroups"]),
			"levels":    orAnyEmpty(draft["levels"]),
		},
		LocationSimple: map[string]any{
			"locationMinimal": orAnyEmpty(draft["locations"]),
		},
		Metadata: map[string]any{
			"information": draft["information"],
			"categories":  orAnyEmpty(draft["categories"]),
			"locations":   orAnyEmpty(draft["locations"]),
			"connections": orAnyEmpty(draft["connections"]),
			"beacons":     orAnyEmpty(draft["beacons"]),
			"images":      orAnyEmpty(draft["images"]),
		},
	}
}

// AssembleLocalizedOverlay produces the per-language overlay the viewer fetches for non-English
// locales: {locations:[{id,...localized}], categories:[{id,...localized}]}.
// The viewer merges each entry onto the corresponding base record by id via Object.assign.
func AssembleLocalizedOverlay(
	locations []*domain.Location,
	cats []*domain.LocationCategory,
) map[string]any {
	locs := make([]map[string]any, 0, len(locations))
	for _, l := range locations {
		var loc map[string]any
		if err := json.Unmarshal(l.Localization, &loc); err == nil && loc != nil {
			entry := map[string]any{"id": l.ID.String()}
			for k, v := range loc {
				entry[k] = v
			}
			locs = append(locs, entry)
		}
	}

	catList := make([]map[string]any, 0, len(cats))
	for _, c := range cats {
		var loc map[string]any
		if err := json.Unmarshal(c.Localization, &loc); err == nil && loc != nil {
			entry := map[string]any{"id": c.ID.String()}
			for k, v := range loc {
				entry[k] = v
			}
			catList = append(catList, entry)
		}
	}

	return map[string]any{
		"locations":  locs,
		"categories": catList,
	}
}

// orAnyEmpty returns v if non-nil, otherwise an empty slice. Prevents viewer .map() from
// throwing when a draft field is missing.
func orAnyEmpty(v any) any {
	if v == nil {
		return []any{}
	}
	return v
}
