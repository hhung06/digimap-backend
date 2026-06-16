package searchindex

import (
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/hhung06/digimap-backend/internal/domain"
	"github.com/hhung06/digimap-backend/internal/localization"
	"github.com/hhung06/digimap-backend/internal/platform/search"
)

// buildDocs constructs all BulkDocs for exhibitors (locations) and products.
func buildDocs(
	venue *domain.Venue,
	locations []*domain.Location,
	products []*domain.Product,
) []search.BulkDoc {
	// Index locations by ID for product lookups.
	locationsByID := make(map[uuid.UUID]*domain.Location, len(locations))
	for _, loc := range locations {
		locationsByID[loc.ID] = loc
	}

	// Group products by location ID for the products_summary field.
	productsByLocationID := make(map[uuid.UUID][]*domain.Product)
	for _, p := range products {
		if p.DeletedAt != nil {
			continue
		}
		if p.LocationID != nil {
			productsByLocationID[*p.LocationID] = append(productsByLocationID[*p.LocationID], p)
		}
	}

	var docs []search.BulkDoc

	// Exhibitor docs — only searchable, non-deleted locations.
	for _, loc := range locations {
		if !loc.IsSearchable || loc.DeletedAt != nil {
			continue
		}
		docs = append(docs, buildLocationDoc(venue, loc, productsByLocationID[loc.ID]))
	}

	// Product docs — all non-deleted products.
	for _, p := range products {
		if p.DeletedAt != nil {
			continue
		}
		var parentLoc *domain.Location
		if p.LocationID != nil {
			parentLoc = locationsByID[*p.LocationID]
		}
		docs = append(docs, buildProductDoc(venue, p, parentLoc))
	}

	return docs
}

func buildLocationDoc(venue *domain.Venue, loc *domain.Location, locProducts []*domain.Product) search.BulkDoc {
	locBlob := localization.DecodeBlob(loc.Localization)
	custom := localization.DecodeBlob(loc.Custom)

	nameJA := localization.Field(locBlob, "common_name", "ja", loc.CommonName)
	nameEN := localization.Field(locBlob, "common_name", "en", "")

	searchTextJA := joinNonEmpty(
		nameJA,
		customString(custom, "exhibitor_name"),
		customString(custom, "exhibition_zone"),
		customString(custom, "country_name"),
		customString(custom, "exhibitor_highlights"),
	)
	searchTextEN := joinNonEmpty(
		nameEN,
		customString(custom, "exhibitor_name_en"),
		customString(custom, "exhibition_zone_en"),
		customString(custom, "country_name_en"),
		customString(custom, "exhibitor_highlights_en"),
	)

	// Build products summary pipe-joined product names for this location.
	var productNames []string
	for _, p := range locProducts {
		if p.DeletedAt != nil {
			continue
		}
		pc := localization.DecodeBlob(p.Custom)
		if name := customString(pc, "name"); name != "" {
			productNames = append(productNames, name)
		}
	}
	productsSummary := strings.Join(productNames, " | ")

	images := map[string]any{}
	if loc.CommonLogo != "" {
		images["logo"] = loc.CommonLogo
	}
	if loc.CommonSmallLogo != "" {
		images["logo_small"] = loc.CommonSmallLogo
	}
	if loc.CommonMediumLogo != "" {
		images["logo_medium"] = loc.CommonMediumLogo
	}
	if loc.CommonLargeLogo != "" {
		images["logo_large"] = loc.CommonLargeLogo
	}
	if loc.TopLogo != "" {
		images["top_logo"] = loc.TopLogo
	}

	source := map[string]any{
		"type": "location",
		"venue": map[string]any{
			"id":          venue.ID,
			"external_id": venue.ExternalID,
			"section":     1,
		},
		"search_text": map[string]any{
			"ja": searchTextJA,
			"en": searchTextEN,
		},
		"exhibitor": map[string]any{
			"id":                   loc.ID,
			"name":                 firstNonEmpty(customString(custom, "exhibitor_name"), nameJA),
			"name_en":              firstNonEmpty(customString(custom, "exhibitor_name_en"), nameEN),
			"name_kana":            customString(custom, "exhibitor_name_kana"),
			"booth_number":         loc.BoothNumber,
			"exhibition_zone":      customString(custom, "exhibition_zone"),
			"exhibition_zone_en":   customString(custom, "exhibition_zone_en"),
			"country":              customString(custom, "country_name"),
			"country_en":           customString(custom, "country_name_en"),
			"description":          customString(custom, "exhibitor_highlights"),
			"description_en":       customString(custom, "exhibitor_highlights_en"),
			"products_summary":     productsSummary,
			"product_count":        len(locProducts),
			"source":               loc.Source,
			"common_location_type": loc.CommonLocationType,
			"common_sub_type":      loc.CommonLocationSubType,
		},
		"images":     images,
		"categories": []any{},
		"indexed_at": time.Now().UTC().Format(time.RFC3339),
	}

	return search.BulkDoc{
		ID:     loc.ID.String(),
		Source: source,
	}
}

func buildProductDoc(venue *domain.Venue, p *domain.Product, parentLoc *domain.Location) search.BulkDoc {
	custom := localization.DecodeBlob(p.Custom)

	var parentCustom map[string]any
	if parentLoc != nil {
		parentCustom = localization.DecodeBlob(parentLoc.Custom)
	}

	searchTextJA := joinNonEmpty(
		customString(custom, "name"),
		customString(custom, "name_kana"),
		customString(custom, "keyword"),
		customString(custom, "specialities"),
		customString(custom, "ingredient"),
		customString(custom, "target"),
		customString(custom, "use"),
		customString(parentCustom, "exhibitor_name"),
		customString(parentCustom, "exhibition_zone"),
		customString(parentCustom, "country_name"),
		customString(custom, "parent_category_ja"),
	)
	searchTextEN := joinNonEmpty(
		customString(custom, "name_en"),
		customString(parentCustom, "exhibitor_name_en"),
		customString(parentCustom, "exhibition_zone_en"),
		customString(parentCustom, "country_name_en"),
		customString(custom, "parent_category_en"),
	)

	exhibitorDoc := map[string]any{}
	if parentLoc != nil {
		exhibitorDoc["id"] = parentLoc.ID
		exhibitorDoc["name"] = customString(parentCustom, "exhibitor_name")
		exhibitorDoc["name_en"] = customString(parentCustom, "exhibitor_name_en")
		exhibitorDoc["booth_number"] = parentLoc.BoothNumber
		exhibitorDoc["exhibition_zone"] = customString(parentCustom, "exhibition_zone")
		exhibitorDoc["exhibition_zone_en"] = customString(parentCustom, "exhibition_zone_en")
		exhibitorDoc["country"] = customString(parentCustom, "country_name")
		exhibitorDoc["country_en"] = customString(parentCustom, "country_name_en")
	}

	source := map[string]any{
		"type": "product",
		"venue": map[string]any{
			"id":          venue.ID,
			"external_id": venue.ExternalID,
			"section":     1,
		},
		"search_text": map[string]any{
			"ja": searchTextJA,
			"en": searchTextEN,
		},
		"product": map[string]any{
			"id":          p.ID,
			"name":        customString(custom, "name"),
			"name_en":     customString(custom, "name_en"),
			"name_kana":   customString(custom, "name_kana"),
			"category":    customString(custom, "parent_category_ja"),
			"category_en": customString(custom, "parent_category_en"),
			"keywords":    customString(custom, "keyword"),
			"description": customString(custom, "specialities"),
			"price":       customString(custom, "price"),
			"ingredients": customString(custom, "ingredient"),
			"target":      customString(custom, "target"),
			"use":         customString(custom, "use"),
		},
		"exhibitor":  exhibitorDoc,
		"images":     map[string]any{},
		"categories": []any{},
		"indexed_at": time.Now().UTC().Format(time.RFC3339),
	}

	return search.BulkDoc{
		ID:     p.ID.String(),
		Source: source,
	}
}

// ── helpers ───────────────────────────────────────────────────────────────────

func customString(m map[string]any, key string) string {
	if m == nil {
		return ""
	}
	if v, ok := m[key]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

func joinNonEmpty(parts ...string) string {
	var filtered []string
	for _, p := range parts {
		if p != "" {
			filtered = append(filtered, p)
		}
	}
	return strings.Join(filtered, " ")
}

func firstNonEmpty(vals ...string) string {
	for _, v := range vals {
		if v != "" {
			return v
		}
	}
	return ""
}
