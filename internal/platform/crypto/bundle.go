package crypto

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
)

// EncryptBundle compresses v as JSON, wraps it as {"compressed_data": "<gzip-base64>"},
// then AES-encrypts the wrapper dict. This mirrors the v2 bundle pipeline in Django:
//   publish_venue_v2.py:710 → AESEncryption.encrypt({"compressed_data": base64-gzip})
//
// Use this for snapshot language bundles.
func EncryptBundle(publicKey string, v any) (string, error) {
	enc, err := NewAESEncryption(publicKey)
	if err != nil {
		return "", err
	}

	gzipped, err := CompressJSON(v)
	if err != nil {
		return "", err
	}

	// Django wraps as {"compressed_data": "<base64-gzip-string>"} then encrypts the JSON of that dict.
	wrapper := map[string]string{
		"compressed_data": base64.StdEncoding.EncodeToString(gzipped),
	}
	wrapJSON, err := json.Marshal(wrapper)
	if err != nil {
		return "", fmt.Errorf("marshal wrapper: %w", err)
	}

	return enc.Encrypt(wrapJSON)
}

// EncryptBytes gzip-compresses raw bytes then AES-encrypts them directly (no dict wrapper).
// Mirrors Django's encrypt_bytes path used by top-location and memo bundles:
//   api/locations/views.py:410 → aes.encrypt_bytes(compress_json(data))
//
// Use this for top-location and memo bundles.
func EncryptBytes(publicKey string, v any) (string, error) {
	enc, err := NewAESEncryption(publicKey)
	if err != nil {
		return "", err
	}

	gzipped, err := CompressJSON(v)
	if err != nil {
		return "", err
	}

	return enc.Encrypt(gzipped)
}
