package service

import (
	"context"

	"github.com/lucaspereirasilva0/list-manager-api/internal/domain"
)

type ItemService interface {
	CreateItem(ctx context.Context, item domain.Item) (domain.Item, error)
	GetItem(ctx context.Context, id string) (domain.Item, error)
	UpdateItem(ctx context.Context, item domain.Item) (domain.Item, error)
	DeleteItem(ctx context.Context, id string) error
	ListItems(ctx context.Context) ([]domain.Item, error)
}
