package service

import (
	"context"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/hhung06/digimap-backend/internal/domain"
	"github.com/hhung06/digimap-backend/internal/platform/storage"
	"github.com/hhung06/digimap-backend/internal/repository"
)

// ── Product service ───────────────────────────────────────────────────────────

type ProductService interface {
	// Categories
	GetCategory(ctx context.Context, id uuid.UUID) (*domain.ProductCategory, error)
	ListCategories(ctx context.Context, venueID uuid.UUID) ([]*domain.ProductCategory, error)
	CreateCategory(ctx context.Context, c *domain.ProductCategory) error
	UpdateCategory(ctx context.Context, c *domain.ProductCategory) error
	DeleteCategory(ctx context.Context, id uuid.UUID) error

	// Products
	Get(ctx context.Context, id uuid.UUID) (*domain.Product, error)
	List(ctx context.Context, venueID uuid.UUID, p domain.Pagination) ([]*domain.Product, int64, error)
	SearchByName(ctx context.Context, venueID uuid.UUID, q string, limit int) ([]*domain.Product, error)
	Create(ctx context.Context, p *domain.Product, categoryIDs []uuid.UUID) error
	Update(ctx context.Context, p *domain.Product, categoryIDs []uuid.UUID) error
	Delete(ctx context.Context, id uuid.UUID) error

	// Attachments
	ListAttachments(ctx context.Context, productID uuid.UUID) ([]*domain.ProductAttachment, error)
	CreateAttachment(ctx context.Context, a *domain.ProductAttachment) error
	DeleteAttachment(ctx context.Context, productID, attachmentID uuid.UUID) error
}

type productService struct {
	repo repository.ProductRepository
}

func NewProductService(repo repository.ProductRepository) ProductService {
	return &productService{repo: repo}
}

func (s *productService) GetCategory(ctx context.Context, id uuid.UUID) (*domain.ProductCategory, error) {
	return s.repo.FindCategoryByID(ctx, id)
}

func (s *productService) ListCategories(ctx context.Context, venueID uuid.UUID) ([]*domain.ProductCategory, error) {
	return s.repo.ListCategories(ctx, venueID)
}

func (s *productService) CreateCategory(ctx context.Context, c *domain.ProductCategory) error {
	if c.Source == "" {
		c.Source = "internal"
	}
	return s.repo.CreateCategory(ctx, c)
}

func (s *productService) UpdateCategory(ctx context.Context, c *domain.ProductCategory) error {
	return s.repo.UpdateCategory(ctx, c)
}

func (s *productService) DeleteCategory(ctx context.Context, id uuid.UUID) error {
	return s.repo.DeleteCategory(ctx, id)
}

func (s *productService) Get(ctx context.Context, id uuid.UUID) (*domain.Product, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *productService) List(ctx context.Context, venueID uuid.UUID, p domain.Pagination) ([]*domain.Product, int64, error) {
	p.Normalize()
	return s.repo.List(ctx, venueID, p)
}

func (s *productService) SearchByName(ctx context.Context, venueID uuid.UUID, q string, limit int) ([]*domain.Product, error) {
	return s.repo.SearchByName(ctx, venueID, q, limit)
}

func (s *productService) Create(ctx context.Context, p *domain.Product, categoryIDs []uuid.UUID) error {
	if p.Source == "" {
		p.Source = "internal"
	}
	if err := s.repo.Create(ctx, p); err != nil {
		return err
	}
	if len(categoryIDs) > 0 {
		return s.repo.SetCategories(ctx, p.ID, categoryIDs)
	}
	return nil
}

func (s *productService) Update(ctx context.Context, p *domain.Product, categoryIDs []uuid.UUID) error {
	if err := s.repo.Update(ctx, p); err != nil {
		return err
	}
	return s.repo.SetCategories(ctx, p.ID, categoryIDs)
}

func (s *productService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}

func (s *productService) ListAttachments(ctx context.Context, productID uuid.UUID) ([]*domain.ProductAttachment, error) {
	return s.repo.ListAttachments(ctx, productID)
}

func (s *productService) CreateAttachment(ctx context.Context, a *domain.ProductAttachment) error {
	fileName := ""
	if a.File != nil {
		fileName = *a.File
	}
	a.FileType = detectFileType(fileName)
	return s.repo.CreateAttachment(ctx, a)
}

func (s *productService) DeleteAttachment(ctx context.Context, productID, attachmentID uuid.UUID) error {
	attachments, err := s.repo.ListAttachments(ctx, productID)
	if err != nil {
		return err
	}
	for _, att := range attachments {
		if att.ID == attachmentID {
			return s.repo.DeleteAttachment(ctx, attachmentID)
		}
	}
	return domain.NewNotFound("attachment not found")
}

// detectFileType mirrors Django's auto-detection based on file extension.
func detectFileType(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".jpg", ".jpeg", ".png", ".gif", ".bmp":
		return "image"
	case ".mp4", ".avi", ".mov", ".mkv":
		return "video"
	default:
		return "document"
	}
}

// ── Storage service ───────────────────────────────────────────────────────────

type StorageService interface {
	// PresignUpload returns a pre-signed URL the client can PUT a file to directly.
	PresignUpload(ctx context.Context, key, contentType string) (string, error)
}

type storageService struct {
	storer storage.Storer
}

func NewStorageService(storer storage.Storer) StorageService {
	return &storageService{storer: storer}
}

func (s *storageService) PresignUpload(ctx context.Context, key, contentType string) (string, error) {
	return s.storer.PresignUpload(ctx, key, contentType, 15*time.Minute)
}
