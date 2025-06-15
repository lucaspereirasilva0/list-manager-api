package service

import (
	"context"

	"github.com/lucaspereirasilva0/list-manager-api/internal/domain"
	// "go.uber.org/zap"
)

type ItemService interface {
	CreateItem(ctx context.Context, name string, active bool) (domain.Item, error)
	GetItem(ctx context.Context, id string) (domain.Item, error)
	UpdateItem(ctx context.Context, item domain.Item) (domain.Item, error)
	DeleteItem(ctx context.Context, id string) error
	ListItems(ctx context.Context) ([]domain.Item, error)
}
