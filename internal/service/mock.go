package service

import (
	"context"

	"github.com/lucaspereirasilva0/list-manager-api/internal/domain"
	"github.com/stretchr/testify/mock"
)

type ItemServiceMock struct {
	mock.Mock
}

func (m *ItemServiceMock) CreateItem(ctx context.Context, name string, active bool) (domain.Item, error) {
	args := m.Called(ctx, name, active)
	return args.Get(0).(domain.Item), args.Error(1)
}

func (m *ItemServiceMock) GetItem(ctx context.Context, id string) (domain.Item, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(domain.Item), args.Error(1)
}

func (m *ItemServiceMock) UpdateItem(ctx context.Context, item domain.Item) (domain.Item, error) {
	args := m.Called(ctx, item)
	return args.Get(0).(domain.Item), args.Error(1)
}

func (m *ItemServiceMock) DeleteItem(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *ItemServiceMock) ListItems(ctx context.Context) ([]domain.Item, error) {
	args := m.Called(ctx)
	return args.Get(0).([]domain.Item), args.Error(1)
}
