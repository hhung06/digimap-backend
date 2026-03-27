package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/hhung06/digimap-backend/internal/domain"
	"github.com/hhung06/digimap-backend/internal/repository"
)

type ConnectionService interface {
	List(ctx context.Context, venueID uuid.UUID, p domain.Pagination) ([]*domain.Connection, int, error)
	Get(ctx context.Context, id uuid.UUID) (*domain.Connection, error)
	Create(ctx context.Context, c *domain.Connection) error
	Update(ctx context.Context, c *domain.Connection) error
	Delete(ctx context.Context, id uuid.UUID) error
	ListLevels(ctx context.Context, connectionID uuid.UUID) ([]*domain.ConnectionLevel, error)
	AddLevel(ctx context.Context, cl *domain.ConnectionLevel) error
	RemoveLevel(ctx context.Context, id uuid.UUID) error
}

type connectionService struct {
	repo repository.ConnectionRepository
}

func NewConnectionService(repo repository.ConnectionRepository) ConnectionService {
	return &connectionService{repo: repo}
}

func (s *connectionService) List(ctx context.Context, venueID uuid.UUID, p domain.Pagination) ([]*domain.Connection, int, error) {
	return s.repo.List(ctx, venueID, p)
}

func (s *connectionService) Get(ctx context.Context, id uuid.UUID) (*domain.Connection, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *connectionService) Create(ctx context.Context, c *domain.Connection) error {
	return s.repo.Create(ctx, c)
}

func (s *connectionService) Update(ctx context.Context, c *domain.Connection) error {
	return s.repo.Update(ctx, c)
}

func (s *connectionService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}

func (s *connectionService) ListLevels(ctx context.Context, connectionID uuid.UUID) ([]*domain.ConnectionLevel, error) {
	return s.repo.ListLevels(ctx, connectionID)
}

func (s *connectionService) AddLevel(ctx context.Context, cl *domain.ConnectionLevel) error {
	return s.repo.AddLevel(ctx, cl)
}

func (s *connectionService) RemoveLevel(ctx context.Context, id uuid.UUID) error {
	return s.repo.RemoveLevel(ctx, id)
}
