package service

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/hhung06/digimap-backend/internal/domain"
)

func (s *notificationService) ListPushTypes(ctx context.Context, venueID uuid.UUID) ([]map[string]any, error) {
	if s.venueRepo == nil {
		return nil, &domain.AppError{Err: domain.ErrInternal, Message: "venue repository not configured"}
	}
	venue, err := s.venueRepo.FindByID(ctx, venueID)
	if err != nil {
		return nil, err
	}

	prefix := expoPrefix(venue.ExternalID)
	pushTypes := cloneMapSlice(domain.PUSH_TYPES)
	if prefix == "" {
		return pushTypes, nil
	}

	optionSets := pushTypeOptionSets(prefix)
	seminars := fetchSeminarOptions(ctx, prefix)
	for _, pushType := range pushTypes {
		key, _ := pushType["key"].(string)
		if key == "seminar_registered" {
			pushType["options"] = seminars
			continue
		}
		if options, ok := optionSets[key]; ok {
			pushType["options"] = cloneMapSlice(options)
		}
	}
	return pushTypes, nil
}

func expoPrefix(externalID string) string {
	value := strings.ToLower(externalID)
	switch {
	case strings.Contains(value, "hcj"):
		return "hcj"
	case strings.Contains(value, "foodex"):
		return "foodex"
	default:
		return ""
	}
}

func pushTypeOptionSets(prefix string) map[string][]map[string]any {
	if prefix == "hcj" {
		return map[string][]map[string]any{
			"topic_visitor_industry":                domain.TOPIC_VISITOR_INDUSTRY_OPTIONS_HCJ,
			"topic_visitor_job":                     domain.TOPIC_VISITOR_JOB_OPTIONS_HCJ,
			"topic_visitor_position":                domain.TOPIC_VISITOR_POSITION_OPTIONS_HCJ,
			"topic_visitor_interest":                domain.TOPIC_VISITOR_INTERESTED_CATEGORIES_HCJ,
			"topic_exhibitor_zone":                  domain.TOPIC_VISITOR_ZONE_OPTIONS_HCJ,
			"topic_visitor_survey_segment_purchase": domain.TOPIC_SURVEY_PURCHASE_OPTIONS_HCJ,
			"topic_visitor_survey_segment_purpose":  domain.TOPIC_SURVEY_VISIT_PURPOSE_HCJ,
			"topic_visitor_survey_segment_budget":   domain.TOPIC_SURVEY_BUDGET_OPTIONS_HCJ,
		}
	}
	return map[string][]map[string]any{
		"topic_visitor_industry":                domain.TOPIC_VISITOR_INDUSTRY_OPTIONS_FOODEX,
		"topic_visitor_job":                     domain.TOPIC_VISITOR_JOB_OPTIONS_FOODEX,
		"topic_visitor_position":                domain.TOPIC_VISITOR_POSITION_OPTIONS_FOODEX,
		"topic_visitor_interest":                domain.TOPIC_VISITOR_INTERESTED_CATEGORIES_FOODEX,
		"topic_exhibitor_zone":                  domain.TOPIC_VISITOR_ZONE_OPTIONS_FOODEX,
		"topic_visitor_survey_segment_purchase": domain.TOPIC_SURVEY_PURCHASE_OPTIONS_FOODEX,
		"topic_visitor_survey_segment_purpose":  domain.TOPIC_SURVEY_VISIT_PURPOSE_FOODEX,
		"topic_visitor_survey_segment_budget":   domain.TOPIC_SURVEY_BUDGET_OPTIONS_FOODEX,
	}
}

func fetchSeminarOptions(ctx context.Context, prefix string) []map[string]any {
	baseURL, path := seminarEndpoint(prefix)
	if baseURL == "" || path == "" {
		return []map[string]any{}
	}
	apiKey := seminarAPIKey(prefix)
	authToken := os.Getenv("BFACE_BASIC_AUTHENTICATION")
	if apiKey == "" || authToken == "" {
		return []map[string]any{}
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, baseURL+path, nil)
	if err != nil {
		slog.Warn("push types: failed to build seminar request", "prefix", prefix, "error", err)
		return []map[string]any{}
	}
	req.Header.Set("Authorization", authToken)
	req.Header.Set("X-API-KEY", apiKey)

	client := &http.Client{Timeout: 10 * time.Second}
	res, err := client.Do(req)
	if err != nil {
		slog.Warn("push types: failed to fetch seminars", "prefix", prefix, "error", err)
		return []map[string]any{}
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		slog.Warn("push types: seminar API returned non-200", "prefix", prefix, "status", res.StatusCode)
		return []map[string]any{}
	}

	var body struct {
		SeminarList []map[string]any `json:"seminar_list"`
	}
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		slog.Warn("push types: failed to decode seminar response", "prefix", prefix, "error", err)
		return []map[string]any{}
	}

	options := make([]map[string]any, 0, len(body.SeminarList))
	for _, seminar := range body.SeminarList {
		id, hasID := seminar["seminar_id"]
		title, hasTitle := seminar["title"]
		if !hasID || !hasTitle {
			continue
		}
		options = append(options, map[string]any{
			"seminar_id":     id,
			"title":          title,
			"start_datetime": seminar["start_datetime"],
			"end_datetime":   seminar["end_datetime"],
			"value":          id,
			"label":          title,
			"label_en":       title,
		})
	}
	return options
}

func seminarEndpoint(prefix string) (string, string) {
	host := "https://test.jma-tradeshow.com"
	switch strings.ToLower(os.Getenv("APP_ENV")) {
	case "prod", "production":
		host = "https://www.jma-tradeshow.com"
	}
	if prefix == "hcj" {
		return host, "/hcj/api/buyer/seminar.php"
	}
	if prefix == "foodex" {
		return host, "/foodex/api/buyer/seminar.php"
	}
	return "", ""
}

func seminarAPIKey(prefix string) string {
	if prefix == "hcj" {
		return os.Getenv("HCJ_X_API_KEY")
	}
	if prefix == "foodex" {
		return os.Getenv("FOODEX_X_API_KEY")
	}
	return ""
}

func cloneMapSlice(items []map[string]any) []map[string]any {
	cloned := make([]map[string]any, len(items))
	for i, item := range items {
		cloned[i] = cloneMap(item)
	}
	return cloned
}

func cloneMap(item map[string]any) map[string]any {
	cloned := make(map[string]any, len(item))
	for key, value := range item {
		switch typed := value.(type) {
		case []map[string]any:
			cloned[key] = cloneMapSlice(typed)
		case map[string]any:
			cloned[key] = cloneMap(typed)
		default:
			cloned[key] = value
		}
	}
	return cloned
}
