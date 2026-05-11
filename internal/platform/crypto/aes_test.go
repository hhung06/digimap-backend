package crypto

import (
	"encoding/base64"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// testKey is a 32-byte URL-safe base64 key (43 chars, no padding), matching Django's venue.public_key format.
const testKey = "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA" // 32 zero bytes encoded

func TestAESEncryption_RoundTrip(t *testing.T) {
	enc, err := NewAESEncryption(testKey)
	require.NoError(t, err)

	cases := [][]byte{
		[]byte("hello"),
		[]byte(`{"compressed_data":"abc123"}`),
		make([]byte, 100),  // large payload
		make([]byte, 1),    // 1-byte payload
		make([]byte, 16),   // exact block size
		make([]byte, 17),   // block+1
	}
	for _, tc := range cases {
		ciphertext, err := enc.Encrypt(tc)
		require.NoError(t, err)

		got, err := enc.Decrypt(ciphertext)
		require.NoError(t, err)
		assert.Equal(t, tc, got)
	}
}

func TestAESEncryption_WireFormat(t *testing.T) {
	enc, err := NewAESEncryption(testKey)
	require.NoError(t, err)

	plaintext := []byte("test payload")
	ciphertext, err := enc.Encrypt(plaintext)
	require.NoError(t, err)

	// Must be valid standard base64
	payload, err := base64.StdEncoding.DecodeString(ciphertext)
	require.NoError(t, err)

	// First 16 bytes are IV; remainder is ciphertext (multiple of 16)
	assert.GreaterOrEqual(t, len(payload), 16+16, "must have at least IV + one block")
	assert.Equal(t, 0, (len(payload)-16)%16, "ciphertext must be block-aligned")
}

func TestAESEncryption_BadKey(t *testing.T) {
	_, err := NewAESEncryption("short")
	assert.Error(t, err)

	_, err = NewAESEncryption(strings.Repeat("A", 43) + "!")
	assert.Error(t, err)
}

func TestCompressJSON(t *testing.T) {
	data := map[string]any{"key": "value", "n": 42}
	gzipped, err := CompressJSON(data)
	require.NoError(t, err)
	assert.Greater(t, len(gzipped), 0)
}

func TestEncryptBundle_RoundTripShape(t *testing.T) {
	data := map[string]any{"search": []string{"a", "b"}}
	ciphertext, err := EncryptBundle(testKey, data)
	require.NoError(t, err)

	// Must be non-empty base64
	payload, err := base64.StdEncoding.DecodeString(ciphertext)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(payload), 32)
}

func TestEncryptBytes_RoundTripShape(t *testing.T) {
	data := []map[string]any{{"id": "abc", "is_top_location": true}}
	ciphertext, err := EncryptBytes(testKey, data)
	require.NoError(t, err)

	payload, err := base64.StdEncoding.DecodeString(ciphertext)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(payload), 32)
}
