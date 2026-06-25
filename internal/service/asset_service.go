package service

import (
	"context"
	"fmt"
	"io"
	"mime"
	"path/filepath"
	"strings"

	"github.com/google/uuid"

	"github.com/hhung06/digimap-backend/internal/domain"
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
	Create(ctx context.Context, in CreateAssetInput) (*domain.Asset, error)
	Update(ctx context.Context, id uuid.UUID, name string) (*domain.Asset, error)
	Delete(ctx context.Context, id uuid.UUID) error
	Upload3D(ctx context.Context, in Upload3DAssetInput) (*domain.Asset, error)
	UploadLibrary(ctx context.Context, in UploadLibraryAssetInput) (*domain.Asset, error)
	UpdateLibrary(ctx context.Context, id uuid.UUID, in UploadLibraryAssetInput) (*domain.Asset, error)
	// AssetURL returns the CloudFront URL for an S3 key, or "" if domain is not configured.
	AssetURL(key string) string
}

type assetService struct {
	repo         repository.AssetRepository
	storer       storage.Storer
	env          string
	cfDomain     string // AWS_CF_ASSETS_DOMAIN, e.g. "assets.example.com"
}

func NewAssetService(repo repository.AssetRepository, storer storage.Storer, env, cfDomain string) AssetService {
	return &assetService{repo: repo, storer: storer, env: env, cfDomain: cfDomain}
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
