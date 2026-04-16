package enricher_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/hhung06/digimap-backend/internal/enricher"
)

// stubVenueRepo satisfies enricher.VenueCustomerResolver for tests.
type stubVenueRepo struct {
	customerID uuid.UUID
	err        error
}

func (s *stubVenueRepo) GetCustomerID(_ context.Context, _ uuid.UUID) (uuid.UUID, error) {
	return s.customerID, s.err
}

func TestRegistry_EnrichForVenue_NilRegistry(t *testing.T) {
	var reg *enricher.Registry
	extras, err := reg.EnrichForVenue(context.Background(), uuid.New(), enricher.ResourceLocation)
	require.NoError(t, err)
	assert.Nil(t, extras)
}

func TestRegistry_EnrichForVenue_EmptyRegistry(t *testing.T) {
	called := false
	repo := &trackingRepo{called: &called}
	reg := enricher.NewRegistry(repo)

	extras, err := reg.EnrichForVenue(context.Background(), uuid.New(), enricher.ResourceLocation)
	require.NoError(t, err)
	assert.Nil(t, extras)
	assert.False(t, called, "GetCustomerID must not be called on empty registry")
}

func TestRegistry_EnrichForVenue_NoMatchingEnricher(t *testing.T) {
	customerID := uuid.New()
	venueID := uuid.New()
	otherID := uuid.New()

	reg := enricher.NewRegistry(&stubVenueRepo{customerID: customerID})
	reg.Register(otherID, enricher.ResourceLocation, func(_ context.Context, _ uuid.UUID) (map[string]any, error) {
		return map[string]any{"key": "val"}, nil
	})

	extras, err := reg.EnrichForVenue(context.Background(), venueID, enricher.ResourceLocation)
	require.NoError(t, err)
	assert.Nil(t, extras)
}

func TestRegistry_EnrichForVenue_ReturnsExtras(t *testing.T) {
	customerID := uuid.New()
	venueID := uuid.New()

	reg := enricher.NewRegistry(&stubVenueRepo{customerID: customerID})
	reg.Register(customerID, enricher.ResourceLocation, func(_ context.Context, _ uuid.UUID) (map[string]any, error) {
		return map[string]any{"crm_zone": "west"}, nil
	})

	extras, err := reg.EnrichForVenue(context.Background(), venueID, enricher.ResourceLocation)
	require.NoError(t, err)
	assert.Equal(t, map[string]any{"crm_zone": "west"}, extras)
}

func TestRegistry_EnrichForVenue_RepoError(t *testing.T) {
	venueID := uuid.New()

	reg := enricher.NewRegistry(&stubVenueRepo{err: errors.New("db error")})
	reg.Register(uuid.New(), enricher.ResourceLocation, func(_ context.Context, _ uuid.UUID) (map[string]any, error) {
		return map[string]any{"key": "val"}, nil
	})

	extras, err := reg.EnrichForVenue(context.Background(), venueID, enricher.ResourceLocation)
	require.NoError(t, err) // repo errors are swallowed — base response unaffected
	assert.Nil(t, extras)
}

func TestMergeInto_NoExtras(t *testing.T) {
	type base struct {
		Name string `json:"name"`
	}
	result := enricher.MergeInto(base{Name: "test"}, nil)
	assert.Equal(t, base{Name: "test"}, result)
}

func TestMergeInto_EmptyExtras(t *testing.T) {
	type base struct {
		Name string `json:"name"`
	}
	result := enricher.MergeInto(base{Name: "test"}, map[string]any{})
	assert.Equal(t, base{Name: "test"}, result)
}

func TestMergeInto_WithExtras(t *testing.T) {
	type base struct {
		Name string `json:"name"`
	}
	result := enricher.MergeInto(base{Name: "test"}, map[string]any{"crm_zone": "west", "custom_label": "Gate"})

	m, ok := result.(map[string]any)
	require.True(t, ok)
	assert.Equal(t, "test", m["name"])
	assert.Equal(t, "west", m["crm_zone"])
	assert.Equal(t, "Gate", m["custom_label"])
}

// trackingRepo records whether GetCustomerID was called.
type trackingRepo struct {
	called *bool
}

func (t *trackingRepo) GetCustomerID(_ context.Context, _ uuid.UUID) (uuid.UUID, error) {
	*t.called = true
	return uuid.Nil, nil
}
