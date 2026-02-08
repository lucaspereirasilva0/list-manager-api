package service

import (
	"github.com/lucaspereirasilva0/list-manager-api/internal/domain"
	"github.com/lucaspereirasilva0/list-manager-api/internal/repository"
)

type parser struct{}

func (p parser) toRepositoryModel(item domain.Item) repository.Item {
	return repository.Item{
		ID:          item.ID,
		Name:        item.Name,
		Active:      item.Active,
		Observation: item.Observation,
		CreatedAt:   item.CreatedAt,
		UpdatedAt:   item.UpdatedAt,
	}
}

func (p parser) toDomainModel(item repository.Item) domain.Item {
	return domain.Item{
		ID:          item.ID,
		Name:        item.Name,
		Active:      item.Active,
		Observation: item.Observation,
		CreatedAt:   item.CreatedAt,
		UpdatedAt:   item.UpdatedAt,
	}
}
