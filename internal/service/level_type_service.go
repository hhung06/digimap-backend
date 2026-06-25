package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/hhung06/digimap-backend/internal/domain"
	"github.com/hhung06/digimap-backend/internal/repository"
)

type LevelTypeService interface {
	List(ctx context.Context) ([]*domain.LevelType, error)
	Get(ctx context.Context, id uuid.UUID) (*domain.LevelType, error)
	Create(ctx context.Context, name, icon string, value int) (*domain.LevelType, error)
	Update(ctx context.Context, id uuid.UUID, name, icon string, value int) (*domain.LevelType, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type levelTypeService struct {
	repo repository.LevelTypeRepository
}

func NewLevelTypeService(repo repository.LevelTypeRepository) LevelTypeService {
	return &levelTypeService{repo: repo}
}

func (s *levelTypeService) List(ctx context.Context) ([]*domain.LevelType, error) {
	return s.repo.List(ctx)
}

func (s *levelTypeService) Get(ctx context.Context, id uuid.UUID) (*domain.LevelType, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *levelTypeService) Create(ctx context.Context, name, icon string, value int) (*domain.LevelType, error) {
	lt := &domain.LevelType{Name: name, Icon: icon, Value: value}
	if err := s.repo.Create(ctx, lt); err != nil {
		return nil, err
	}
	return lt, nil
}

func (s *levelTypeService) Update(ctx context.Context, id uuid.UUID, name, icon string, value int) (*domain.LevelType, error) {
	lt, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	lt.Name = name
	lt.Icon = icon
	lt.Value = value
	if err := s.repo.Update(ctx, lt); err != nil {
		return nil, err
	}
	return lt, nil
}

func (s *levelTypeService) Delete(ctx context.Context, id uuid.UUID) error {
	if _, err := s.repo.FindByID(ctx, id); err != nil {
		return err
	}
	return s.repo.Delete(ctx, id)
}
