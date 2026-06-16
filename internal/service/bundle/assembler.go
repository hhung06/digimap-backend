// Package bundle implements the v2 language-fanned snapshot publisher.
// Each method mirrors Django's prepare_* helpers in publish_venue_v2.py.
package bundle

import (
	"encoding/json"
	"strings"

	"github.com/google/uuid"

	"github.com/hhung06/digimap-backend/internal/domain"
	"github.com/hhung06/digimap-backend/internal/localization"
)

// AssembleLanguageBundle builds the six-key per-language bundle dict that gets encrypted.
// Mirrors Django's combine_data (publish_venue_v2.py:724).
//
// metadata is the raw map decoded from the existing snapshot draft JSON (from S3).
// It serves as the base for the "map" key (with per-language localization layered on top).
func AssembleLanguageBundle(
	venue *domain.Venue,
	lang string,
	locations []*domain.Location,
	locationCats []*domain.LocationCategory,
	products []*domain.Product,
	productCats []*domain.ProductCategory,
	metadata map[string]any,
) map[string]any {
	return map[string]any{
		"search":                   prepareSearchData(lang, locations, products),
		"map":                      prepareMapData(lang, locations, locationCats, metadata),
		"products":                 prepareProductData(lang, products),
		"productsGroupedByCountry": prepareProductsGroupedByCountry(products),
		"exhibitorSearchOptions":   prepareExhibitorSearchOptions(venue, lang, locations, locationCats),
		"productSearchOptions":     prepareProductSearchOptions(venue, lang, locationCats, productCats),
	}
}

// ── Search ────────────────────────────────────────────────────────────────────

type searchLocation struct {
	ID              string          `json:"id"`
	Name            string          `json:"name"`
	BoothNumber     string          `json:"booth_number,omitempty"`
	Custom          map[string]any  `json:"custom,omitempty"`
	CommonLocationType int          `json:"common_location_type"`
}

type searchProduct struct {
	ID     string         `json:"id"`
	Name   string         `json:"name"`
	Custom map[string]any `json:"custom,omitempty"`
}

// prepareSearchData mirrors publish_venue_v2.py:338 prepare_search_data.
// NOTE: Japanese dictionary tokenization (janome + wanakana) is not yet implemented.
// The `dictionaries` field is always empty in this port. Track as a known gap.
func prepareSearchData(lang string, locations []*domain.Location, products []*domain.Product) map[string]any {
	isJa := lang == "ja"

	var locs []searchLocation
	for _, l := range locations {
		custom := decodeCustom(l.Custom)
		// Non-ja: only include exhibitor_english_status == 1 (Django line ~365)
		if !isJa && customInt(custom, "exhibitor_english_status") != 1 {
			continue
		}
		name := localizedName(l.Localization, l.CommonName, lang)
		locs = append(locs, searchLocation{
			ID:                 l.ID.String(),
			Name:               name,
			BoothNumber:        l.BoothNumber,
			Custom:             custom,
			CommonLocationType: l.CommonLocationType,
		})
	}

	var prods []searchProduct
	for _, p := range products {
		custom := decodeCustom(p.Custom)
		if !isJa && customInt(custom, "exhibitor_english_status") != 1 {
			continue
		}
		name := localizedProductName(p.Localization, p.Name, lang)
		prods = append(prods, searchProduct{
			ID:     p.ID.String(),
			Name:   name,
			Custom: custom,
		})
	}

	// NOTE: dictionaries (Japanese morphological tokenization) not yet ported.
	// Requires github.com/ikawaha/kagome + kana conversion. Tracked separately.
	return map[string]any{
		"locations":   orEmpty(locs),
		"products":    orEmpty(prods),
		"dictionaries": map[string]any{},
	}
}

// ── Map ───────────────────────────────────────────────────────────────────────

// prepareMapData mirrors publish_venue_v2.py:284 prepare_map_data.
// metadata is the decoded snapshot JSON (from S3 draft) which already contains
// the venue map structure (levels, location geometry, etc.).
// Per-language localization is merged on top for matching location IDs.
func prepareMapData(
	lang string,
	locations []*domain.Location,
	cats []*domain.LocationCategory,
	metadata map[string]any,
) map[string]any {
	// Build a localization lookup by location ID.
	locByID := make(map[string]map[string]any, len(locations))
	for _, l := range locations {
		var loc map[string]any
		if err := json.Unmarshal(l.Localization, &loc); err == nil && loc != nil {
			locByID[l.ID.String()] = loc
		}
	}

	// Build per-language category localization list.
	catLocalizations := make([]map[string]any, 0, len(cats))
	for _, c := range cats {
		var loc map[string]any
		if err := json.Unmarshal(c.Localization, &loc); err == nil && loc != nil {
			entry := map[string]any{"id": c.ID.String()}
			for k, v := range loc {
				entry[k] = v
			}
			catLocalizations = append(catLocalizations, entry)
		}
	}

	// Merge location localization into the metadata locations array if present.
	// Django lines 313-318: for each bundle location, overlay the detail serializer fields.
	if meta, ok := metadata["locations"].([]any); ok {
		for _, raw := range meta {
			if m, ok := raw.(map[string]any); ok {
				if id, ok := m["id"].(string); ok {
					if langData, ok := locByID[id]; ok {
						for k, v := range langData {
							m[k] = v
						}
					}
				}
			}
		}
	}

	return map[string]any{
		"data": metadata,
		"localization": map[string]any{
			"locations":  locByID,
			"categories": catLocalizations,
		},
	}
}

// ── Products ──────────────────────────────────────────────────────────────────

type productItem struct {
	ID           string         `json:"id"`
	Name         string         `json:"name"`
	LocationID   string         `json:"location_id,omitempty"`
	CategoryID   string         `json:"category_id,omitempty"`
	Custom       map[string]any `json:"custom,omitempty"`
	Localization map[string]any `json:"localization,omitempty"`
}

// prepareProductData mirrors publish_venue_v2.py:567 prepare_product_data.
// Section filter: ja → section==1, else → section==2 AND exhibitor_english_status==1.
func prepareProductData(lang string, products []*domain.Product) []productItem {
	isJa := lang == "ja"
	var out []productItem
	for _, p := range products {
		custom := decodeCustom(p.Custom)
		section := customInt(custom, "section")
		engStatus := customInt(custom, "exhibitor_english_status")
		if isJa && section != 1 {
			continue
		}
		if !isJa && (section != 2 || engStatus != 1) {
			continue
		}
		item := productItem{
			ID:     p.ID.String(),
			Name:   localizedProductName(p.Localization, p.Name, lang),
			Custom: stripNilKeys(custom),
		}
		if p.LocationID != nil {
			item.LocationID = p.LocationID.String()
		}
		if p.MainCategoryID != nil {
			item.CategoryID = p.MainCategoryID.String()
		}
		var locMap map[string]any
		if err := json.Unmarshal(p.Localization, &locMap); err == nil {
			item.Localization = locMap
		}
		out = append(out, item)
	}
	return out
}

// ── Products grouped by country ───────────────────────────────────────────────

type countryGroup struct {
	CountryName   string   `json:"country_name"`
	CountryNameEn string   `json:"country_name_en"`
	Products      []string `json:"products"` // product UUIDs
}

// prepareProductsGroupedByCountry mirrors publish_venue_v2.py:599.
// Filters to source=="external", groups by (country_name, country_name_en) from product.location.custom.
func prepareProductsGroupedByCountry(products []*domain.Product) []countryGroup {
	type key struct{ name, nameEn string }
	groups := make(map[key][]string)
	order := []key{}
	seen := map[key]bool{}

	for _, p := range products {
		if p.Source != "external" {
			continue
		}
		custom := decodeCustom(p.Custom)
		name := customStr(custom, "country_name")
		nameEn := customStr(custom, "country_name_en")
		k := key{name, nameEn}
		if !seen[k] {
			seen[k] = true
			order = append(order, k)
		}
		groups[k] = append(groups[k], p.ID.String())
	}

	out := make([]countryGroup, 0, len(order))
	for _, k := range order {
		out = append(out, countryGroup{
			CountryName:   k.name,
			CountryNameEn: k.nameEn,
			Products:      groups[k],
		})
	}
	return out
}

// ── Exhibitor search options ──────────────────────────────────────────────────

// prepareExhibitorSearchOptions mirrors publish_venue_v2.py:154.
func prepareExhibitorSearchOptions(
	venue *domain.Venue,
	lang string,
	locations []*domain.Location,
	cats []*domain.LocationCategory,
) map[string]any {
	exportMap := SelectExportableMapping(venue.ExternalID)

	// Build exhibitor_categories from location categories (external only).
	catItems := make([]map[string]any, 0, len(cats))
	for _, c := range cats {
		if c.Source != "external" {
			continue
		}
		catItems = append(catItems, map[string]any{
			"id":   c.ID.String(),
			"name": localizedCatName(c.Localization, c.Name, lang),
		})
	}

	// Build exhibitor_english_status options from exportMap.
	isJa := lang == "ja"
	statusOpts := make([]map[string]any, 0, len(exportMap))
	for k, v := range exportMap {
		label := v[0]
		if isJa {
			label = v[1]
		}
		statusOpts = append(statusOpts, map[string]any{"id": k, "label": label})
	}

	// Build country list from location.custom.country_name.
	type ckey struct{ name, nameEn string }
	seen := map[ckey]bool{}
	var countries []map[string]any
	for _, l := range locations {
		custom := decodeCustom(l.Custom)
		name := customStr(custom, "country_name")
		nameEn := customStr(custom, "country_name_en")
		if name == "" && nameEn == "" {
			continue
		}
		k := ckey{name, nameEn}
		if !seen[k] {
			seen[k] = true
			countries = append(countries, map[string]any{
				"country_name":    name,
				"country_name_en": nameEn,
			})
		}
	}

	return map[string]any{
		"exhibitor_categories":     catItems,
		"exhibitor_english_status": statusOpts,
		"exhibitor_country":        orEmptySlice(countries),
	}
}

// ── Product search options ────────────────────────────────────────────────────

// prepareProductSearchOptions mirrors publish_venue_v2.py:196.
func prepareProductSearchOptions(
	venue *domain.Venue,
	lang string,
	locationCats []*domain.LocationCategory,
	productCats []*domain.ProductCategory,
) map[string]any {
	isJa := lang == "ja"
	exportMap := SelectExportableMapping(venue.ExternalID)
	importerMap := SelectImporterMapping(venue.ExternalID)

	// exhibitor_categories (from location categories)
	catItems := make([]map[string]any, 0, len(locationCats))
	for _, c := range locationCats {
		if c.Source != "external" {
			continue
		}
		catItems = append(catItems, map[string]any{
			"id":   c.ID.String(),
			"name": localizedCatName(c.Localization, c.Name, lang),
		})
	}

	// exhibitor_english_status
	statusOpts := make([]map[string]any, 0, len(exportMap))
	for k, v := range exportMap {
		label := v[0]
		if isJa {
			label = v[1]
		}
		statusOpts = append(statusOpts, map[string]any{"id": k, "label": label})
	}

	// product_importer
	importerOpts := make([]map[string]any, 0, len(importerMap))
	for k, v := range importerMap {
		lbl := v[0]
		if isJa {
			lbl = v[1]
		}
		importerOpts = append(importerOpts, map[string]any{"id": k, "label": lbl})
	}

	// product_requisite_documents
	reqDocs := make([]string, len(ProductRequisiteDocuments))
	for i, d := range ProductRequisiteDocuments {
		if isJa {
			reqDocs[i] = d[1]
		} else {
			reqDocs[i] = d[0]
		}
	}

	// product_categories (all) and grouped_product_categories (parent _0 entries)
	var allCats []map[string]any
	var groupedCats []map[string]any
	for _, c := range productCats {
		name := localizedProductCatName(c.Localization, c.Name, lang)
		item := map[string]any{"id": c.ID.String(), "name": name, "external_id": c.ExternalID}
		allCats = append(allCats, item)
		if strings.HasSuffix(c.ExternalID, "_0") {
			groupedCats = append(groupedCats, item)
		}
	}

	return map[string]any{
		"exhibitor_categories":        catItems,
		"exhibitor_english_status":    statusOpts,
		"product_importer":            importerOpts,
		"product_requisite_documents": reqDocs,
		"product_categories":          orEmptySlice(allCats),
		"grouped_product_categories":  orEmptySlice(groupedCats),
	}
}

// ── Helpers ───────────────────────────────────────────────────────────────────

func decodeCustom(raw json.RawMessage) map[string]any {
	return localization.DecodeBlob(raw)
}

func customStr(m map[string]any, key string) string {
	if m == nil {
		return ""
	}
	if v, ok := m[key].(string); ok {
		return v
	}
	return ""
}

func customInt(m map[string]any, key string) int {
	if m == nil {
		return 0
	}
	switch v := m[key].(type) {
	case float64:
		return int(v)
	case int:
		return v
	}
	return 0
}

func localizedName(raw json.RawMessage, fallback, lang string) string {
	return localization.Field(localization.DecodeBlob(raw), "common_name", lang, fallback)
}

func localizedProductName(raw json.RawMessage, fallback, lang string) string {
	return localization.Field(localization.DecodeBlob(raw), "name", lang, fallback)
}

func localizedCatName(raw json.RawMessage, fallback, lang string) string {
	return localization.Field(localization.DecodeBlob(raw), "name", lang, fallback)
}

func localizedProductCatName(raw json.RawMessage, fallback, lang string) string {
	return localizedProductName(raw, fallback, lang)
}

func stripNilKeys(m map[string]any) map[string]any {
	if m == nil {
		return nil
	}
	out := make(map[string]any, len(m))
	for k, v := range m {
		if v != nil {
			out[k] = v
		}
	}
	return out
}

func orEmpty[T any](s []T) any {
	if s == nil {
		return []T{}
	}
	return s
}

func orEmptySlice(s []map[string]any) []map[string]any {
	if s == nil {
		return []map[string]any{}
	}
	return s
}

// Ensure unused uuid import is used.
var _ = uuid.Nil
