package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/hhung06/digimap-backend/internal/domain"
	"github.com/hhung06/digimap-backend/internal/dto"
	"github.com/hhung06/digimap-backend/internal/repository"
	"github.com/hhung06/digimap-backend/internal/service"
)

type webhookHandler struct {
	venueRepo            repository.VenueRepository
	locationRepo         repository.LocationRepository
	productRepo          repository.ProductRepository
	locationCategoryRepo repository.LocationCategoryRepository
	notificationSvc      service.NotificationService
}

func newWebhookHandler(
	venueRepo repository.VenueRepository,
	locationRepo repository.LocationRepository,
	productRepo repository.ProductRepository,
	locationCategoryRepo repository.LocationCategoryRepository,
	notificationSvc service.NotificationService,
) *webhookHandler {
	return &webhookHandler{
		venueRepo:            venueRepo,
		locationRepo:         locationRepo,
		productRepo:          productRepo,
		locationCategoryRepo: locationCategoryRepo,
		notificationSvc:      notificationSvc,
	}
}

// validateJMAHeaders validates TOKEN, YEAR, and SYSTEM-CODE headers.
// Returns (systemCode, year, token, true) on success, or writes an error response and returns ("","","",false).
func (h *webhookHandler) validateJMAHeaders(c *gin.Context) (systemCode, year, token string, ok bool) {
	token = c.GetHeader("TOKEN")
	year = c.GetHeader("YEAR")
	systemCode = strings.ToUpper(c.GetHeader("SYSTEM-CODE"))

	if token == "" {
		c.JSON(http.StatusUnauthorized, dto.Fail(dto.CodeAuthRequired, "TOKEN header is required"))
		return "", "", "", false
	}
	if year == "" {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "YEAR header is required"))
		return "", "", "", false
	}
	yearInt, err := strconv.Atoi(year)
	if err != nil || yearInt < 2020 || yearInt > 2050 {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "YEAR must be a valid 4-digit year (2020–2050)"))
		return "", "", "", false
	}
	if systemCode != "FOODEX" && systemCode != "HCJ" {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "SYSTEM-CODE must be FOODEX or HCJ"))
		return "", "", "", false
	}
	return systemCode, year, token, true
}

// lookupVenueByToken fetches the venue by private key (TOKEN header).
// On failure it writes the error response and returns nil.
func (h *webhookHandler) lookupVenueByToken(c *gin.Context, token string) *domain.Venue {
	venue, err := h.venueRepo.FindByPrivateKey(c.Request.Context(), token)
	if err != nil {
		var appErr *domain.AppError
		if errors.As(err, &appErr) && errors.Is(appErr.Err, domain.ErrNotFound) {
			c.JSON(http.StatusUnauthorized, dto.Fail(dto.CodeAuthRequired, "invalid TOKEN"))
			return nil
		}
		c.JSON(http.StatusInternalServerError, dto.Fail(dto.CodeInternalError, "internal error"))
		return nil
	}
	return venue
}

// isActive converts the raw status field (bool or int) to a boolean.
func isActive(status any) bool {
	switch v := status.(type) {
	case bool:
		return v
	case float64: // JSON numbers decode as float64
		return v != 0
	case int:
		return v != 0
	}
	return false
}

// ── Exhibitor webhook ─────────────────────────────────────────────────────────

type jmaExhibitorItem struct {
	ExhibitorID string `json:"exhibitor_id" binding:"required"`
	Status      any    `json:"status"` // bool or int (0/1)
}

type jmaExhibitorPayload struct {
	ExhibitorList []jmaExhibitorItem `json:"exhibitor_list" binding:"required,min=1"`
}

type webhookResults struct {
	Processed int      `json:"processed"`
	Created   []string `json:"created"`
	Updated   []string `json:"updated"`
	Deleted   []string `json:"deleted"`
	Errors    []any    `json:"errors"`
}

// JMAExhibitorUpdate handles POST /api/v1/webhooks/jma/exhibitors
func (h *webhookHandler) JMAExhibitorUpdate(c *gin.Context) {
	systemCode, _, token, ok := h.validateJMAHeaders(c)
	if !ok {
		return
	}

	venue := h.lookupVenueByToken(c, token)
	if venue == nil {
		return
	}

	var payload jmaExhibitorPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "exhibitor_list is required and must be non-empty"))
		return
	}

	ctx := c.Request.Context()
	results := webhookResults{
		Created: []string{},
		Updated: []string{},
		Deleted: []string{},
		Errors:  []any{},
	}

	_ = systemCode // available for future per-system field mapping

	for i, item := range payload.ExhibitorList {
		if item.ExhibitorID == "" {
			results.Errors = append(results.Errors, map[string]any{
				"index":        i,
				"exhibitor_id": item.ExhibitorID,
				"error":        "exhibitor_id is required",
			})
			continue
		}

		if !isActive(item.Status) {
			// Delete
			existing, err := h.locationRepo.FindByExternalID(ctx, venue.ID, item.ExhibitorID, domain.LocationTypeBooth)
			if err == nil {
				if delErr := h.locationRepo.Delete(ctx, existing.ID); delErr != nil {
					results.Errors = append(results.Errors, map[string]any{
						"index":        i,
						"exhibitor_id": item.ExhibitorID,
						"error":        delErr.Error(),
					})
					continue
				}
				results.Deleted = append(results.Deleted, item.ExhibitorID)
				results.Processed++
			}
			// Not found on delete is not an error (idempotent)
			continue
		}

		// Build location fields from raw payload (pass through custom data)
		rawJSON, _ := json.Marshal(item)
		rawPayload := map[string]any{}
		_ = json.Unmarshal(rawJSON, &rawPayload)

		name := getString(rawPayload, "exhibitor_name_en", "exhibitor_name")
		desc := getString(rawPayload, "exhibitor_highlights_en")
		website := getString(rawPayload, "exhibitor_webguide_url")
		boothNum := getString(rawPayload, "booth_number")
		customBytes, _ := json.Marshal(rawPayload)
		localizationBytes := buildExhibitorLocalization(rawPayload)

		existing, err := h.locationRepo.FindByExternalID(ctx, venue.ID, item.ExhibitorID, domain.LocationTypeBooth)
		if err != nil {
			// Create
			loc := &domain.Location{
				VenueID:             venue.ID,
				CommonName:          name,
				CommonDescription:   desc,
				CommonSocialWebsite: website,
				BoothNumber:         boothNum,
				CommonLocationType:  domain.LocationTypeBooth,
				ExternalID:          item.ExhibitorID,
				Source:              "external",
				Custom:              json.RawMessage(customBytes),
				Localization:        json.RawMessage(localizationBytes),
			}
			if createErr := h.locationRepo.Create(ctx, loc); createErr != nil {
				results.Errors = append(results.Errors, map[string]any{
					"index":        i,
					"exhibitor_id": item.ExhibitorID,
					"error":        createErr.Error(),
				})
				continue
			}
			// Assign category after create
			catID := h.upsertExhibitorCategory(ctx, venue.ID, rawPayload)
			if catID != nil {
				loc.MainCategoryID = catID
				_ = h.locationRepo.Update(ctx, loc)
				_ = h.locationRepo.SetCategories(ctx, loc.ID, []uuid.UUID{*catID})
			}
			results.Created = append(results.Created, item.ExhibitorID)
		} else {
			// Update
			existing.CommonName = name
			existing.CommonDescription = desc
			existing.CommonSocialWebsite = website
			existing.BoothNumber = boothNum
			existing.Custom = json.RawMessage(customBytes)
			existing.Localization = json.RawMessage(localizationBytes)

			catID := h.upsertExhibitorCategory(ctx, venue.ID, rawPayload)
			if catID != nil {
				existing.MainCategoryID = catID
			}
			if updateErr := h.locationRepo.Update(ctx, existing); updateErr != nil {
				results.Errors = append(results.Errors, map[string]any{
					"index":        i,
					"exhibitor_id": item.ExhibitorID,
					"error":        updateErr.Error(),
				})
				continue
			}
			if catID != nil {
				_ = h.locationRepo.SetCategories(ctx, existing.ID, []uuid.UUID{*catID})
			}
			results.Updated = append(results.Updated, item.ExhibitorID)
		}
		results.Processed++
	}

	c.JSON(http.StatusOK, dto.OK(results))
}

// upsertExhibitorCategory finds or creates a LocationCategory for the exhibition zone
// from the raw exhibitor payload. Returns nil if no zone name is present.
func (h *webhookHandler) upsertExhibitorCategory(ctx context.Context, venueID uuid.UUID, raw map[string]any) *uuid.UUID {
	zoneEn := getString(raw, "exhibition_zone_en")
	zoneJa := getString(raw, "exhibition_zone")
	if zoneEn == "" {
		return nil
	}

	locBytes, _ := json.Marshal(map[string]any{"name_en": zoneEn, "name_ja": zoneJa})

	cat, err := h.locationCategoryRepo.FindByNameAndVenue(ctx, venueID, zoneEn, "external")
	if err != nil {
		// Create
		cat = &domain.LocationCategory{
			VenueID:      venueID,
			Name:         zoneEn,
			Color:        "#0000FF",
			Visible:      true,
			Source:       "external",
			Localization: json.RawMessage(locBytes),
		}
		if createErr := h.locationCategoryRepo.Create(ctx, cat); createErr != nil {
			return nil
		}
	} else {
		cat.Localization = json.RawMessage(locBytes)
		_ = h.locationCategoryRepo.Update(ctx, cat)
	}
	return &cat.ID
}

// buildExhibitorLocalization builds the Localization JSON for a location from raw exhibitor payload.
func buildExhibitorLocalization(raw map[string]any) []byte {
	nameEn := getString(raw, "exhibitor_name_en")
	nameJa := getString(raw, "exhibitor_name")
	descEn := getString(raw, "exhibitor_highlights_en")
	descJa := getString(raw, "exhibitor_highlights")
	b, _ := json.Marshal(map[string]any{
		"common_name_en":        nameEn,
		"common_name_ja":        nameJa,
		"common_short_name_en":  nameEn,
		"common_short_name_ja":  nameJa,
		"common_description_en": descEn,
		"common_description_ja": descJa,
	})
	return b
}

// ── Product webhook ───────────────────────────────────────────────────────────

type jmaProductItem struct {
	ProductID   string `json:"product_id" binding:"required"`
	ExhibitorID string `json:"exhibitor_id" binding:"required"`
	Status      any    `json:"status"`  // bool or int (0/1)
	Section     *int   `json:"section"` // optional; default 1
}

type jmaProductPayload struct {
	ProductList []jmaProductItem `json:"product_list" binding:"required,min=1"`
}

// JMAProductUpdate handles POST /api/v1/webhooks/jma/products
func (h *webhookHandler) JMAProductUpdate(c *gin.Context) {
	systemCode, _, token, ok := h.validateJMAHeaders(c)
	if !ok {
		return
	}

	venue := h.lookupVenueByToken(c, token)
	if venue == nil {
		return
	}

	var payload jmaProductPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "product_list is required and must be non-empty"))
		return
	}

	ctx := c.Request.Context()
	results := webhookResults{
		Created: []string{},
		Updated: []string{},
		Deleted: []string{},
		Errors:  []any{},
	}

	_ = systemCode // available for future per-system field mapping

	for i, item := range payload.ProductList {
		if item.ProductID == "" || item.ExhibitorID == "" {
			results.Errors = append(results.Errors, map[string]any{
				"index":        i,
				"product_id":   item.ProductID,
				"exhibitor_id": item.ExhibitorID,
				"error":        "product_id and exhibitor_id are required",
			})
			continue
		}

		section := 1
		if item.Section != nil {
			section = *item.Section
		}
		productCode := fmt.Sprintf("%s_%d_%s", item.ExhibitorID, section, item.ProductID)

		if !isActive(item.Status) {
			// Delete
			existing, err := h.productRepo.FindByCode(ctx, venue.ID, productCode, "external")
			if err == nil {
				if delErr := h.productRepo.Delete(ctx, existing.ID); delErr != nil {
					results.Errors = append(results.Errors, map[string]any{
						"index":        i,
						"product_id":   item.ProductID,
						"exhibitor_id": item.ExhibitorID,
						"error":        delErr.Error(),
					})
					continue
				}
				results.Deleted = append(results.Deleted, productCode)
				results.Processed++
			}
			continue
		}

		// Find the associated exhibitor location
		location, err := h.locationRepo.FindByExternalID(ctx, venue.ID, item.ExhibitorID, domain.LocationTypeBooth)
		if err != nil {
			results.Errors = append(results.Errors, map[string]any{
				"index":        i,
				"product_id":   item.ProductID,
				"exhibitor_id": item.ExhibitorID,
				"error":        fmt.Sprintf("exhibitor location not found: %s", item.ExhibitorID),
			})
			continue
		}

		rawJSON, _ := json.Marshal(item)
		rawPayload := map[string]any{}
		_ = json.Unmarshal(rawJSON, &rawPayload)

		name := getString(rawPayload, "name")
		size := getString(rawPayload, "size")
		price := getString(rawPayload, "price")
		country := getString(rawPayload, "product_country_of_origin")
		expiration := getString(rawPayload, "expiration")
		desc := getString(rawPayload, "specialities")
		customBytes, _ := json.Marshal(rawPayload)
		localizationBytes := buildProductLocalization(name, size, price, country, expiration, desc)

		existing, findErr := h.productRepo.FindByCode(ctx, venue.ID, productCode, "external")
		if findErr != nil {
			// Create
			loc := location
			prod := &domain.Product{
				VenueID:      venue.ID,
				LocationID:   &loc.ID,
				ExternalID:   productCode,
				Name:         name,
				Size:         size,
				Price:        price,
				Country:      country,
				Expiration:   expiration,
				Description:  desc,
				Source:       "external",
				Custom:       json.RawMessage(customBytes),
				Localization: json.RawMessage(localizationBytes),
			}
			if createErr := h.productRepo.Create(ctx, prod); createErr != nil {
				results.Errors = append(results.Errors, map[string]any{
					"index":        i,
					"product_id":   item.ProductID,
					"exhibitor_id": item.ExhibitorID,
					"error":        createErr.Error(),
				})
				continue
			}
			results.Created = append(results.Created, productCode)
		} else {
			// Update
			existing.Name = name
			existing.Size = size
			existing.Price = price
			existing.Country = country
			existing.Expiration = expiration
			existing.Description = desc
			existing.Custom = json.RawMessage(customBytes)
			existing.Localization = json.RawMessage(localizationBytes)
			if updateErr := h.productRepo.Update(ctx, existing); updateErr != nil {
				results.Errors = append(results.Errors, map[string]any{
					"index":        i,
					"product_id":   item.ProductID,
					"exhibitor_id": item.ExhibitorID,
					"error":        updateErr.Error(),
				})
				continue
			}
			results.Updated = append(results.Updated, productCode)
		}
		results.Processed++
	}

	c.JSON(http.StatusOK, dto.OK(results))
}

// buildProductLocalization builds the Localization JSON for a product.
func buildProductLocalization(name, size, price, country, expiration, desc string) []byte {
	b, _ := json.Marshal(map[string]any{
		"name_en":        name,
		"name_ja":        name,
		"size_en":        size,
		"size_ja":        size,
		"price_en":       price,
		"price_ja":       price,
		"country_en":     country,
		"country_ja":     country,
		"expiration_en":  expiration,
		"expiration_ja":  expiration,
		"description_en": desc,
		"description_ja": desc,
	})
	return b
}

// ── Push notification webhook ─────────────────────────────────────────────────

type jmaPushPayload struct {
	ExpoID       string         `json:"expo_id" binding:"required"`
	DeviceTokens []string       `json:"device_tokens" binding:"required,min=1"`
	Title        string         `json:"title" binding:"required"`
	Content      string         `json:"content" binding:"required"`
	Data         map[string]any `json:"data"`
	ScheduledAt  string         `json:"scheduled_at"`
}

// JMAPushNotification handles POST /api/v1/webhooks/jma/push
// Auth: TOKEN header only (no YEAR or SYSTEM-CODE required).
func (h *webhookHandler) JMAPushNotification(c *gin.Context) {
	token := c.GetHeader("TOKEN")
	if token == "" {
		c.JSON(http.StatusUnauthorized, dto.Fail(dto.CodeAuthRequired, "TOKEN header is required"))
		return
	}

	venue := h.lookupVenueByToken(c, token)
	if venue == nil {
		return
	}

	var payload jmaPushPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid payload: expo_id, device_tokens, title, and content are required"))
		return
	}

	expoID := strings.ToLower(payload.ExpoID)
	if expoID != "foodex" && expoID != "hcj" {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "expo_id must be foodex or hcj"))
		return
	}

	ctx := c.Request.Context()
	dataJSON, _ := json.Marshal(payload.Data)
	segFiltersJSON, _ := json.Marshal([]map[string]any{{"key": "jma_webhook_device_tokens", "type": nil}})

	chunks := chunkStrings(payload.DeviceTokens, 2000)
	for _, chunk := range chunks {
		tokensJSON, _ := json.Marshal(chunk)
		notif := &domain.Notification{
			VenueID:        &venue.ID,
			Title:          payload.Title,
			Content:        payload.Content,
			TargetApp:      expoID,
			SendType:       domain.NotifTypeImmediate,
			Status:         domain.NotifStatusUnsent,
			Kind:           domain.NotifKindNormal,
			SendStatus:     domain.NotifSendPending,
			Data:           json.RawMessage(dataJSON),
			SegmentFilters: json.RawMessage(segFiltersJSON),
			DeviceTokens:   json.RawMessage(tokensJSON),
		}
		if err := h.notificationSvc.Create(ctx, notif); err != nil {
			c.JSON(http.StatusInternalServerError, dto.Fail(dto.CodeInternalError, "failed to queue notification"))
			return
		}
	}

	c.JSON(http.StatusOK, dto.OK(gin.H{
		"success": true,
		"message": fmt.Sprintf("notification queued for %d tokens", len(payload.DeviceTokens)),
		"expo_id": expoID,
	}))
}

// ── helpers ───────────────────────────────────────────────────────────────────

// getString returns the first non-empty string from the given keys in the map.
func getString(m map[string]any, keys ...string) string {
	for _, k := range keys {
		if v, ok := m[k]; ok {
			if s, ok := v.(string); ok && s != "" {
				return s
			}
		}
	}
	return ""
}

// chunkStrings splits a slice into chunks of at most size n.
func chunkStrings(s []string, n int) [][]string {
	var chunks [][]string
	for i := 0; i < len(s); i += n {
		end := i + n
		if end > len(s) {
			end = len(s)
		}
		chunks = append(chunks, s[i:end])
	}
	return chunks
}
