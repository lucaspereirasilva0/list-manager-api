package handlers

import "github.com/lucaspereirasilva0/list-manager-api/internal/domain"

type parser struct{}

func (p parser) toApiModel(item domain.Item) Item {
	return Item{
		ID:          item.ID,
		Name:        item.Name,
		Active:      item.Active,
		Observation: item.Observation,
		CreatedAt:   item.CreatedAt,
		UpdatedAt:   item.UpdatedAt,
	}
}

func (p parser) toDomainModel(item Item) domain.Item {
	return domain.Item{
		ID:          item.ID,
		Name:        item.Name,
		Active:      item.Active,
		Observation: item.Observation,
		CreatedAt:   item.CreatedAt,
		UpdatedAt:   item.UpdatedAt,
	}
}
