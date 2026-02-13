package repository

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type RepositoryMock struct {
	mock.Mock
}

func (m *RepositoryMock) Create(ctx context.Context, item Item) (Item, error) {
	args := m.Called(ctx, item)
	return args.Get(0).(Item), args.Error(1)
}

func (m *RepositoryMock) Update(ctx context.Context, item Item) (Item, error) {
	args := m.Called(ctx, item)
	return args.Get(0).(Item), args.Error(1)
}

func (m *RepositoryMock) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *RepositoryMock) GetByID(ctx context.Context, id string) (Item, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(Item), args.Error(1)
}

func (m *RepositoryMock) List(ctx context.Context) ([]Item, error) {
	args := m.Called(ctx)
	return args.Get(0).([]Item), args.Error(1)
}

func (m *RepositoryMock) BulkUpdateActive(ctx context.Context, active bool) (int64, int64, error) {
	args := m.Called(ctx, active)
	return args.Get(0).(int64), args.Get(1).(int64), args.Error(2)
}
