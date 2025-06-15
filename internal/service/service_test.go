package service_test

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/lucaspereirasilva0/list-manager-api/internal/domain"
	"github.com/lucaspereirasilva0/list-manager-api/internal/repository"
	"github.com/lucaspereirasilva0/list-manager-api/internal/service"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

const (
	_dummyID = "123"
)

var (
	errDummy = errors.New("dummy error")
)

func TestCreateItem(t *testing.T) {
	tests := []struct {
		name                string
		givenItem           domain.Item
		givenRepositoryItem repository.Item
		wantServiceItem     domain.Item
		wantErr             error
	}{
		{
			name:                "Given_Item_When_CreateItem_Then_ExpectedSuccess",
			givenItem:           mockServiceItem(),
			givenRepositoryItem: mockOutputRepositoryItem(),
			wantServiceItem:     mockServiceItem(),
		},
		{
			name:                "Given_Item_When_CreateItem_Then_ExpectedInternalError",
			givenItem:           mockServiceItem(),
			givenRepositoryItem: mockOutputRepositoryItem(),
			wantErr:             mockInternalServerError(repository.NewRepositoryError(errDummy)),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			mockRepo := &repository.RepositoryMock{}
			mockRepo.On("Create", ctx, mock.MatchedBy(validateRepositoryItem(tt.givenRepositoryItem))).Return(tt.givenRepositoryItem, tt.wantErr)

			service := service.NewItemService(mockRepo)
			item, err := service.CreateItem(ctx, tt.givenItem.Name, tt.givenItem.Active)

			if tt.wantErr != nil {
				require.ErrorContains(t, err, tt.wantErr.Error())
			} else {
				require.EqualValues(t, tt.wantServiceItem, item)
				require.NoError(t, err)
			}
		})
	}
}

func TestGetItem(t *testing.T) {
	tests := []struct {
		name                string
		givenID             string
		wantItem            domain.Item
		givenRepositoryItem repository.Item
		wantErr             error
	}{
		{
			name:                "Given_Item_When_GetItem_Then_ExpectedSuccess",
			givenID:             _dummyID,
			givenRepositoryItem: mockOutputRepositoryItem(),
			wantItem:            mockServiceItem(),
		},
		{
			name:    "Given_Item_When_GetItem_Then_ExpectedErrFailedGetItem",
			givenID: _dummyID,
			wantErr: mockNotFoundRepositoryError(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			mockRepo := &repository.RepositoryMock{}
			mockRepo.On("GetByID", ctx, tt.givenID).Return(tt.givenRepositoryItem, tt.wantErr)

			service := service.NewItemService(mockRepo)
			item, err := service.GetItem(ctx, tt.givenID)

			require.Equal(t, tt.wantItem, item)
			if tt.wantErr != nil {
				require.ErrorContains(t, err, tt.wantErr.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestUpdateItem(t *testing.T) {
	type mockUpdate struct {
		givenOutputItem repository.Item
		givenErr        error
	}

	type mockGetByID struct {
		givenItem repository.Item
		givenErr  error
	}

	type givenMockRepository struct {
		mockUpdate
		mockGetByID
	}

	tests := []struct {
		name      string
		givenItem domain.Item
		givenMockRepository
		wantServiceItem domain.Item
		wantErr         error
	}{
		{
			name:      "Given_ValidItem_When_UpdateItem_Then_ExpectedSuccess",
			givenItem: domain.Item{ID: _dummyID, Name: "updated-name", Active: false},
			givenMockRepository: givenMockRepository{
				mockUpdate: mockUpdate{
					givenOutputItem: mockOutputRepositoryItem(),
				},
				mockGetByID: mockGetByID{
					givenItem: mockRepositoryItem(),
				},
			},
			wantServiceItem: mockServiceItem(),
		},
		{
			name:      "Given_ItemNotFound_When_UpdateItem_Then_ExpectedNotFoundError",
			givenItem: domain.Item{ID: _dummyID, Name: "updated-name", Active: false},
			givenMockRepository: givenMockRepository{
				mockGetByID: mockGetByID{
					givenErr: mockNotFoundRepositoryError(),
				},
			},
			wantErr: mockNotFoundRepositoryError(),
		},
		{
			name:      "Given_InternalError_When_UpdateItem_Then_ExpectedInternalServerError",
			givenItem: domain.Item{ID: _dummyID, Name: "updated-name", Active: false},
			givenMockRepository: givenMockRepository{
				mockUpdate: mockUpdate{
					givenOutputItem: mockOutputRepositoryItem(),
					givenErr:        repository.NewRepositoryError(errDummy),
				},
			},
			wantErr: mockInternalServerError(repository.NewRepositoryError(errDummy)),
		},
		{
			name:      "Given_EmptyItem_When_UpdateItem_Then_ExpectedEmptyItemError",
			givenItem: domain.Item{ID: "", Name: "any-name", Active: true},
			wantErr:   service.NewErrorEmptyItem(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			mockRepo := &repository.RepositoryMock{}
			mockRepo.On("GetByID", ctx, mock.AnythingOfType("string")).
				Return(tt.givenMockRepository.mockGetByID.givenItem, tt.givenMockRepository.mockGetByID.givenErr)
			mockRepo.On("Update", ctx, mock.AnythingOfType("repository.Item")).
				Return(tt.givenMockRepository.mockUpdate.givenOutputItem, tt.givenMockRepository.mockUpdate.givenErr)

			service := service.NewItemService(mockRepo)
			item, err := service.UpdateItem(ctx, tt.givenItem)

			if tt.wantErr != nil {
				require.ErrorContains(t, err, tt.wantErr.Error())
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.wantServiceItem, item)
			}
		})
	}
}

func TestDeleteItem(t *testing.T) {
	tests := []struct {
		name    string
		givenID string
		wantErr error
	}{
		{
			name:    "Given_Item_When_DeleteItem_Then_ExpectedSuccess",
			givenID: _dummyID,
		},
		{
			name:    "Given_Item_When_DeleteItem_Then_ExpectedErrFailedDeleteItem",
			givenID: _dummyID,
			wantErr: mockInternalServerError(repository.NewRepositoryError(errDummy)),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			mockRepo := &repository.RepositoryMock{}
			mockRepo.On("Delete", ctx, tt.givenID).Return(tt.wantErr)

			service := service.NewItemService(mockRepo)
			err := service.DeleteItem(ctx, tt.givenID)

			if tt.wantErr != nil {
				require.ErrorContains(t, err, tt.wantErr.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestListItems(t *testing.T) {
	tests := []struct {
		name                 string
		givenRepositoryItems []repository.Item
		wantServiceItems     []domain.Item
		wantErr              error
	}{
		{
			name:                 "Given_Items_When_ListItems_Then_ExpectedSuccess",
			givenRepositoryItems: []repository.Item{mockOutputRepositoryItem()},
			wantServiceItems:     []domain.Item{mockServiceItem()},
		},
		{
			name:                 "Given_NoItems_When_ListItems_Then_ExpectedEmptyList",
			givenRepositoryItems: []repository.Item{},
			wantServiceItems:     []domain.Item{},
		},
		{
			name:                 "Given_Error_When_ListItems_Then_ExpectedInternalError",
			givenRepositoryItems: []repository.Item{},
			wantServiceItems:     []domain.Item{},
			wantErr:              mockInternalServerError(repository.NewRepositoryError(errDummy)),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			mockRepo := &repository.RepositoryMock{}
			mockRepo.On("List", ctx).Return(tt.givenRepositoryItems, tt.wantErr)

			service := service.NewItemService(mockRepo)
			items, err := service.ListItems(ctx)

			if tt.wantErr != nil {
				require.ErrorContains(t, err, tt.wantErr.Error())
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.wantServiceItems, items)
			}
		})
	}
}

func mockOutputRepositoryItem() repository.Item {
	return repository.Item{
		ID:     _dummyID,
		Name:   "updated-name",
		Active: true,
	}
}

func mockRepositoryItem() repository.Item {
	return repository.Item{
		ID:     _dummyID,
		Name:   "test",
		Active: true,
	}
}

func mockServiceItem() domain.Item {
	return domain.Item{
		ID:     _dummyID,
		Name:   "updated-name",
		Active: true,
	}
}

func mockNotFoundRepositoryError() error {
	return service.NewErrorService(
		repository.ErrItemNotFound,
		"item not found",
		http.StatusNotFound,
	)
}

func mockInternalServerError(err error) error {
	return service.NewErrorService(
		err,
		"internal server error",
		http.StatusInternalServerError,
	)
}

func validateRepositoryItem(expected repository.Item) func(item repository.Item) bool {
	return func(actual repository.Item) bool {
		return actual.Name == expected.Name && actual.Active == expected.Active
	}
}
