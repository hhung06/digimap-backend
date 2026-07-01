package service

import (
	"bytes"
	"context"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"mime"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/hhung06/digimap-backend/internal/domain"
	"github.com/hhung06/digimap-backend/internal/platform/storage"
)

const maxMediaUploadBytes int64 = 10 << 20
const maxMediaImageWidth int64 = 4096
const maxMediaImageHeight int64 = 4096
const maxMediaImagePixels int64 = 16_000_000

type MediaTarget struct {
	Entity   string
	RecordID uuid.UUID
	Field    string
}

type MediaUpload struct {
	Filename    string
	ContentType string
	Size        int64
	Reader      io.Reader
}

type MediaService interface {
	Upload(context.Context, MediaTarget, MediaUpload) (string, error)
	URL(context.Context, MediaTarget, string) *string
	DeleteOwned(context.Context, MediaTarget, string) error
}

type mediaService struct {
	storer storage.Storer
	env    string
}

func NewMediaService(storer storage.Storer, env string) MediaService {
	return &mediaService{storer: storer, env: env}
}

func (s *mediaService) Upload(ctx context.Context, target MediaTarget, upload MediaUpload) (string, error) {
	if err := validateMediaTarget(s.env, target); err != nil {
		return "", err
	}
	body, contentType, ext, err := validateMediaUpload(upload)
	if err != nil {
		return "", err
	}

	key := storage.MediaKey(s.env, target.Entity, target.RecordID, target.Field, uuid.New(), ext)
	if key == "" {
		return "", domain.NewValidation(map[string]string{"target": "invalid media target"})
	}
	if err := s.storer.PutMedia(ctx, key, contentType, int64(len(body)), bytes.NewReader(body)); err != nil {
		return "", err
	}
	return key, nil
}

func (s *mediaService) URL(ctx context.Context, target MediaTarget, key string) *string {
	if strings.TrimSpace(key) == "" {
		return nil
	}
	if err := validateMediaTarget(s.env, target); err != nil {
		return nil
	}
	if !storage.OwnsMediaKey(s.env, target.Entity, target.RecordID, target.Field, key) {
		return nil
	}
	url, err := s.storer.PresignDownload(ctx, key, 15*time.Minute)
	if err != nil {
		return nil
	}
	return &url
}

func (s *mediaService) DeleteOwned(ctx context.Context, target MediaTarget, key string) error {
	if err := validateMediaTarget(s.env, target); err != nil {
		return err
	}
	if strings.TrimSpace(key) == "" {
		return domain.NewValidation(map[string]string{"key": "required"})
	}
	if !storage.OwnsMediaKey(s.env, target.Entity, target.RecordID, target.Field, key) {
		return domain.NewValidation(map[string]string{"key": "media key does not belong to target"})
	}
	return s.storer.DeleteObject(ctx, key)
}

func validateMediaTarget(env string, target MediaTarget) error {
	if target.RecordID == uuid.Nil {
		return domain.NewValidation(map[string]string{"record_id": "required"})
	}
	if storage.MediaKey(env, target.Entity, target.RecordID, target.Field, uuid.New(), "png") == "" {
		return domain.NewValidation(map[string]string{"target": "invalid media target"})
	}
	return nil
}

func validateMediaUpload(upload MediaUpload) ([]byte, string, string, error) {
	details := map[string]string{}
	if strings.TrimSpace(upload.Filename) == "" {
		details["filename"] = "required"
	}
	if strings.TrimSpace(upload.ContentType) == "" {
		details["content_type"] = "required"
	}
	if upload.Size <= 0 {
		details["size"] = "must be positive"
	} else if upload.Size > maxMediaUploadBytes {
		details["size"] = "too large"
	}
	if upload.Reader == nil {
		details["reader"] = "required"
	}
	if len(details) > 0 {
		return nil, "", "", domain.NewValidation(details)
	}

	declaredContentType, _, err := mime.ParseMediaType(upload.ContentType)
	if err != nil {
		return nil, "", "", domain.NewValidation(map[string]string{"content_type": "invalid"})
	}
	declaredContentType = strings.ToLower(declaredContentType)

	body, err := io.ReadAll(io.LimitReader(upload.Reader, maxMediaUploadBytes+1))
	if err != nil {
		return nil, "", "", err
	}
	if int64(len(body)) > maxMediaUploadBytes {
		return nil, "", "", domain.NewValidation(map[string]string{"size": "too large"})
	}
	if int64(len(body)) != upload.Size {
		return nil, "", "", domain.NewValidation(map[string]string{"size": "does not match body length"})
	}

	detectedContentType := http.DetectContentType(body)
	if ext, ok := documentExtensionByContentType(declaredContentType, upload.Filename, detectedContentType); ok {
		return body, declaredContentType, ext, nil
	}
	if _, ok := mediaExtensionByContentType(declaredContentType); !ok {
		return nil, "", "", domain.NewValidation(map[string]string{"content_type": "unsupported media type"})
	}
	if detectedContentType != declaredContentType {
		return nil, "", "", domain.NewValidation(map[string]string{"content_type": "does not match file content"})
	}
	cfg, format, err := image.DecodeConfig(bytes.NewReader(body))
	if err != nil {
		return nil, "", "", domain.NewValidation(map[string]string{"image": "undecodable"})
	}
	if !validMediaDimensions(cfg) {
		return nil, "", "", domain.NewValidation(map[string]string{"image": "dimensions too large"})
	}
	ext, ok := mediaExtensionByImageFormat(format)
	if !ok {
		return nil, "", "", domain.NewValidation(map[string]string{"content_type": "unsupported image type"})
	}
	if _, _, err := image.Decode(bytes.NewReader(body)); err != nil {
		return nil, "", "", domain.NewValidation(map[string]string{"image": "undecodable"})
	}
	return body, detectedContentType, ext, nil
}

func documentExtensionByContentType(contentType, filename, detectedContentType string) (string, bool) {
	ext := strings.ToLower(filepath.Ext(filename))
	switch contentType {
	case "application/pdf":
		if ext == ".pdf" && strings.HasPrefix(detectedContentType, "application/pdf") {
			return "pdf", true
		}
	case "application/vnd.openxmlformats-officedocument.wordprocessingml.document":
		if ext == ".docx" && (detectedContentType == "application/zip" || detectedContentType == "application/octet-stream") {
			return "docx", true
		}
	}
	return "", false
}

func validMediaDimensions(cfg image.Config) bool {
	if cfg.Width <= 0 || cfg.Height <= 0 {
		return false
	}
	width := int64(cfg.Width)
	height := int64(cfg.Height)
	if width > maxMediaImageWidth || height > maxMediaImageHeight {
		return false
	}
	return width <= maxMediaImagePixels/height
}

func mediaExtensionByContentType(contentType string) (string, bool) {
	switch contentType {
	case "image/gif":
		return "gif", true
	case "image/jpeg":
		return "jpg", true
	case "image/png":
		return "png", true
	default:
		return "", false
	}
}

func mediaExtensionByImageFormat(format string) (string, bool) {
	switch format {
	case "gif", "png":
		return format, true
	case "jpeg":
		return "jpg", true
	default:
		return "", false
	}
}
