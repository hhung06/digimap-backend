package service

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/hhung06/digimap-backend/internal/domain"
	"github.com/hhung06/digimap-backend/internal/repository"
)

type VisitorSurveySubmission struct {
	PublicKey      string
	FullName       string
	Email          string
	PhoneNumber    string
	VisitorType    int
	BusinessName   string
	Interests      []int
	OtherInterests string
	IsConsented    bool
	IPAddress      string
	UserAgent      string
}

type VisitorSurveySubmissionService interface {
	Submit(ctx context.Context, req VisitorSurveySubmission) (*domain.AppUser, error)
}

type visitorSurveySubmissionService struct {
	venues   repository.VenueRepository
	appUsers repository.AppUserRepository
}

func NewVisitorSurveySubmissionService(venues repository.VenueRepository, appUsers repository.AppUserRepository) VisitorSurveySubmissionService {
	return &visitorSurveySubmissionService{venues: venues, appUsers: appUsers}
}

func (s *visitorSurveySubmissionService) Submit(ctx context.Context, req VisitorSurveySubmission) (*domain.AppUser, error) {
	if strings.TrimSpace(req.PublicKey) == "" {
		return nil, domain.NewValidation(map[string]string{"public_key": "is required"})
	}
	if req.VisitorType == 0 {
		return nil, domain.NewValidation(map[string]string{"visitor_type": "is required"})
	}

	venue, err := s.venues.FindByPublicKey(ctx, req.PublicKey)
	if err != nil {
		return nil, err
	}

	firstName, lastName := splitFullName(req.FullName)
	interests, err := json.Marshal(dedupeInterests(req.Interests))
	if err != nil {
		return nil, err
	}
	visitorType := req.VisitorType

	u := &domain.AppUser{
		VenueID:        &venue.ID,
		Source:         "visitor_survey",
		FirstName:      firstName,
		LastName:       lastName,
		Email:          strings.TrimSpace(req.Email),
		Phone:          strings.TrimSpace(req.PhoneNumber),
		VisitorType:    &visitorType,
		Interests:      interests,
		OtherInterests: strings.TrimSpace(req.OtherInterests),
		IsConsented:    req.IsConsented,
		IPAddress:      strings.TrimSpace(req.IPAddress),
		UserAgent:      strings.TrimSpace(req.UserAgent),
		BusinessName:   strings.TrimSpace(req.BusinessName),
	}
	if err := s.appUsers.Create(ctx, u); err != nil {
		return nil, err
	}
	return u, nil
}

func splitFullName(v string) (string, string) {
	parts := strings.Fields(strings.TrimSpace(v))
	if len(parts) == 0 {
		return "", ""
	}
	if len(parts) == 1 {
		return parts[0], ""
	}
	return parts[0], strings.Join(parts[1:], " ")
}

func dedupeInterests(values []int) []int {
	if len(values) == 0 {
		return []int{}
	}
	out := make([]int, 0, len(values))
	seen := make(map[int]struct{}, len(values))
	for _, value := range values {
		if _, ok := seen[value]; ok {
			continue
		}
		seen[value] = struct{}{}
		out = append(out, value)
	}
	return out
}
