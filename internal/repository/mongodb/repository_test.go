package mongodb_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	dbmongo "github.com/lucaspereirasilva0/list-manager-api/internal/database/mongodb"
	"github.com/lucaspereirasilva0/list-manager-api/internal/repository"
	mongorepo "github.com/lucaspereirasilva0/list-manager-api/internal/repository/mongodb"
)

// -------------------- hand-written mocks --------------------

type mockCursor struct {
	allFunc   func(ctx context.Context, results interface{}) error
	closeFunc func(ctx context.Context) error
}

func (m *mockCursor) All(ctx context.Context, results interface{}) error {
	return m.allFunc(ctx, results)
}
func (m *mockCursor) Close(ctx context.Context) error {
	if m.closeFunc != nil {
		return m.closeFunc(ctx)
	}
	return nil
}

type mockCollection struct {
	insertOneFunc func(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error)
	findOneFunc   func(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) *mongo.SingleResult
	updateOneFunc func(ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error)
	deleteOneFunc func(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error)
	findFunc      func(ctx context.Context, filter interface{}, opts ...*options.FindOptions) (dbmongo.MongoCursorOperations, error)
}

func (m *mockCollection) InsertOne(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
	return m.insertOneFunc(ctx, document, opts...)
}
func (m *mockCollection) FindOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) *mongo.SingleResult {
	return m.findOneFunc(ctx, filter, opts...)
}
func (m *mockCollection) UpdateOne(ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	return m.updateOneFunc(ctx, filter, update, opts...)
}
func (m *mockCollection) DeleteOne(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	return m.deleteOneFunc(ctx, filter, opts...)
}
func (m *mockCollection) Find(ctx context.Context, filter interface{}, opts ...*options.FindOptions) (dbmongo.MongoCursorOperations, error) {
	return m.findFunc(ctx, filter, opts...)
}
func (m *mockCollection) DeleteMany(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	return nil, nil
}

// mock MongoClientOperations
type mockMongoClientOps struct {
	err error
}

func (mc *mockMongoClientOps) StartSession(opts ...*options.SessionOptions) (mongo.Session, error) {
	return nil, mc.err
}
func (mc *mockMongoClientOps) Database(name string, opts ...*options.DatabaseOptions) dbmongo.MongoDatabaseOperations {
	return nil
}

type mockClient struct {
	itemsCol dbmongo.MongoCollectionOperations
	usersCol dbmongo.MongoCollectionOperations
	mcOps    dbmongo.MongoClientOperations
}

func (c *mockClient) GetCollection(name string) dbmongo.MongoCollectionOperations {
	switch name {
	case mongorepo.CollectionItems:
		return c.itemsCol
	case mongorepo.CollectionUsers:
		if c.usersCol != nil {
			return c.usersCol
		}
	}
	return nil
}
func (c *mockClient) Disconnect(ctx context.Context) error  { return nil }
func (c *mockClient) Client() dbmongo.MongoClientOperations { return c.mcOps }

// -------------------- test suite --------------------

type MongoSuite struct {
	suite.Suite
	ctx        context.Context
	itemsCol   *mockCollection
	client     *mockClient
	repository repository.ItemRepository
}

func (s *MongoSuite) SetupTest() {
	s.ctx = context.Background()
	s.itemsCol = &mockCollection{}
	s.client = &mockClient{itemsCol: s.itemsCol}
	s.repository = mongorepo.NewMongoDBItemRepository(s.client)
}

// ---- Create ----
func (s *MongoSuite) TestCreate() {
	s.itemsCol.insertOneFunc = func(ctx context.Context, doc interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
		return &mongo.InsertOneResult{InsertedID: primitive.NewObjectID()}, nil
	}

	item := repository.Item{ID: primitive.NewObjectID().Hex(), Name: "Item 1", Active: true}
	created, err := s.repository.Create(s.ctx, item)

	s.NoError(err)
	s.Equal(item.Name, created.Name)
	s.Equal(item.ID, created.ID)
}

// ---- GetByID ----
func (s *MongoSuite) TestGetByID() {
	id := primitive.NewObjectID().Hex()
	expected := repository.Item{ID: id, Name: "Found", Active: true}

	s.itemsCol.findOneFunc = func(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) *mongo.SingleResult {
		return mongo.NewSingleResultFromDocument(expected, nil, nil)
	}

	item, err := s.repository.GetByID(s.ctx, id)
	s.NoError(err)
	s.Equal(expected, item)

	// not found
	s.itemsCol.findOneFunc = func(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) *mongo.SingleResult {
		return mongo.NewSingleResultFromDocument(bson.D{}, mongo.ErrNoDocuments, nil)
	}
	_, err = s.repository.GetByID(s.ctx, primitive.NewObjectID().Hex())
	s.ErrorIs(err, repository.ErrNotFound)

	// invalid ID
	_, err = s.repository.GetByID(s.ctx, "bad-id")
	s.Error(err)
}

// ---- Update ----
func (s *MongoSuite) TestUpdate() {
	id := primitive.NewObjectID().Hex()
	item := repository.Item{ID: id, Name: "Upd", Active: false}

	s.itemsCol.updateOneFunc = func(ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
		return &mongo.UpdateResult{MatchedCount: 1, ModifiedCount: 1}, nil
	}
	updated, err := s.repository.Update(s.ctx, item)
	s.NoError(err)
	s.Equal(item, updated)

	// not found scenario
	s.itemsCol.updateOneFunc = func(ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
		return &mongo.UpdateResult{MatchedCount: 0, ModifiedCount: 0}, nil
	}
	_, err = s.repository.Update(s.ctx, item)
	s.ErrorIs(err, repository.ErrNotFound)

	// invalid ID
	badItem := repository.Item{ID: "bad-id", Name: "x", Active: true}
	_, err = s.repository.Update(s.ctx, badItem)
	s.Error(err)
}

// ---- Delete ----
func (s *MongoSuite) TestDelete() {
	id := primitive.NewObjectID().Hex()

	s.itemsCol.deleteOneFunc = func(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
		return &mongo.DeleteResult{DeletedCount: 1}, nil
	}
	err := s.repository.Delete(s.ctx, id)
	s.NoError(err)

	// not found
	s.itemsCol.deleteOneFunc = func(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
		return &mongo.DeleteResult{DeletedCount: 0}, nil
	}
	err = s.repository.Delete(s.ctx, primitive.NewObjectID().Hex())
	s.ErrorIs(err, repository.ErrNotFound)

	// invalid ID
	err = s.repository.Delete(s.ctx, "bad-id")
	s.Error(err)
}

// ---- List ----
func (s *MongoSuite) TestList() {
	items := []repository.Item{
		{ID: primitive.NewObjectID().Hex(), Name: "A", Active: true},
		{ID: primitive.NewObjectID().Hex(), Name: "B", Active: false},
	}

	cursor := &mockCursor{}
	cursor.allFunc = func(ctx context.Context, results interface{}) error {
		dest := results.(*[]repository.Item)
		*dest = items
		return nil
	}

	s.itemsCol.findFunc = func(ctx context.Context, filter interface{}, opts ...*options.FindOptions) (dbmongo.MongoCursorOperations, error) {
		return cursor, nil
	}

	got, err := s.repository.List(s.ctx)
	s.NoError(err)
	s.Equal(items, got)
}

// ---- Additional Error Path Tests ----

func (s *MongoSuite) TestCreate_ErrorInsert() {
	expectErr := mongo.CommandError{Message: "insert failed"}
	s.itemsCol.insertOneFunc = func(ctx context.Context, doc interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
		return nil, expectErr
	}
	_, err := s.repository.Create(s.ctx, repository.Item{ID: primitive.NewObjectID().Hex(), Name: "X", Active: true})
	s.Error(err)
	s.Contains(err.Error(), "insert failed")
}

func (s *MongoSuite) TestGetByID_DatabaseError() {
	dbErr := mongo.CommandError{Message: "db down"}
	s.itemsCol.findOneFunc = func(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) *mongo.SingleResult {
		return mongo.NewSingleResultFromDocument(bson.D{}, dbErr, nil)
	}
	_, err := s.repository.GetByID(s.ctx, primitive.NewObjectID().Hex())
	s.Error(err)
	s.Contains(err.Error(), "db down")
}

func (s *MongoSuite) TestUpdate_DBError() {
	id := primitive.NewObjectID().Hex()
	dbErr := mongo.CommandError{Message: "update failed"}
	s.itemsCol.updateOneFunc = func(ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
		return nil, dbErr
	}
	_, err := s.repository.Update(s.ctx, repository.Item{ID: id, Name: "Y", Active: true})
	s.Error(err)
	s.Contains(err.Error(), "update failed")
}

func (s *MongoSuite) TestDelete_DBError() {
	id := primitive.NewObjectID().Hex()
	dbErr := mongo.CommandError{Message: "delete failed"}
	s.itemsCol.deleteOneFunc = func(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
		return nil, dbErr
	}
	err := s.repository.Delete(s.ctx, id)
	s.Error(err)
	s.Contains(err.Error(), "delete failed")
}

func (s *MongoSuite) TestList_FindError() {
	dbErr := mongo.CommandError{Message: "find failed"}
	s.itemsCol.findFunc = func(ctx context.Context, filter interface{}, opts ...*options.FindOptions) (dbmongo.MongoCursorOperations, error) {
		return nil, dbErr
	}
	_, err := s.repository.List(s.ctx)
	s.Error(err)
	s.Contains(err.Error(), "find failed")
}

func (s *MongoSuite) TestList_DecodeError() {
	cursor := &mockCursor{}
	cursor.allFunc = func(ctx context.Context, results interface{}) error {
		return mongo.CommandError{Message: "decode error"}
	}
	s.itemsCol.findFunc = func(ctx context.Context, filter interface{}, opts ...*options.FindOptions) (dbmongo.MongoCursorOperations, error) {
		return cursor, nil
	}
	_, err := s.repository.List(s.ctx)
	s.Error(err)
	s.Contains(err.Error(), "decode error")
}

// ---- CreateItemWithUser (error path) ----

func (s *MongoSuite) TestCreateItemWithUser_StartSessionError() {
	startErr := mongo.CommandError{Message: "cannot start session"}
	s.client.mcOps = &mockMongoClientOps{err: startErr}

	_, _, err := s.repository.(*mongorepo.MongoDBItemRepository).CreateItemWithUser(s.ctx, repository.Item{Name: "X"}, repository.User{})
	s.Error(err)
	s.Contains(err.Error(), "cannot start session")
}

// -------------------- run suite --------------------

func TestMongoDBSuite(t *testing.T) {
	suite.Run(t, new(MongoSuite))
}
