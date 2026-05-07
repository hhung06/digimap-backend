package service_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/hhung06/digimap-backend/internal/repository/mocks"
	"github.com/hhung06/digimap-backend/internal/service"
)

func TestLevelService_DeleteMapGroup_DelegatesToRepository(t *testing.T) {
	repo := &mocks.LevelRepository{}
	svc := service.NewLevelService(repo)
	ctx := context.Background()
	id := uuid.New()

	repo.On("DeleteMapGroup", ctx, id).Return(nil)

	err := svc.DeleteMapGroup(ctx, id)

	require.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestLevelService_ListMapGroups_DelegatesToRepository(t *testing.T) {
	repo := &mocks.LevelRepository{}
	svc := service.NewLevelService(repo)
	ctx := context.Background()
	venueID := uuid.New()

	repo.On("ListMapGroups", ctx, venueID).Return(nil, nil)

	_, err := svc.ListMapGroups(ctx, venueID)

	require.NoError(t, err)
	repo.AssertExpectations(t)
}
