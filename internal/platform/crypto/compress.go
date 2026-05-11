package crypto

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
)

// CompressJSON serializes v as compact JSON (no whitespace), then gzip-compresses it at best compression.
// Mirrors Django's compress_json (indoormap-backend/utils/encrypting.py:118).
func CompressJSON(v any) ([]byte, error) {
	raw, err := json.Marshal(v)
	if err != nil {
		return nil, fmt.Errorf("json marshal: %w", err)
	}
	var buf bytes.Buffer
	w, err := gzip.NewWriterLevel(&buf, gzip.BestCompression)
	if err != nil {
		return nil, fmt.Errorf("gzip writer: %w", err)
	}
	if _, err := w.Write(raw); err != nil {
		return nil, fmt.Errorf("gzip write: %w", err)
	}
	if err := w.Close(); err != nil {
		return nil, fmt.Errorf("gzip close: %w", err)
	}
	return buf.Bytes(), nil
}
