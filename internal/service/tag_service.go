package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/hhung06/digimap-backend/internal/domain"
	"github.com/hhung06/digimap-backend/internal/repository"
)

type TagService interface {
	List(ctx context.Context, p domain.Pagination) ([]*domain.Tag, int, error)
	Get(ctx context.Context, id uuid.UUID) (*domain.Tag, error)
	Create(ctx context.Context, t *domain.Tag) error
	Update(ctx context.Context, t *domain.Tag) error
	Delete(ctx context.Context, id uuid.UUID) error
	AttachTag(ctx context.Context, et *domain.EntityTag) error
	DetachTag(ctx context.Context, tagID uuid.UUID, entityType string, entityID uuid.UUID) error
	ListEntityTags(ctx context.Context, entityType string, entityID uuid.UUID) ([]*domain.Tag, error)
}

type tagService struct {
	repo repository.TagRepository
}

func NewTagService(repo repository.TagRepository) TagService {
	return &tagService{repo: repo}
}

func (s *tagService) List(ctx context.Context, p domain.Pagination) ([]*domain.Tag, int, error) {
	return s.repo.List(ctx, p)
}

func (s *tagService) Get(ctx context.Context, id uuid.UUID) (*domain.Tag, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *tagService) Create(ctx context.Context, t *domain.Tag) error {
	return s.repo.Create(ctx, t)
}

func (s *tagService) Update(ctx context.Context, t *domain.Tag) error {
	return s.repo.Update(ctx, t)
}

func (s *tagService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}

func (s *tagService) AttachTag(ctx context.Context, et *domain.EntityTag) error {
	return s.repo.AttachTag(ctx, et)
}

func (s *tagService) DetachTag(ctx context.Context, tagID uuid.UUID, entityType string, entityID uuid.UUID) error {
	return s.repo.DetachTag(ctx, tagID, entityType, entityID)
}

func (s *tagService) ListEntityTags(ctx context.Context, entityType string, entityID uuid.UUID) ([]*domain.Tag, error) {
	return s.repo.ListEntityTags(ctx, entityType, entityID)
}
