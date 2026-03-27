package search

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/opensearch-project/opensearch-go/v4"
	"github.com/opensearch-project/opensearch-go/v4/opensearchapi"
)

// Hit is a single OpenSearch document hit.
type Hit struct {
	ID     string
	Score  float32
	Source json.RawMessage
}

// SearchResult holds the response from an OpenSearch query.
type SearchResult struct {
	Total int64
	Hits  []Hit
}

// Searcher abstracts full-text search operations against an OpenSearch cluster.
type Searcher interface {
	// Search performs a query on the given index. query is raw JSON (e.g. a match/multi_match query body).
	Search(ctx context.Context, index string, query map[string]interface{}, from, size int) (*SearchResult, error)
	// Index writes or updates a document.
	Index(ctx context.Context, index string, docID string, doc interface{}) error
	// Delete removes a document by ID.
	Delete(ctx context.Context, index string, docID string) error
}

type osSearcher struct {
	client *opensearchapi.Client
}

// NewSearcher creates a Searcher backed by a real OpenSearch cluster.
func NewSearcher(endpoint, username, password string) (Searcher, error) {
	cfg := opensearchapi.Config{
		Client: opensearch.Config{
			Addresses: []string{endpoint},
		},
	}
	if username != "" {
		cfg.Client.Username = username
		cfg.Client.Password = password
	}
	client, err := opensearchapi.NewClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("create opensearch client: %w", err)
	}
	return &osSearcher{client: client}, nil
}

func (s *osSearcher) Search(ctx context.Context, index string, query map[string]interface{}, from, size int) (*SearchResult, error) {
	body, err := json.Marshal(map[string]interface{}{
		"from":  from,
		"size":  size,
		"query": query,
	})
	if err != nil {
		return nil, err
	}
	req := &opensearchapi.SearchReq{
		Indices: []string{index},
		Body:    bytes.NewReader(body),
	}
	resp, err := s.client.Search(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("opensearch search: %w", err)
	}
	result := &SearchResult{Total: int64(resp.Hits.Total.Value)}
	for _, h := range resp.Hits.Hits {
		result.Hits = append(result.Hits, Hit{ID: h.ID, Score: h.Score, Source: h.Source})
	}
	return result, nil
}

func (s *osSearcher) Index(ctx context.Context, index, docID string, doc interface{}) error {
	body, err := json.Marshal(doc)
	if err != nil {
		return err
	}
	req := opensearchapi.IndexReq{
		Index:      index,
		DocumentID: docID,
		Body:       bytes.NewReader(body),
	}
	_, err = s.client.Index(ctx, req)
	if err != nil {
		return fmt.Errorf("opensearch index: %w", err)
	}
	return nil
}

func (s *osSearcher) Delete(ctx context.Context, index, docID string) error {
	req := opensearchapi.DocumentDeleteReq{
		Index:      index,
		DocumentID: docID,
	}
	_, err := s.client.Document.Delete(ctx, req)
	if err != nil {
		return fmt.Errorf("opensearch delete: %w", err)
	}
	return nil
}

// LogSearcher is a no-op Searcher for local development.
type LogSearcher struct{}

func NewLogSearcher() Searcher { return &LogSearcher{} }

func (s *LogSearcher) Search(_ context.Context, index string, _ map[string]interface{}, from, size int) (*SearchResult, error) {
	fmt.Printf("[DEV SEARCH] index=%s from=%d size=%d\n", index, from, size)
	return &SearchResult{}, nil
}

func (s *LogSearcher) Index(_ context.Context, index, docID string, _ interface{}) error {
	fmt.Printf("[DEV SEARCH] index doc index=%s id=%s\n", index, docID)
	return nil
}

func (s *LogSearcher) Delete(_ context.Context, index, docID string) error {
	fmt.Printf("[DEV SEARCH] delete doc index=%s id=%s\n", index, docID)
	return nil
}
