package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
)

// AESEncryption provides AES-256-CBC encrypt/decrypt matching Django's AESEncryption class
// (indoormap-backend/utils/encrypting.py). The key is the venue's URL-safe base64 public_key
// (43 chars, no padding, raw 32 bytes when decoded).
//
// Wire format: base64StdEncoding(iv[16] || ciphertext), PKCS7-padded to block size 16.
type AESEncryption struct {
	key []byte // 32 bytes (AES-256)
}

// NewAESEncryption decodes a URL-safe base64 key (no padding) to 32 raw key bytes.
func NewAESEncryption(urlSafeBase64Key string) (*AESEncryption, error) {
	key, err := base64.RawURLEncoding.DecodeString(urlSafeBase64Key)
	if err != nil {
		return nil, fmt.Errorf("decode aes key: %w", err)
	}
	if len(key) != 32 {
		return nil, fmt.Errorf("aes key must be 32 bytes, got %d", len(key))
	}
	return &AESEncryption{key: key}, nil
}

// Encrypt encrypts arbitrary bytes using AES-256-CBC with a random IV.
// Returns base64StdEncoding(iv || ciphertext) as a string, matching Django's encrypt/encrypt_bytes output.
func (a *AESEncryption) Encrypt(plaintext []byte) (string, error) {
	padded := pkcs7Pad(plaintext, aes.BlockSize)

	block, err := aes.NewCipher(a.key)
	if err != nil {
		return "", fmt.Errorf("new aes cipher: %w", err)
	}

	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", fmt.Errorf("generate iv: %w", err)
	}

	out := make([]byte, len(padded))
	cipher.NewCBCEncrypter(block, iv).CryptBlocks(out, padded)

	payload := append(iv, out...)
	return base64.StdEncoding.EncodeToString(payload), nil
}

// Decrypt reverses Encrypt. Returns the original unpadded plaintext.
func (a *AESEncryption) Decrypt(b64Ciphertext string) ([]byte, error) {
	payload, err := base64.StdEncoding.DecodeString(b64Ciphertext)
	if err != nil {
		return nil, fmt.Errorf("base64 decode: %w", err)
	}
	if len(payload) < aes.BlockSize {
		return nil, fmt.Errorf("ciphertext too short")
	}

	iv := payload[:aes.BlockSize]
	ciphertext := payload[aes.BlockSize:]

	if len(ciphertext)%aes.BlockSize != 0 {
		return nil, fmt.Errorf("ciphertext length not a multiple of block size")
	}

	block, err := aes.NewCipher(a.key)
	if err != nil {
		return nil, fmt.Errorf("new aes cipher: %w", err)
	}

	plaintext := make([]byte, len(ciphertext))
	cipher.NewCBCDecrypter(block, iv).CryptBlocks(plaintext, ciphertext)

	return pkcs7Unpad(plaintext)
}

// pkcs7Pad pads src to a multiple of blockSize using PKCS#7.
func pkcs7Pad(src []byte, blockSize int) []byte {
	pad := blockSize - len(src)%blockSize
	padded := make([]byte, len(src)+pad)
	copy(padded, src)
	for i := len(src); i < len(padded); i++ {
		padded[i] = byte(pad)
	}
	return padded
}

func pkcs7Unpad(src []byte) ([]byte, error) {
	if len(src) == 0 {
		return nil, fmt.Errorf("empty input")
	}
	pad := int(src[len(src)-1])
	if pad == 0 || pad > aes.BlockSize {
		return nil, fmt.Errorf("invalid pkcs7 padding value %d", pad)
	}
	return src[:len(src)-pad], nil
}
