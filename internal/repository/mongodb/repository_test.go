package mongodb_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/lucaspereirasilva0/list-manager-api/internal/repository"

	dbmongo "github.com/lucaspereirasilva0/list-manager-api/internal/database/mongodb"
	mongorepo "github.com/lucaspereirasilva0/list-manager-api/internal/repository/mongodb"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	errDatabase  = errors.New("database error")
	testObjectID = primitive.NewObjectID()
)

// --- Mock Data Functions (Parameter-less) ---

func mockCreateItemInput() repository.Item {
	return repository.Item{ID: testObjectID.Hex(), Name: "Test Item", Active: true}
}

func mockUpdateItemInput() repository.Item {
	return repository.Item{ID: testObjectID.Hex(), Name: "Updated Item", Active: false}
}

func mockInvalidHexIDItemInput() repository.Item {
	return repository.Item{ID: "invalid-hex-id", Name: "Test Item", Active: false}
}

func mockCreateItemOutput() repository.Item {
	return repository.Item{ID: testObjectID.Hex(), Name: "Test Item", Active: true, CreatedAt: time.Time{}, UpdatedAt: time.Time{}}
}

func mockFoundItemOutput() repository.Item {
	return repository.Item{ID: testObjectID.Hex(), Name: "Found Item", Active: true}
}

func mockUpdateItemOutput() repository.Item {
	return repository.Item{ID: testObjectID.Hex(), Name: "Updated Item", Active: false, UpdatedAt: time.Time{}}
}

func mockItemListOutput() []repository.Item {
	return []repository.Item{
		{ID: primitive.NewObjectID().Hex(), Name: "First Item", Active: true},
		{ID: primitive.NewObjectID().Hex(), Name: "Second Item", Active: false},
	}
}

func mockInsertOneResult() *mongo.InsertOneResult {
	return &mongo.InsertOneResult{InsertedID: testObjectID}
}

func mockSuccessfulUpdateOneResult() *mongo.UpdateResult {
	return &mongo.UpdateResult{MatchedCount: 1}
}

func mockNotFoundUpdateOneResult() *mongo.UpdateResult {
	return &mongo.UpdateResult{MatchedCount: 0}
}

func mockSuccessfulDeleteOneResult() *mongo.DeleteResult {
	return &mongo.DeleteResult{DeletedCount: 1}
}

func mockNotFoundDeleteOneResult() *mongo.DeleteResult {
	return &mongo.DeleteResult{DeletedCount: 0}
}

func mockSuccessfulUpdateManyResult() *mongo.UpdateResult {
	return &mongo.UpdateResult{MatchedCount: 5, ModifiedCount: 5}
}

func mockEmptyUpdateManyResult() *mongo.UpdateResult {
	return &mongo.UpdateResult{MatchedCount: 0, ModifiedCount: 0}
}

func mockPartialUpdateManyResult() *mongo.UpdateResult {
	return &mongo.UpdateResult{MatchedCount: 5, ModifiedCount: 3}
}

func mockSuccessfulFindOneResult() *mongo.SingleResult {
	item := mockFoundItemOutput()
	bsonBytes, _ := bson.Marshal(item)
	return mongo.NewSingleResultFromDocument(bsonBytes, nil, nil)
}

func mockNotFoundFindOneResult() *mongo.SingleResult {
	bsonBytes, _ := bson.Marshal(repository.Item{})
	return mongo.NewSingleResultFromDocument(bsonBytes, mongo.ErrNoDocuments, nil)
}

func mockDBErrorFindOneResult() *mongo.SingleResult {
	bsonBytes, _ := bson.Marshal(repository.Item{})
	return mongo.NewSingleResultFromDocument(bsonBytes, errDatabase, nil)
}

// --- Tests ---

func TestCreate(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name                     string
		givenItem                repository.Item
		givenMockInsertOneResult *mongo.InsertOneResult
		givenMockInsertOneError  error
		wantErr                  error
		wantCreatedItem          repository.Item
	}{
		{
			name:                     "Given_ValidItem_When_Create_Then_ExpectedSuccess",
			givenItem:                mockCreateItemInput(),
			givenMockInsertOneResult: mockInsertOneResult(),
			wantCreatedItem:          mockCreateItemOutput(),
		},
		{
			name:                    "Given_DatabaseError_When_Create_Then_ExpectedInternalError",
			givenItem:               mockCreateItemInput(),
			givenMockInsertOneError: errDatabase,
			wantErr:                 errDatabase,
		},
		{
			name:      "Given_ItemWithInvalidHexID_When_Create_Then_ExpectedInvalidIDError",
			givenItem: mockInvalidHexIDItemInput(),
			wantErr:   repository.NewInvalidHexIDError(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			collectionMock := new(dbmongo.MockMongoCollectionOperations)
			clientMock := new(dbmongo.MockClientOperations)

			if tt.givenMockInsertOneResult != nil || tt.givenMockInsertOneError != nil {
				collectionMock.On("InsertOne", ctx, mock.Anything).Return(tt.givenMockInsertOneResult, tt.givenMockInsertOneError)
			}

			clientMock.On("GetCollection", mongorepo.CollectionItems).Return(collectionMock)

			repo := mongorepo.NewMongoDBItemRepository(clientMock)

			createdItem, err := repo.Create(ctx, tt.givenItem)

			if tt.wantErr != nil {
				require.ErrorContains(t, err, tt.wantErr.Error())
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.wantCreatedItem, createdItem)
			}

			collectionMock.AssertExpectations(t)
			clientMock.AssertExpectations(t)
		})
	}
}

func TestGetByID(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name                   string
		givenID                string
		givenMockFindOneResult *mongo.SingleResult
		wantErr                error
		wantItem               repository.Item
	}{
		{
			name:                   "Given_ValidID_When_GetByID_Then_ExpectedSuccess",
			givenID:                testObjectID.Hex(),
			givenMockFindOneResult: mockSuccessfulFindOneResult(),
			wantItem:               mockFoundItemOutput(),
		},
		{
			name:                   "Given_ValidID_When_GetByID_And_ItemNotFound_Then_ExpectedNotFoundError",
			givenID:                testObjectID.Hex(),
			givenMockFindOneResult: mockNotFoundFindOneResult(),
			wantErr:                repository.NewItemNotFoundError(),
		},
		{
			name:                   "Given_ValidID_When_GetByID_And_DatabaseError_Then_ExpectedInternalError",
			givenID:                testObjectID.Hex(),
			givenMockFindOneResult: mockDBErrorFindOneResult(),
			wantErr:                errDatabase,
		},
		{
			name:    "Given_InvalidID_When_GetByID_Then_ExpectedInvalidIDError",
			givenID: "invalid-id",
			wantErr: repository.NewInvalidHexIDError(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			collectionMock := new(dbmongo.MockMongoCollectionOperations)
			clientMock := new(dbmongo.MockClientOperations)

			if tt.givenMockFindOneResult != nil {
				collectionMock.On("FindOne", ctx, mock.Anything).Return(tt.givenMockFindOneResult)
			}

			clientMock.On("GetCollection", mongorepo.CollectionItems).Return(collectionMock)

			repo := mongorepo.NewMongoDBItemRepository(clientMock)

			item, err := repo.GetByID(ctx, tt.givenID)

			if tt.wantErr != nil {
				require.ErrorContains(t, err, tt.wantErr.Error())
			} else {
				require.NoError(t, err)
			}
			require.Equal(t, tt.wantItem, item)

			collectionMock.AssertExpectations(t)
			clientMock.AssertExpectations(t)
		})
	}
}

func TestUpdate(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name                     string
		givenItem                repository.Item
		givenMockUpdateOneResult *mongo.UpdateResult
		givenMockUpdateOneError  error
		wantErr                  error
		wantUpdatedItem          repository.Item
	}{
		{
			name:                     "Given_ValidItem_When_Update_Then_ExpectedSuccess",
			givenItem:                mockUpdateItemInput(),
			givenMockUpdateOneResult: mockSuccessfulUpdateOneResult(),
			wantUpdatedItem:          mockUpdateItemOutput(),
		},
		{
			name:                     "Given_ValidItem_When_Update_And_ItemNotFound_Then_ExpectedNotFoundError",
			givenItem:                mockUpdateItemInput(),
			givenMockUpdateOneResult: mockNotFoundUpdateOneResult(),
			wantErr:                  repository.NewItemNotFoundError(),
		},
		{
			name:                    "Given_ValidItem_When_Update_And_DatabaseError_Then_ExpectedInternalError",
			givenItem:               mockUpdateItemInput(),
			givenMockUpdateOneError: errDatabase,
			wantErr:                 errDatabase,
		},
		{
			name:      "Given_InvalidID_When_Update_Then_ExpectedInvalidIDError",
			givenItem: mockInvalidHexIDItemInput(),
			wantErr:   repository.NewInvalidHexIDError(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			collectionMock := new(dbmongo.MockMongoCollectionOperations)
			clientMock := new(dbmongo.MockClientOperations)

			if tt.givenMockUpdateOneResult != nil || tt.givenMockUpdateOneError != nil {
				collectionMock.On("UpdateOne", ctx, mock.Anything, mock.Anything).Return(tt.givenMockUpdateOneResult, tt.givenMockUpdateOneError)
			}

			clientMock.On("GetCollection", mongorepo.CollectionItems).Return(collectionMock)

			repo := mongorepo.NewMongoDBItemRepository(clientMock)

			updatedItem, err := repo.Update(ctx, tt.givenItem)

			if tt.wantErr != nil {
				require.ErrorContains(t, err, tt.wantErr.Error())
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.wantUpdatedItem, updatedItem)
			}

			collectionMock.AssertExpectations(t)
			clientMock.AssertExpectations(t)
		})
	}
}

func TestDelete(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name                     string
		givenID                  string
		givenMockDeleteOneResult *mongo.DeleteResult
		givenMockDeleteOneError  error
		wantErr                  error
	}{
		{
			name:                     "Given_ValidID_When_Delete_Then_ExpectedSuccess",
			givenID:                  testObjectID.Hex(),
			givenMockDeleteOneResult: mockSuccessfulDeleteOneResult(),
		},
		{
			name:                     "Given_ValidID_When_Delete_And_ItemNotFound_Then_ExpectedNotFoundError",
			givenID:                  testObjectID.Hex(),
			givenMockDeleteOneResult: mockNotFoundDeleteOneResult(),
			wantErr:                  repository.NewItemNotFoundError(),
		},
		{
			name:                    "Given_ValidID_When_Delete_And_DatabaseError_Then_ExpectedInternalError",
			givenID:                 testObjectID.Hex(),
			givenMockDeleteOneError: errDatabase,
			wantErr:                 errDatabase,
		},
		{
			name:    "Given_InvalidID_When_Delete_Then_ExpectedInvalidIDError",
			givenID: "invalid-id",
			wantErr: repository.NewInvalidHexIDError(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			collectionMock := new(dbmongo.MockMongoCollectionOperations)
			clientMock := new(dbmongo.MockClientOperations)

			if tt.givenMockDeleteOneResult != nil || tt.givenMockDeleteOneError != nil {
				collectionMock.On("DeleteOne", ctx, mock.Anything).Return(tt.givenMockDeleteOneResult, tt.givenMockDeleteOneError)
			}

			clientMock.On("GetCollection", mongorepo.CollectionItems).Return(collectionMock)

			repo := mongorepo.NewMongoDBItemRepository(clientMock)

			err := repo.Delete(ctx, tt.givenID)

			if tt.wantErr != nil {
				require.ErrorContains(t, err, tt.wantErr.Error())
			} else {
				require.NoError(t, err)
			}

			collectionMock.AssertExpectations(t)
			clientMock.AssertExpectations(t)
		})
	}
}

func TestBulkUpdateActive(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name                      string
		givenActive               bool
		givenMockUpdateManyResult *mongo.UpdateResult
		givenMockUpdateManyError  error
		wantErr                   error
		wantMatchedCount          int64
		wantModifiedCount         int64
	}{
		{
			name:                      "Given_TrueActive_When_BulkUpdateActive_Then_ReturnsSuccess",
			givenActive:               true,
			givenMockUpdateManyResult: mockSuccessfulUpdateManyResult(),
			wantMatchedCount:          5,
			wantModifiedCount:         5,
		},
		{
			name:                      "Given_FalseActive_When_BulkUpdateActive_Then_ReturnsSuccess",
			givenActive:               false,
			givenMockUpdateManyResult: mockSuccessfulUpdateManyResult(),
			wantMatchedCount:          5,
			wantModifiedCount:         5,
		},
		{
			name:                      "Given_EmptyCollection_When_BulkUpdateActive_Then_ReturnsZeroCounts",
			givenActive:               true,
			givenMockUpdateManyResult: mockEmptyUpdateManyResult(),
			wantMatchedCount:          0,
			wantModifiedCount:         0,
		},
		{
			name:                      "Given_PartialUpdate_When_BulkUpdateActive_Then_ReturnsPartialCounts",
			givenActive:               true,
			givenMockUpdateManyResult: mockPartialUpdateManyResult(),
			wantMatchedCount:          5,
			wantModifiedCount:         3,
		},
		{
			name:                     "Given_DatabaseError_When_BulkUpdateActive_Then_ReturnsError",
			givenActive:              true,
			givenMockUpdateManyError: errDatabase,
			wantErr:                  errDatabase,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			collectionMock := new(dbmongo.MockMongoCollectionOperations)
			clientMock := new(dbmongo.MockClientOperations)

			if tt.givenMockUpdateManyResult != nil || tt.givenMockUpdateManyError != nil {
				collectionMock.On("UpdateMany", ctx, mock.Anything, mock.Anything, mock.Anything).Return(tt.givenMockUpdateManyResult, tt.givenMockUpdateManyError)
			}

			clientMock.On("GetCollection", mongorepo.CollectionItems).Return(collectionMock)

			repo := mongorepo.NewMongoDBItemRepository(clientMock)

			matchedCount, modifiedCount, err := repo.BulkUpdateActive(ctx, tt.givenActive)

			if tt.wantErr != nil {
				require.ErrorContains(t, err, tt.wantErr.Error())
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.wantMatchedCount, matchedCount)
				require.Equal(t, tt.wantModifiedCount, modifiedCount)
			}

			collectionMock.AssertExpectations(t)
			clientMock.AssertExpectations(t)
		})
	}
}

func TestList(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name           string
		givenFindError error
		givenAllError  error
		wantItems      []repository.Item
		wantErr        error
	}{
		{
			name:      "Given_ItemsExist_When_List_Then_ExpectedSuccess",
			wantItems: mockItemListOutput(),
		},
		{
			name:      "Given_NoItemsExist_When_List_Then_ExpectedEmptyList",
			wantItems: []repository.Item{},
		},
		{
			name:           "Given_FindError_When_List_Then_ExpectedInternalError",
			givenFindError: errDatabase,
			wantErr:        errDatabase,
		},
		{
			name:          "Given_CursorAllError_When_List_Then_ExpectedInternalError",
			givenAllError: errDatabase,
			wantErr:       errDatabase,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			collectionMock := new(dbmongo.MockMongoCollectionOperations)
			cursorMock := new(dbmongo.MockMongoCursorOperations)
			clientMock := new(dbmongo.MockClientOperations)

			if tt.givenFindError != nil {
				collectionMock.On("Find", ctx, mock.Anything).Return((*dbmongo.MockMongoCursorOperations)(nil), tt.givenFindError)
			} else {
				collectionMock.On("Find", ctx, mock.Anything).Return(cursorMock, nil)
				cursorMock.On("All", ctx, mock.Anything).Return(tt.givenAllError).Run(func(args mock.Arguments) {
					if tt.givenAllError == nil {
						results := args.Get(1).(*[]repository.Item)
						*results = tt.wantItems
					}
				})
				cursorMock.On("Close", ctx).Return(nil)
			}

			clientMock.On("GetCollection", mongorepo.CollectionItems).Return(collectionMock)

			repo := mongorepo.NewMongoDBItemRepository(clientMock)

			items, err := repo.List(ctx)

			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
			} else {
				require.NoError(t, err)
			}
			require.Equal(t, tt.wantItems, items)

			collectionMock.AssertExpectations(t)
			cursorMock.AssertExpectations(t)
		})
	}
}
