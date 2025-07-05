package mongodb

import (
	"context"

	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MockMongoCursorOperations is a mock for MongoCursorOperations.
type MockMongoCursorOperations struct {
	mock.Mock
}

func NewMockCollectionWrapper(mt *mtest.T) *mongoCollectionWrapper {
	return &mongoCollectionWrapper{
		collection: mt.Coll,
	}
}

func NewMockMongoClientWrappers(mt *mtest.T) *mongoClientWrapper {
	return &mongoClientWrapper{
		client: mt.Client,
	}
}

// All implements MongoCursorOperations.
func (m *MockMongoCursorOperations) All(ctx context.Context, results interface{}) error {
	args := m.Called(ctx, results)
	return args.Error(0)
}

// Close implements MongoCursorOperations.
func (m *MockMongoCursorOperations) Close(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

// MockMongoCollectionOperations is a mock for MongoCollectionOperations.
type MockMongoCollectionOperations struct {
	mock.Mock
}

// InsertOne implements MongoCollectionOperations.
func (m *MockMongoCollectionOperations) InsertOne(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
	args := m.Called(ctx, document)
	return args.Get(0).(*mongo.InsertOneResult), args.Error(1)
}

// FindOne implements MongoCollectionOperations.
func (m *MockMongoCollectionOperations) FindOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) *mongo.SingleResult {
	args := m.Called(ctx, filter)
	return args.Get(0).(*mongo.SingleResult)
}

// UpdateOne implements MongoCollectionOperations.
func (m *MockMongoCollectionOperations) UpdateOne(ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	args := m.Called(ctx, filter, update)
	return args.Get(0).(*mongo.UpdateResult), args.Error(1)
}

// DeleteOne implements MongoCollectionOperations.
func (m *MockMongoCollectionOperations) DeleteOne(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).(*mongo.DeleteResult), args.Error(1)
}

// Find implements MongoCollectionOperations.
func (m *MockMongoCollectionOperations) Find(ctx context.Context, filter interface{}, opts ...*options.FindOptions) (MongoCursorOperations, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).(MongoCursorOperations), args.Error(1)
}

// DeleteMany implements MongoCollectionOperations.
func (m *MockMongoCollectionOperations) DeleteMany(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	args := m.Called(ctx, filter, opts)
	return args.Get(0).(*mongo.DeleteResult), args.Error(1)
}

// MockMongoDatabaseOperations is a mock for MongoDatabaseOperations.
type MockMongoDatabaseOperations struct {
	mock.Mock
}

// Collection implements MongoDatabaseOperations.
func (m *MockMongoDatabaseOperations) Collection(name string, opts ...*options.CollectionOptions) MongoCollectionOperations {
	args := m.Called(name, opts)
	return args.Get(0).(MongoCollectionOperations)
}

// MockMongoClientOperations is a mock for MongoClientOperations.
type MockMongoClientOperations struct {
	mock.Mock
}

// Database implements MongoClientOperations.
func (m *MockMongoClientOperations) Database(name string, opts ...*options.DatabaseOptions) MongoDatabaseOperations {
	args := m.Called(name, opts)
	return args.Get(0).(MongoDatabaseOperations)
}

// MockClientOperations is a mock for ClientOperations.
type MockClientOperations struct {
	mock.Mock
}

// GetCollection implements ClientOperations.
func (m *MockClientOperations) GetCollection(collectionName string) MongoCollectionOperations {
	args := m.Called(collectionName)
	return args.Get(0).(MongoCollectionOperations)
}

// Disconnect implements ClientOperations.
func (m *MockClientOperations) Disconnect(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

// Client implements ClientOperations.
func (m *MockClientOperations) Client() MongoClientOperations {
	args := m.Called()
	return args.Get(0).(MongoClientOperations)
}
