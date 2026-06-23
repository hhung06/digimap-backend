package dto_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/hhung06/digimap-backend/internal/domain"
	"github.com/hhung06/digimap-backend/internal/dto"
)

func TestArticleRequestPeriodFieldsAcceptDateOnly(t *testing.T) {
	var req dto.ArticleRequest

	err := json.Unmarshal([]byte(`{
		"title":"News",
		"published_period_start":"2026-06-01",
		"published_period_end":"2026-06-30"
	}`), &req)

	require.NoError(t, err)
	require.NotNil(t, req.PublishedPeriodStart)
	require.NotNil(t, req.PublishedPeriodEnd)
	assert.Equal(t, time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC), req.PublishedPeriodStart.Time())
	assert.Equal(t, time.Date(2026, 6, 30, 0, 0, 0, 0, time.UTC), req.PublishedPeriodEnd.Time())
}

func TestArticleRequestPeriodFieldsRejectDatetimeStrings(t *testing.T) {
	var req dto.ArticleRequest

	err := json.Unmarshal([]byte(`{
		"title":"News",
		"published_period_start":"2026-06-01T00:00:00Z"
	}`), &req)

	require.Error(t, err)
}

func TestArticleResponsePeriodFieldsEmitDateOnly(t *testing.T) {
	start := time.Date(2026, 6, 1, 15, 45, 0, 0, time.FixedZone("JST", 9*60*60))
	end := time.Date(2026, 6, 30, 23, 59, 0, 0, time.UTC)
	resp := dto.ArticleToResponse(&domain.Article{
		ID:                   uuid.New(),
		Title:                "News",
		PublishedPeriodStart: &start,
		PublishedPeriodEnd:   &end,
	})

	body, err := json.Marshal(resp)
	require.NoError(t, err)

	assert.Contains(t, string(body), `"published_period_start":"2026-06-01"`)
	assert.Contains(t, string(body), `"published_period_end":"2026-06-30"`)
	assert.NotContains(t, string(body), "T15:45")
	assert.NotContains(t, string(body), "T23:59")
}

func TestArticleImageResponseIncludesStoredKeyAndNullableURL(t *testing.T) {
	resp := dto.ArticleImageToResponse(&domain.ArticleImage{
		ID:        uuid.New(),
		ArticleID: uuid.New(),
		Image:     "develop/media/articles/article-id/images/image-id.png",
	})

	body, err := json.Marshal(resp)
	require.NoError(t, err)

	assert.Contains(t, string(body), `"image":"develop/media/articles/article-id/images/image-id.png"`)
	assert.Contains(t, string(body), `"image_url":null`)
}

func TestUpdateArticleRequestRemoveImagesDistinguishesOmittedFalseAndTrue(t *testing.T) {
	var omitted dto.UpdateArticleRequest
	require.NoError(t, json.Unmarshal([]byte(`{"title":"keep"}`), &omitted))
	assert.Nil(t, omitted.RemoveImages)

	var keep dto.UpdateArticleRequest
	require.NoError(t, json.Unmarshal([]byte(`{"remove_images":false}`), &keep))
	require.NotNil(t, keep.RemoveImages)
	assert.False(t, *keep.RemoveImages)

	var clear dto.UpdateArticleRequest
	require.NoError(t, json.Unmarshal([]byte(`{"remove_images":true}`), &clear))
	require.NotNil(t, clear.RemoveImages)
	assert.True(t, *clear.RemoveImages)
}
