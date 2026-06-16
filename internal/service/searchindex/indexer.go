package searchindex

import (
	"context"
	"fmt"
	"time"

	"github.com/hhung06/digimap-backend/internal/domain"
	"github.com/hhung06/digimap-backend/internal/platform/search"
	applog "github.com/hhung06/digimap-backend/log"
)

// IndexVenueAsync indexes all exhibitor and product documents for a venue into
// OpenSearch using a blue/green alias swap. It must be called from a goroutine
// (publishV2 is already async). Errors are logged and never propagated.
func IndexVenueAsync(
	ctx context.Context,
	logger applog.Logger,
	searcher search.Searcher,
	venue *domain.Venue,
	locations []*domain.Location,
	products []*domain.Product,
	environment string,
) {
	alias := fmt.Sprintf("%s_%s", venue.ExternalID, environment)
	index := fmt.Sprintf("%s_%s_v%d", venue.ExternalID, environment, time.Now().Unix())

	// 1. Create the new index with the venue mapping.
	if err := searcher.CreateIndex(ctx, index, search.VenueIndexMapping()); err != nil {
		logger.Errorf("[searchindex] create index %q: %v", index, err)
		return
	}

	// 2. Build and bulk-index all documents.
	docs := buildDocs(venue, locations, products)
	if err := searcher.BulkIndex(ctx, index, docs); err != nil {
		logger.Errorf("[searchindex] bulk index venue=%s index=%s docs=%d: %v", venue.ExternalID, index, len(docs), err)
		return
	}

	logger.Infof("[searchindex] indexed venue=%s index=%s docs=%d", venue.ExternalID, index, len(docs))

	// 3. Swap (or create) the alias.
	aliasExists, err := searcher.AliasExists(ctx, alias)
	if err != nil {
		logger.Errorf("[searchindex] alias exists check alias=%q: %v", alias, err)
		return
	}

	if aliasExists {
		oldIndex, err := searcher.GetIndexForAlias(ctx, alias)
		if err != nil {
			logger.Errorf("[searchindex] get index for alias=%q: %v", alias, err)
			return
		}
		if err := searcher.SwapAlias(ctx, oldIndex, index, alias); err != nil {
			logger.Errorf("[searchindex] swap alias=%q %s->%s: %v", alias, oldIndex, index, err)
			return
		}
		// Best-effort delete old index.
		if err := searcher.DeleteIndex(ctx, oldIndex); err != nil {
			logger.Warnf("[searchindex] delete old index=%q: %v", oldIndex, err)
		}
	} else {
		if err := searcher.PutAlias(ctx, index, alias); err != nil {
			logger.Errorf("[searchindex] put alias=%q -> %q: %v", alias, index, err)
			return
		}
	}

	logger.Infof("[searchindex] alias=%q now points to index=%q", alias, index)
}
