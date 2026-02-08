package local

//TODO this file is used only for save in memory database
// import (
// 	"context"

// 	"github.com/lucaspereirasilva0/list-manager-api/internal/repository"
// )

// type localRepository struct {
// 	items map[string]repository.Item
// }

// func NewLocalRepository() repository.ItemRepository {
// 	return &localRepository{
// 		items: make(map[string]repository.Item),
// 	}
// }

// func (r *localRepository) Create(ctx context.Context, item repository.Item) (repository.Item, error) {
// 	r.items[item.ID] = item
// 	return item, nil
// }

// func (r *localRepository) Update(ctx context.Context, item repository.Item) (repository.Item, error) {
// 	r.items[item.ID] = item
// 	return item, nil
// }

// func (r *localRepository) Delete(ctx context.Context, id string) error {
// 	delete(r.items, id)
// 	return nil
// }

// func (r *localRepository) GetByID(ctx context.Context, id string) (repository.Item, error) {
// 	item, ok := r.items[id]
// 	if !ok {
// 		return repository.Item{}, repository.NewRepositoryError(repository.ErrItemNotFound)
// 	}
// 	return item, nil
// }

// func (r *localRepository) List(ctx context.Context) ([]repository.Item, error) {
// 	items := make([]repository.Item, 0, len(r.items))
// 	for _, item := range r.items {
// 		items = append(items, item)
// 	}
// 	return items, nil
// }
