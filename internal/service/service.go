package service

import (
	"context"
	"log"

	"github.com/lucaspereirasilva0/list-manager-api/internal/domain"
	"github.com/lucaspereirasilva0/list-manager-api/internal/repository"
)

type itemService struct {
	repository repository.ItemRepository
	parser     parser
}

func NewItemService(repository repository.ItemRepository) ItemService {
	return &itemService{
		repository: repository,
		parser:     parser{},
	}
}

func (s *itemService) CreateItem(ctx context.Context, item domain.Item) (domain.Item, error) {
	repositoryItem := s.parser.toRepositoryModel(domain.NewItem(item.Name, item.Active))

	createdRepositoryItem, err := s.repository.Create(ctx, repositoryItem)
	if err != nil {
		log.Printf("failed to create item: %s: %w"+item.Name, err)
		return domain.Item{}, handleError(err)
	}

	return s.parser.toDomainModel(createdRepositoryItem), nil
}

func (s *itemService) UpdateItem(ctx context.Context, item domain.Item) (domain.Item, error) {
	if item.IsEmpty() {
		return domain.Item{}, NewErrorEmptyItem()
	}
	_, err := s.repository.GetByID(ctx, item.ID)
	if err != nil {
		log.Printf("failed to get item: %s: %w"+item.ID, err)
		return domain.Item{}, handleError(err)
	}

	repositoryItem := s.parser.toRepositoryModel(item)

	updatedItem, err := s.repository.Update(ctx, repositoryItem)
	if err != nil {
		log.Printf("failed to update item: %s: %w"+item.ID, err)
		return domain.Item{}, handleError(err)
	}

	return s.parser.toDomainModel(updatedItem), nil
}

func (s *itemService) GetItem(ctx context.Context, id string) (domain.Item, error) {
	item, err := s.repository.GetByID(ctx, id)
	if err != nil {
		log.Printf("failed to get item: %s: %w"+id, err)
		return domain.Item{}, handleError(err)
	}

	return s.parser.toDomainModel(item), nil
}

func (s *itemService) DeleteItem(ctx context.Context, id string) error {
	err := s.repository.Delete(ctx, id)
	if err != nil {
		log.Printf("failed to delete item: %s: %w"+id, err)
		return handleError(err)
	}
	return nil
}

func (s *itemService) ListItems(ctx context.Context) ([]domain.Item, error) {
	items, err := s.repository.List(ctx)
	if err != nil {
		log.Printf("failed to list items: %v", err)
		return nil, handleError(err)
	}

	domainItems := make([]domain.Item, len(items))
	for i, item := range items {
		domainItems[i] = s.parser.toDomainModel(item)
	}

	return domainItems, nil
}
