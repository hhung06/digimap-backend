package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
)

const gcmNonceSize = 12

// AESGCMEncryption provides AES-256-GCM encrypt/decrypt matching the Python AESGCM implementation
// in indoormap-backend/utils/encrypting.py (encrypt_data / decrypt_data).
//
// Wire format: base64StdEncoding(nonce[12] || ciphertext+auth_tag[n+16]).
// The key is 32 bytes encoded as standard or URL-safe base64 (Python tries URL-safe first).
type AESGCMEncryption struct {
	key []byte
}

// NewAESGCMEncryption decodes a base64 key (URL-safe or standard) to 32 raw key bytes.
func NewAESGCMEncryption(base64Key string) (*AESGCMEncryption, error) {
	key, err := base64.URLEncoding.DecodeString(base64Key)
	if err != nil {
		key, err = base64.StdEncoding.DecodeString(base64Key)
		if err != nil {
			return nil, fmt.Errorf("decode aes-gcm key: %w", err)
		}
	}
	if len(key) != 32 {
		return nil, fmt.Errorf("aes-gcm key must be 32 bytes, got %d", len(key))
	}
	return &AESGCMEncryption{key: key}, nil
}

// Encrypt encrypts plaintext and returns base64StdEncoding(nonce || ciphertext+tag).
func (a *AESGCMEncryption) Encrypt(plaintext string) (string, error) {
	block, err := aes.NewCipher(a.key)
	if err != nil {
		return "", fmt.Errorf("new aes cipher: %w", err)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("new gcm: %w", err)
	}
	nonce := make([]byte, gcmNonceSize)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("generate nonce: %w", err)
	}
	ciphertext := gcm.Seal(nil, nonce, []byte(plaintext), nil)
	return base64.StdEncoding.EncodeToString(append(nonce, ciphertext...)), nil
}

// Decrypt reverses Encrypt. Accepts base64StdEncoding(nonce || ciphertext+tag).
func (a *AESGCMEncryption) Decrypt(encoded string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return "", fmt.Errorf("base64 decode: %w", err)
	}
	if len(data) < gcmNonceSize {
		return "", fmt.Errorf("encrypted data too short")
	}
	nonce := data[:gcmNonceSize]
	ciphertext := data[gcmNonceSize:]

	block, err := aes.NewCipher(a.key)
	if err != nil {
		return "", fmt.Errorf("new aes cipher: %w", err)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("new gcm: %w", err)
	}
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", fmt.Errorf("gcm open: %w", err)
	}
	return string(plaintext), nil
}
