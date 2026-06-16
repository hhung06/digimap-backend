package search

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

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

// BulkDoc is a document to be written in a bulk indexing operation.
type BulkDoc struct {
	ID     string
	Source any
}

// Searcher abstracts full-text search operations against an OpenSearch cluster.
type Searcher interface {
	// Search performs a query on the given index. query is raw JSON (e.g. a match/multi_match query body).
	Search(ctx context.Context, index string, query map[string]interface{}, from, size int) (*SearchResult, error)
	// Index writes or updates a document.
	Index(ctx context.Context, index string, docID string, doc interface{}) error
	// Delete removes a document by ID.
	Delete(ctx context.Context, index string, docID string) error

	// Index management
	CreateIndex(ctx context.Context, name string, body map[string]any) error
	DeleteIndex(ctx context.Context, name string) error

	// Bulk indexing
	BulkIndex(ctx context.Context, index string, docs []BulkDoc) error

	// Alias management
	AliasExists(ctx context.Context, name string) (bool, error)
	GetIndexForAlias(ctx context.Context, name string) (string, error)
	PutAlias(ctx context.Context, index, alias string) error
	SwapAlias(ctx context.Context, oldIndex, newIndex, alias string) error
}

// compile-time interface checks
var _ Searcher = (*osSearcher)(nil)
var _ Searcher = (*LogSearcher)(nil)

const bulkChunkSize = 250

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

func (s *osSearcher) CreateIndex(ctx context.Context, name string, body map[string]any) error {
	b, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("marshal index body: %w", err)
	}
	_, err = s.client.Indices.Create(ctx, opensearchapi.IndicesCreateReq{
		Index: name,
		Body:  bytes.NewReader(b),
	})
	if err != nil {
		return fmt.Errorf("opensearch create index %q: %w", name, err)
	}
	return nil
}

func (s *osSearcher) DeleteIndex(ctx context.Context, name string) error {
	_, err := s.client.Indices.Delete(ctx, opensearchapi.IndicesDeleteReq{
		Indices: []string{name},
	})
	if err != nil {
		return fmt.Errorf("opensearch delete index %q: %w", name, err)
	}
	return nil
}

func (s *osSearcher) BulkIndex(ctx context.Context, index string, docs []BulkDoc) error {
	for i := 0; i < len(docs); i += bulkChunkSize {
		end := i + bulkChunkSize
		if end > len(docs) {
			end = len(docs)
		}
		chunk := docs[i:end]

		var buf bytes.Buffer
		for _, d := range chunk {
			action := map[string]any{
				"index": map[string]any{
					"_index": index,
					"_id":    d.ID,
				},
			}
			actionLine, err := json.Marshal(action)
			if err != nil {
				return fmt.Errorf("marshal bulk action: %w", err)
			}
			sourceLine, err := json.Marshal(d.Source)
			if err != nil {
				return fmt.Errorf("marshal bulk source: %w", err)
			}
			buf.Write(actionLine)
			buf.WriteByte('\n')
			buf.Write(sourceLine)
			buf.WriteByte('\n')
		}

		resp, err := s.client.Bulk(ctx, opensearchapi.BulkReq{
			Body: bytes.NewReader(buf.Bytes()),
		})
		if err != nil {
			return fmt.Errorf("opensearch bulk (chunk %d): %w", i/bulkChunkSize, err)
		}
		if resp.Errors {
			if err := bulkRespItemsError(resp.Items); err != nil {
				return fmt.Errorf("opensearch bulk item failures (chunk %d): %w", i/bulkChunkSize, err)
			}
			return fmt.Errorf("opensearch bulk reported item failures (chunk %d)", i/bulkChunkSize)
		}
	}
	return nil
}

func bulkItemsError(body []byte) error {
	var resp struct {
		Errors bool                                    `json:"errors"`
		Items  []map[string]opensearchapi.BulkRespItem `json:"items"`
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		return fmt.Errorf("decode bulk response: %w", err)
	}
	if !resp.Errors {
		return nil
	}
	return bulkRespItemsError(resp.Items)
}

func bulkRespItemsError(items []map[string]opensearchapi.BulkRespItem) error {
	var failures []string
	for _, item := range items {
		for op, ri := range item {
			if ri.Error != nil {
				failures = append(failures, fmt.Sprintf("op=%s id=%s type=%s reason=%s", op, ri.ID, ri.Error.Type, ri.Error.Reason))
			}
		}
	}
	if len(failures) == 0 {
		return fmt.Errorf("bulk response contains item failures")
	}
	return errors.New(strings.Join(failures, "; "))
}

func (s *osSearcher) AliasExists(ctx context.Context, name string) (bool, error) {
	resp, err := s.client.Indices.Alias.Exists(ctx, opensearchapi.AliasExistsReq{
		Alias: []string{name},
	})
	if err != nil {
		// A 404 from the API is returned as an error by the client; treat it as false.
		return false, nil
	}
	return resp.StatusCode == 200, nil
}

func (s *osSearcher) GetIndexForAlias(ctx context.Context, name string) (string, error) {
	resp, err := s.client.Indices.Alias.Get(ctx, opensearchapi.AliasGetReq{
		Alias: []string{name},
	})
	if err != nil {
		return "", fmt.Errorf("opensearch get alias %q: %w", name, err)
	}
	for indexName := range resp.Indices {
		return indexName, nil
	}
	return "", fmt.Errorf("alias %q not found", name)
}

func (s *osSearcher) PutAlias(ctx context.Context, index, alias string) error {
	_, err := s.client.Indices.Alias.Put(ctx, opensearchapi.AliasPutReq{
		Indices: []string{index},
		Alias:   alias,
	})
	if err != nil {
		return fmt.Errorf("opensearch put alias %q -> %q: %w", alias, index, err)
	}
	return nil
}

func (s *osSearcher) SwapAlias(ctx context.Context, oldIndex, newIndex, alias string) error {
	body, err := json.Marshal(map[string]any{
		"actions": []any{
			map[string]any{"remove": map[string]any{"index": oldIndex, "alias": alias}},
			map[string]any{"add": map[string]any{"index": newIndex, "alias": alias}},
		},
	})
	if err != nil {
		return fmt.Errorf("marshal swap alias body: %w", err)
	}
	resp, err := s.client.Aliases(ctx, opensearchapi.AliasesReq{
		Body: bytes.NewReader(body),
	})
	if err != nil {
		return fmt.Errorf("opensearch swap alias %q (%s->%s): %w", alias, oldIndex, newIndex, err)
	}
	if !resp.Acknowledged {
		return fmt.Errorf("opensearch swap alias %q not acknowledged", alias)
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

func (s *LogSearcher) CreateIndex(_ context.Context, name string, _ map[string]any) error {
	fmt.Printf("[SEARCH] CreateIndex name=%s\n", name)
	return nil
}

func (s *LogSearcher) DeleteIndex(_ context.Context, name string) error {
	fmt.Printf("[SEARCH] DeleteIndex name=%s\n", name)
	return nil
}

func (s *LogSearcher) BulkIndex(_ context.Context, index string, docs []BulkDoc) error {
	fmt.Printf("[SEARCH] BulkIndex index=%s docs=%d\n", index, len(docs))
	return nil
}

func (s *LogSearcher) AliasExists(_ context.Context, name string) (bool, error) {
	fmt.Printf("[SEARCH] AliasExists name=%s\n", name)
	return false, nil
}

func (s *LogSearcher) GetIndexForAlias(_ context.Context, name string) (string, error) {
	fmt.Printf("[SEARCH] GetIndexForAlias name=%s\n", name)
	return "", nil
}

func (s *LogSearcher) PutAlias(_ context.Context, index, alias string) error {
	fmt.Printf("[SEARCH] PutAlias index=%s alias=%s\n", index, alias)
	return nil
}

func (s *LogSearcher) SwapAlias(_ context.Context, oldIndex, newIndex, alias string) error {
	fmt.Printf("[SEARCH] SwapAlias old=%s new=%s alias=%s\n", oldIndex, newIndex, alias)
	return nil
}
