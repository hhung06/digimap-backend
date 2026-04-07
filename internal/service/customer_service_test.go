package service_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/hhung06/digimap-backend/internal/domain"
	"github.com/hhung06/digimap-backend/internal/repository/mocks"
	"github.com/hhung06/digimap-backend/internal/service"
)

func TestCustomerService_Get(t *testing.T) {
	repo := &mocks.CustomerRepository{}
	svc := service.NewCustomerService(repo)

	ctx := context.Background()
	id := uuid.New()
	expected := &domain.Customer{ID: id, Name: "Acme Corp"}

	repo.On("FindByID", ctx, id).Return(expected, nil)

	c, err := svc.Get(ctx, id)
	require.NoError(t, err)
	assert.Equal(t, "Acme Corp", c.Name)
	repo.AssertExpectations(t)
}

func TestCustomerService_Get_NotFound(t *testing.T) {
	repo := &mocks.CustomerRepository{}
	svc := service.NewCustomerService(repo)

	ctx := context.Background()
	id := uuid.New()

	repo.On("FindByID", ctx, id).Return((*domain.Customer)(nil), domain.NewNotFound("customer not found"))

	_, err := svc.Get(ctx, id)
	require.Error(t, err)
	assert.ErrorIs(t, err, domain.ErrNotFound)
}

func TestCustomerService_List(t *testing.T) {
	repo := &mocks.CustomerRepository{}
	svc := service.NewCustomerService(repo)

	ctx := context.Background()
	p := domain.Pagination{Page: 1, PageSize: 10}
	expected := []*domain.Customer{{Name: "Acme"}, {Name: "Globex"}}

	// List normalizes pagination before calling repo — use mockAny to match any pagination
	repo.On("List", ctx, mockAny).Return(expected, 2, nil)

	customers, total, err := svc.List(ctx, p)
	require.NoError(t, err)
	assert.Len(t, customers, 2)
	assert.Equal(t, int64(2), total)
	repo.AssertExpectations(t)
}

func TestCustomerService_Create(t *testing.T) {
	repo := &mocks.CustomerRepository{}
	svc := service.NewCustomerService(repo)

	ctx := context.Background()
	c := &domain.Customer{Name: "New Corp"}

	repo.On("Create", ctx, c).Return(nil)

	err := svc.Create(ctx, c)
	require.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestCustomerService_Update(t *testing.T) {
	repo := &mocks.CustomerRepository{}
	svc := service.NewCustomerService(repo)

	ctx := context.Background()
	c := &domain.Customer{ID: uuid.New(), Name: "Updated Corp"}

	repo.On("Update", ctx, c).Return(nil)

	err := svc.Update(ctx, c)
	require.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestCustomerService_Delete(t *testing.T) {
	repo := &mocks.CustomerRepository{}
	svc := service.NewCustomerService(repo)

	ctx := context.Background()
	id := uuid.New()

	repo.On("Delete", ctx, id).Return(nil)

	err := svc.Delete(ctx, id)
	require.NoError(t, err)
	repo.AssertExpectations(t)
}
