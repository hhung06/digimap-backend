package service

import (
	"context"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/google/uuid"

	"github.com/hhung06/digimap-backend/internal/domain"
	"github.com/hhung06/digimap-backend/internal/localization"
	"github.com/hhung06/digimap-backend/internal/repository"
	"github.com/hhung06/digimap-backend/internal/service/bundle"
)

type SearchOptionsService interface {
	Get(ctx context.Context, venueID uuid.UUID, origin, lang string) (map[string]any, error)
}

type searchOptionsService struct {
	venueRepo            repository.VenueRepository
	locationRepo         repository.LocationRepository
	locationCategoryRepo repository.LocationCategoryRepository
	productRepo          repository.ProductRepository
}

func NewSearchOptionsService(
	venueRepo repository.VenueRepository,
	locationRepo repository.LocationRepository,
	locationCategoryRepo repository.LocationCategoryRepository,
	productRepo repository.ProductRepository,
) SearchOptionsService {
	return &searchOptionsService{
		venueRepo:            venueRepo,
		locationRepo:         locationRepo,
		locationCategoryRepo: locationCategoryRepo,
		productRepo:          productRepo,
	}
}

func (s *searchOptionsService) Get(ctx context.Context, venueID uuid.UUID, origin, lang string) (map[string]any, error) {
	if lang == "" {
		lang = "en"
	}
	origin = strings.TrimSpace(origin)
	if origin == "" {
		origin = "exhibitors"
	}

	venue, err := s.venueRepo.FindByID(ctx, venueID)
	if err != nil {
		return nil, err
	}
	locationCats, err := s.locationCategoryRepo.List(ctx, venueID)
	if err != nil {
		return nil, err
	}

	switch origin {
	case "exhibitors", "locations":
		locations, _, err := s.locationRepo.List(ctx, venueID, nil, nil, domain.Pagination{Page: 1, PageSize: 100000})
		if err != nil {
			return nil, err
		}
		return buildExhibitorSearchOptions(venue, lang, locations, locationCats), nil
	case "products":
		productCats, err := s.productRepo.ListCategories(ctx, venueID)
		if err != nil {
			return nil, err
		}
		return buildProductSearchOptions(venue, lang, locationCats, productCats), nil
	default:
		return nil, domain.NewValidation(map[string]string{"origin": "must be products, exhibitors, or locations"})
	}
}

func buildExhibitorSearchOptions(venue *domain.Venue, lang string, locations []*domain.Location, cats []*domain.LocationCategory) map[string]any {
	return map[string]any{
		"exhibitor_categories":     locationCategoryOptions(cats, lang),
		"exhibitor_english_status": mappingOptions(bundle.SelectExportableMapping(venue.ExternalID), lang),
		"exhibitor_country":        exhibitorCountries(locations, lang),
	}
}

func buildProductSearchOptions(venue *domain.Venue, lang string, locationCats []*domain.LocationCategory, productCats []*domain.ProductCategory) map[string]any {
	return map[string]any{
		"exhibitor_categories":        locationCategoryOptions(locationCats, lang),
		"exhibitor_english_status":    mappingOptions(bundle.SelectExportableMapping(venue.ExternalID), lang),
		"product_importer":            mappingOptions(bundle.SelectImporterMapping(venue.ExternalID), lang),
		"product_requisite_documents": requisiteDocumentOptions(lang),
		"product_categories":          productCategoryChildOptions(productCats, lang),
		"grouped_product_categories":  groupedProductCategoryOptions(productCats, lang),
	}
}

func locationCategoryOptions(cats []*domain.LocationCategory, lang string) []map[string]any {
	out := make([]map[string]any, 0, len(cats))
	for _, c := range cats {
		if c == nil || c.Source != "external" {
			continue
		}
		name := localization.Field(localization.DecodeBlob(c.Localization), "name", lang, c.Name)
		shortName := localization.Field(localization.DecodeBlob(c.Localization), "shortName", lang, c.ShortName)
		out = append(out, map[string]any{
			"id":        c.ID.String(),
			"name":      name,
			"type":      c.Type,
			"shortName": shortName,
		})
	}
	return out
}

func mappingOptions[T ~[2]string](mapping map[int]T, lang string) []map[string]string {
	keys := make([]int, 0, len(mapping))
	for key := range mapping {
		keys = append(keys, key)
	}
	sort.Ints(keys)
	out := make([]map[string]string, 0, len(keys))
	for _, key := range keys {
		entry := mapping[key]
		label := entry[0]
		if lang == "ja" {
			label = entry[1]
		}
		out = append(out, map[string]string{strconv.Itoa(key): label})
	}
	return out
}

func exhibitorCountries(locations []*domain.Location, lang string) []string {
	field := "country_name_en"
	if lang == "ja" {
		field = "country_name"
	}
	seen := map[string]bool{}
	for _, l := range locations {
		if l == nil {
			continue
		}
		value, _ := localization.DecodeBlob(l.Custom)[field].(string)
		value = strings.TrimSpace(value)
		if value != "" {
			seen[value] = true
		}
	}
	out := make([]string, 0, len(seen))
	for value := range seen {
		out = append(out, value)
	}
	sort.Slice(out, func(i, j int) bool { return strings.ToLower(out[i]) < strings.ToLower(out[j]) })
	return out
}

func requisiteDocumentOptions(lang string) []string {
	out := make([]string, len(bundle.ProductRequisiteDocuments))
	for i, doc := range bundle.ProductRequisiteDocuments {
		out[i] = doc[0]
		if lang == "ja" {
			out[i] = doc[1]
		}
	}
	return out
}

func productCategoryChildOptions(cats []*domain.ProductCategory, lang string) []map[string]any {
	children := make([]*domain.ProductCategory, 0, len(cats))
	for _, c := range cats {
		if c == nil || c.Source != "external" || isParentProductCategory(c) {
			continue
		}
		children = append(children, c)
	}
	sortProductCategories(children)
	out := make([]map[string]any, 0, len(children))
	for _, c := range children {
		out = append(out, productCategoryOption(c, lang, nil))
	}
	return out
}

func groupedProductCategoryOptions(cats []*domain.ProductCategory, lang string) []map[string]any {
	childrenByParent := map[string][]*domain.ProductCategory{}
	parents := make([]*domain.ProductCategory, 0, len(cats))
	for _, c := range cats {
		if c == nil || c.Source != "external" || c.ExternalID == nil {
			continue
		}
		if isParentProductCategory(c) {
			parents = append(parents, c)
			continue
		}
		childrenByParent[productCategoryParentKey(*c.ExternalID)] = append(childrenByParent[productCategoryParentKey(*c.ExternalID)], c)
	}
	sortProductCategories(parents)
	out := make([]map[string]any, 0, len(parents))
	for _, parent := range parents {
		children := childrenByParent[productCategoryParentKey(*parent.ExternalID)]
		sortProductCategories(children)
		childItems := make([]map[string]any, 0, len(children))
		for _, child := range children {
			childItems = append(childItems, productCategoryOption(child, lang, nil))
		}
		out = append(out, productCategoryOption(parent, lang, childItems))
	}
	return out
}

func productCategoryOption(c *domain.ProductCategory, lang string, children []map[string]any) map[string]any {
	externalID := ""
	if c.ExternalID != nil {
		externalID = *c.ExternalID
	}
	item := map[string]any{
		"id":         c.ID.String(),
		"externalId": externalID,
		"name":       localization.Field(localization.DecodeBlob(c.Localization), "name", lang, c.Name),
	}
	if children != nil {
		item["children"] = children
	}
	return item
}

func isParentProductCategory(c *domain.ProductCategory) bool {
	return c.ExternalID != nil && strings.HasSuffix(*c.ExternalID, "_0")
}

func productCategoryParentKey(externalID string) string {
	idx := strings.Index(externalID, "_")
	if idx < 0 {
		return externalID
	}
	return externalID[:idx]
}

func sortProductCategories(cats []*domain.ProductCategory) {
	sort.Slice(cats, func(i, j int) bool {
		return productCategorySortParts(cats[i]) < productCategorySortParts(cats[j])
	})
}

func productCategorySortParts(c *domain.ProductCategory) string {
	if c == nil || c.ExternalID == nil {
		return ""
	}
	parts := strings.SplitN(*c.ExternalID, "_", 2)
	parent, _ := strconv.Atoi(parts[0])
	child := 0
	if len(parts) > 1 {
		child, _ = strconv.Atoi(parts[1])
	}
	return fmt.Sprintf("%08d_%08d", parent, child)
}
