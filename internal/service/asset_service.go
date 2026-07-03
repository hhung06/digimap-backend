package service

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"mime"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/google/uuid"

	"github.com/hhung06/digimap-backend/internal/domain"
	"github.com/hhung06/digimap-backend/internal/platform/cdn"
	"github.com/hhung06/digimap-backend/internal/platform/storage"
	"github.com/hhung06/digimap-backend/internal/repository"
)

type CreateAssetInput struct {
	VenueID     *uuid.UUID
	Name        string
	Key         string
	ContentType string
	SizeBytes   int64
	URL         string
	CreatedBy   *uuid.UUID
}

type Upload3DAssetInput struct {
	VenueID     *uuid.UUID
	ID          string // optional; if set and asset exists, it's an update
	Name        string
	Description string
	FileType    string
	Width       *float64
	Height      *float64
	CreatedBy   *uuid.UUID

	// Each FileUpload field is optional (nil = not provided / no change).
	File      *FileUpload
	Thumbnail *FileUpload
	Material  *FileUpload
}

type FileUpload struct {
	Filename    string
	ContentType string
	Size        int64
	Reader      io.Reader
}

type UploadLibraryAssetInput struct {
	VenueID   *uuid.UUID
	Status    string
	AssetType string
	CreatedBy *uuid.UUID
	File      *FileUpload
}

type AssetService interface {
	List(ctx context.Context, venueID uuid.UUID, p domain.Pagination, assetType string) ([]*domain.Asset, int64, error)
	ListAll(ctx context.Context, assetType string) ([]*domain.Asset, error)
	ListLibrary(ctx context.Context, status, assetType string) ([]*domain.Asset, error)
	Get(ctx context.Context, id uuid.UUID) (*domain.Asset, error)
	// GetByIDs returns assets by their UUIDs, scoped to the given venue.
	GetByIDs(ctx context.Context, venueID uuid.UUID, ids []uuid.UUID) ([]*domain.Asset, error)
	Create(ctx context.Context, in CreateAssetInput) (*domain.Asset, error)
	// Upload2D decodes a base64 data-URL, uploads it to S3, and registers the asset record.
	// id is the client-assigned identifier; non-UUID strings (e.g. SHA-1 hashes) are mapped
	// to a deterministic UUID v5 so the same file always resolves to the same record.
	Upload2D(ctx context.Context, venueID uuid.UUID, id, dataURL string) (*domain.Asset, error)
	Update(ctx context.Context, id uuid.UUID, name string) (*domain.Asset, error)
	Delete(ctx context.Context, id uuid.UUID) error
	Upload3D(ctx context.Context, in Upload3DAssetInput) (*domain.Asset, error)
	UploadLibrary(ctx context.Context, in UploadLibraryAssetInput) (*domain.Asset, error)
	UpdateLibrary(ctx context.Context, id uuid.UUID, in UploadLibraryAssetInput) (*domain.Asset, error)
	// AssetURL returns the CloudFront URL for an S3 key, or "" if domain is not configured.
	AssetURL(key string) string
}

type assetService struct {
	repo        repository.AssetRepository
	storer      storage.Storer
	invalidator cdn.Invalidator
	env         string
	cfDomain    string // AWS_CF_ASSETS_DOMAIN, e.g. "assets.example.com"
}

func NewAssetService(
	repo repository.AssetRepository,
	storer storage.Storer,
	invalidator cdn.Invalidator,
	env, cfDomain string,
) AssetService {
	return &assetService{repo: repo, storer: storer, invalidator: invalidator, env: env, cfDomain: cfDomain}
}

func (s *assetService) AssetURL(key string) string {
	if key == "" || s.cfDomain == "" {
		return ""
	}
	return "https://" + s.cfDomain + "/" + key
}

func (s *assetService) List(ctx context.Context, venueID uuid.UUID, p domain.Pagination, assetType string) ([]*domain.Asset, int64, error) {
	return s.repo.List(ctx, venueID, p, assetType)
}

func (s *assetService) ListAll(ctx context.Context, assetType string) ([]*domain.Asset, error) {
	return s.repo.ListAll(ctx, assetType)
}

func (s *assetService) ListLibrary(ctx context.Context, status, assetType string) ([]*domain.Asset, error) {
	return s.repo.ListLibrary(ctx, s.env, status, assetType)
}

func (s *assetService) Get(ctx context.Context, id uuid.UUID) (*domain.Asset, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *assetService) GetByIDs(ctx context.Context, venueID uuid.UUID, ids []uuid.UUID) ([]*domain.Asset, error) {
	return s.repo.FindByIDs(ctx, venueID, ids)
}

func (s *assetService) Create(ctx context.Context, in CreateAssetInput) (*domain.Asset, error) {
	a := &domain.Asset{
		VenueID:     in.VenueID,
		Name:        in.Name,
		Key:         in.Key,
		ContentType: in.ContentType,
		SizeBytes:   in.SizeBytes,
		URL:         in.URL,
		AssetType:   domain.AssetType2D,
		CreatedBy:   in.CreatedBy,
	}
	if err := s.repo.Create(ctx, a); err != nil {
		return nil, err
	}
	return a, nil
}

// allowed2DMimeTypes is the allowlist for canvas image uploads.
// SVG is excluded because it can embed JavaScript.
var allowed2DMimeTypes = map[string]string{
	"image/png":  "png",
	"image/jpeg": "jpg",
	"image/webp": "webp",
	"image/gif":  "gif",
	"image/avif": "avif",
	"image/heic": "heic",
	"image/heif": "heif",
}

func (s *assetService) Upload2D(ctx context.Context, venueID uuid.UUID, id, dataURL string) (*domain.Asset, error) {
	assetID, err := uuid.Parse(id)
	if err != nil {
		// Client sent a non-UUID identifier (e.g. SHA-1 hex hash from the editor).
		// Derive a deterministic UUID v5 so the same file always maps to the same record.
		assetID = uuid.NewSHA1(uuid.NameSpaceOID, []byte(id))
	}

	// If the asset already exists, verify it belongs to this venue (prevent cross-venue overwrite).
	existing, existingErr := s.repo.FindByID(ctx, assetID)
	if existingErr == nil {
		if existing.VenueID == nil || *existing.VenueID != venueID {
			return nil, domain.NewValidation(map[string]string{"id": "asset belongs to a different venue"})
		}
	}

	// Parse "data:<contentType>;base64,<data>" — only used to verify client claim.
	comma := strings.Index(dataURL, ",")
	if comma < 0 {
		return nil, domain.NewValidation(map[string]string{"file": "invalid data URL"})
	}
	encoded := dataURL[comma+1:]

	raw, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return nil, domain.NewValidation(map[string]string{"file": "invalid base64 encoding"})
	}

	// Detect content-type from actual bytes — do not trust the data-URL header.
	detected := http.DetectContentType(raw)
	mediaType, _, _ := mime.ParseMediaType(detected)
	ext, ok := allowed2DMimeTypes[mediaType]
	if !ok {
		// AVIF, HEIC, and HEIF use the ISOBMFF container; Go's detector returns "video/mp4" for them.
		// Fall back to the claimed MIME type from the data-URL header for these formats.
		if mediaType == "video/mp4" || mediaType == "application/octet-stream" {
			if semi := strings.Index(dataURL, ";"); semi > 5 {
				claimedType := dataURL[5:semi] // "data:<type>;base64,..." → type
				if claimedExt, allowed := allowed2DMimeTypes[claimedType]; allowed {
					mediaType = claimedType
					ext = claimedExt
					ok = true
				}
			}
		}
	}
	if !ok {
		return nil, domain.NewValidation(map[string]string{"file": "unsupported file type; allowed: png, jpeg, webp, gif, avif"})
	}

	// The key embeds the detected extension, so re-uploading different content under the
	// same id (e.g. replacing a JPEG floor plan with a PNG) can produce a DIFFERENT key
	// than the one already stored — writing the old key never gets overwritten. Always
	// write to the freshly detected key, then persist it below.
	key := storage.Asset2DKey(s.env, assetID, ext)
	if err := s.storer.PutMedia(ctx, key, mediaType, int64(len(raw)), bytes.NewReader(raw)); err != nil {
		return nil, fmt.Errorf("upload 2D asset: %w", err)
	}

	// CloudFront caches by URL/key; invalidate unconditionally on every upload so a
	// same-key overwrite (same format re-uploaded) is picked up immediately, mirroring
	// Django's invalidate_cloudfront call (indoormap-backend/api/assets/views.py:213).
	if _, err := s.invalidator.Invalidate(ctx, []string{"/" + key}); err != nil {
		fmt.Printf("asset cf invalidation failed key=%s: %v\n", key, err)
	}

	a := &domain.Asset{
		ID:          assetID,
		VenueID:     &venueID,
		Name:        id,
		Key:         key,
		ContentType: mediaType,
		SizeBytes:   int64(len(raw)),
		URL:         s.AssetURL(key),
		AssetType:   domain.AssetType2D,
	}
	if existingErr == nil {
		if updateErr := s.repo.Update(ctx, a); updateErr != nil {
			return nil, fmt.Errorf("update 2D asset: %w", updateErr)
		}
		// The new content resolved to a different key (format changed) — drop the orphaned
		// old object so it doesn't linger in storage, matching UpdateLibrary's cleanup pattern.
		if existing.Key != "" && existing.Key != key {
			_ = s.storer.DeleteObject(ctx, existing.Key)
		}
	} else if createErr := s.repo.Create(ctx, a); createErr != nil {
		return nil, fmt.Errorf("create 2D asset: %w", createErr)
	}
	return a, nil
}

func (s *assetService) Update(ctx context.Context, id uuid.UUID, name string) (*domain.Asset, error) {
	a, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	a.Name = name
	if err := s.repo.Update(ctx, a); err != nil {
		return nil, err
	}
	return a, nil
}

func (s *assetService) Delete(ctx context.Context, id uuid.UUID) error {
	a, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}
	// Best-effort S3 cleanup for 3D assets (non-fatal).
	if a.AssetType == domain.AssetType3D && s.storer != nil {
		for _, key := range []string{a.Key, a.Thumbnail, a.Material} {
			if key != "" {
				_ = s.storer.DeleteObject(ctx, key)
			}
		}
	}
	return nil
}

func (s *assetService) Upload3D(ctx context.Context, in Upload3DAssetInput) (*domain.Asset, error) {
	assetID := uuid.New()
	if in.ID != "" {
		parsed, err := uuid.Parse(in.ID)
		if err != nil {
			return nil, domain.NewValidation(map[string]string{"id": "invalid UUID"})
		}
		assetID = parsed
	}

	existing, _ := s.repo.FindByID(ctx, assetID)
	if existing == nil && in.File == nil {
		return nil, domain.NewValidation(map[string]string{"file": "required when creating a new 3D asset"})
	}

	// Build S3 keys.
	var fileKey, thumbnailKey, materialKey string
	if in.File != nil {
		ext := sanitizeExt(filepath.Ext(in.File.Filename))
		fileKey = storage.Asset3DKey(s.env, assetID, "", ext)
	} else if existing != nil {
		fileKey = existing.Key
	}

	if fileKey == "" {
		return nil, domain.NewValidation(map[string]string{"file": "required when creating a new 3D asset"})
	}

	if in.Thumbnail != nil {
		ext := sanitizeExt(filepath.Ext(in.Thumbnail.Filename))
		thumbnailKey = storage.Asset3DKey(s.env, assetID, "_thumbnail", ext)
	} else if existing != nil {
		thumbnailKey = existing.Thumbnail
	}

	if in.Material != nil {
		ext := sanitizeExt(filepath.Ext(in.Material.Filename))
		materialKey = storage.Asset3DKey(s.env, assetID, "_material", ext)
	} else if existing != nil {
		materialKey = existing.Material
	}

	// Upload files to S3 first.
	if in.File != nil {
		if err := s.storer.PutMedia(ctx, fileKey, in.File.ContentType, in.File.Size, in.File.Reader); err != nil {
			return nil, fmt.Errorf("upload 3D file: %w", err)
		}
	}
	if in.Thumbnail != nil {
		if err := s.storer.PutMedia(ctx, thumbnailKey, in.Thumbnail.ContentType, in.Thumbnail.Size, in.Thumbnail.Reader); err != nil {
			return nil, fmt.Errorf("upload thumbnail: %w", err)
		}
	}
	if in.Material != nil {
		if err := s.storer.PutMedia(ctx, materialKey, in.Material.ContentType, in.Material.Size, in.Material.Reader); err != nil {
			return nil, fmt.Errorf("upload material: %w", err)
		}
	}

	// Merge with existing or set defaults.
	name := in.Name
	if name == "" && existing != nil {
		name = existing.Name
	}
	description := in.Description
	if description == "" && existing != nil {
		description = existing.Description
	}
	fileType := in.FileType
	if fileType == "" && existing != nil {
		fileType = existing.FileType
	}
	if fileType == "" && in.File != nil {
		fileType = mimeToFileType(in.File.ContentType)
	}
	width := in.Width
	if width == nil && existing != nil {
		width = existing.Width
	}
	height := in.Height
	if height == nil && existing != nil {
		height = existing.Height
	}

	var sizeBytes int64
	if in.File != nil {
		sizeBytes = in.File.Size
	} else if existing != nil {
		sizeBytes = existing.SizeBytes
	}

	a := &domain.Asset{
		ID:          assetID,
		VenueID:     in.VenueID,
		Name:        name,
		Key:         fileKey,
		ContentType: contentTypeOrEmpty(in.File),
		SizeBytes:   sizeBytes,
		URL:         s.AssetURL(fileKey),
		AssetType:   domain.AssetType3D,
		Description: description,
		FileType:    fileType,
		Thumbnail:   thumbnailKey,
		Material:    materialKey,
		Width:       width,
		Height:      height,
		CreatedBy:   in.CreatedBy,
	}

	if existing != nil {
		if err := s.repo.Update(ctx, a); err != nil {
			return nil, err
		}
	} else {
		if err := s.repo.Create(ctx, a); err != nil {
			return nil, err
		}
	}
	return a, nil
}

func (s *assetService) UploadLibrary(ctx context.Context, in UploadLibraryAssetInput) (*domain.Asset, error) {
	if in.File == nil {
		return nil, domain.NewValidation(map[string]string{"file": "required"})
	}
	id := uuid.New()
	key := storage.LibraryAssetKey(s.env, id, in.File.Filename)
	if key == "" {
		return nil, domain.NewValidation(map[string]string{"file": "invalid filename"})
	}
	if err := s.storer.PutMedia(ctx, key, in.File.ContentType, in.File.Size, in.File.Reader); err != nil {
		return nil, fmt.Errorf("upload library asset: %w", err)
	}
	status := in.Status
	if status == "" {
		status = domain.AssetStatusUnpublished
	}
	assetType := in.AssetType
	if assetType == "" {
		assetType = domain.AssetType2D
	}
	a := &domain.Asset{
		ID:        id,
		VenueID:   in.VenueID,
		Name:      in.File.Filename,
		Key:       key,
		URL:       s.AssetURL(key),
		AssetType: assetType,
		Status:    status,
		CreatedBy: in.CreatedBy,
	}
	if err := s.repo.Create(ctx, a); err != nil {
		return nil, err
	}
	return a, nil
}

func (s *assetService) UpdateLibrary(ctx context.Context, id uuid.UUID, in UploadLibraryAssetInput) (*domain.Asset, error) {
	a, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if in.Status != "" {
		a.Status = in.Status
	}
	if in.AssetType != "" {
		a.AssetType = in.AssetType
	}
	if in.File != nil {
		key := storage.LibraryAssetKey(s.env, id, in.File.Filename)
		if key == "" {
			return nil, domain.NewValidation(map[string]string{"file": "invalid filename"})
		}
		if err := s.storer.PutMedia(ctx, key, in.File.ContentType, in.File.Size, in.File.Reader); err != nil {
			return nil, fmt.Errorf("upload library asset: %w", err)
		}
		// Best-effort cleanup of old file.
		if a.Key != "" {
			_ = s.storer.DeleteObject(ctx, a.Key)
		}
		a.Key = key
		a.URL = s.AssetURL(key)
		a.Name = in.File.Filename
	}
	if err := s.repo.Update(ctx, a); err != nil {
		return nil, err
	}
	return a, nil
}

func sanitizeExt(ext string) string {
	ext = strings.TrimPrefix(strings.ToLower(ext), ".")
	var clean []byte
	for i := 0; i < len(ext) && i < 10; i++ {
		c := ext[i]
		if (c >= 'a' && c <= 'z') || (c >= '0' && c <= '9') {
			clean = append(clean, c)
		}
	}
	return string(clean)
}

func mimeToFileType(contentType string) string {
	mt, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		return contentType
	}
	parts := strings.SplitN(mt, "/", 2)
	if len(parts) == 2 {
		return parts[1]
	}
	return mt
}

func contentTypeOrEmpty(f *FileUpload) string {
	if f == nil {
		return ""
	}
	return f.ContentType
}
