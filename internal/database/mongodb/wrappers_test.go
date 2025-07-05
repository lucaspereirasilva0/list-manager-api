package mongodb_test

import (
	"context"
	"testing"

	"github.com/lucaspereirasilva0/list-manager-api/internal/database/mongodb"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

func TestCreateWrapper(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("test", func(mt *mtest.T) {
		mt.AddMockResponses(mtest.CreateSuccessResponse())

		collection := mongodb.NewMockCollectionWrapper(mt)
		result, err := collection.InsertOne(context.Background(), bson.D{{Key: "name", Value: "test"}})
		require.NoError(t, err)
		require.NotNil(t, result)
	})
}

func TestFindOneWrapper(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("test", func(mt *mtest.T) {
		mt.AddMockResponses(mtest.CreateCursorResponse(1, "foo.bar", mtest.FirstBatch, bson.D{{Key: "name", Value: "test"}}))

		collection := mongodb.NewMockCollectionWrapper(mt)
		result := collection.FindOne(context.Background(), bson.D{{Key: "name", Value: "test"}})
		require.NotNil(t, result)

		var doc bson.M
		err := result.Decode(&doc)
		require.NoError(t, err)
		require.Equal(t, "test", doc["name"])
	})
}

func TestUpdateOneWrapper(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("test", func(mt *mtest.T) {
		mt.AddMockResponses(bson.D{{Key: "ok", Value: 1}, {Key: "n", Value: 1}, {Key: "nModified", Value: 1}})

		collection := mongodb.NewMockCollectionWrapper(mt)
		result, err := collection.UpdateOne(context.Background(), bson.D{{Key: "name", Value: "test"}}, bson.D{{Key: "$set", Value: bson.D{{Key: "name", Value: "new_test"}}}})
		require.NoError(t, err)
		require.NotNil(t, result)
		require.Equal(t, int64(1), result.ModifiedCount)
	})
}

func TestDeleteOneWrapper(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("test", func(mt *mtest.T) {
		mt.AddMockResponses(bson.D{{Key: "ok", Value: 1}, {Key: "n", Value: 1}})

		collection := mongodb.NewMockCollectionWrapper(mt)
		result, err := collection.DeleteOne(context.Background(), bson.D{{Key: "name", Value: "test"}})
		require.NoError(t, err)
		require.NotNil(t, result)
		require.Equal(t, int64(1), result.DeletedCount)
	})
}

func TestFindWrapper(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("test", func(mt *mtest.T) {
		mt.AddMockResponses(
			mtest.CreateCursorResponse(1, "foo.bar", mtest.FirstBatch, bson.D{{Key: "name", Value: "test1"}}, bson.D{{Key: "name", Value: "test2"}}),
			mtest.CreateCursorResponse(0, "foo.bar", mtest.NextBatch),
		)

		collection := mongodb.NewMockCollectionWrapper(mt)
		cursor, err := collection.Find(context.Background(), bson.D{})
		require.NoError(t, err)
		require.NotNil(t, cursor)

		var results []bson.M
		err = cursor.All(context.Background(), &results)
		require.NoError(t, err)
		require.Len(t, results, 2)
	})
}

func TestDeleteManyWrapper(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("test", func(mt *mtest.T) {
		mt.AddMockResponses(bson.D{{Key: "ok", Value: 1}, {Key: "n", Value: 2}})

		collection := mongodb.NewMockCollectionWrapper(mt)
		result, err := collection.DeleteMany(context.Background(), bson.D{{Key: "name", Value: "test"}})
		require.NoError(t, err)
		require.NotNil(t, result)
		require.Equal(t, int64(2), result.DeletedCount)
	})
}

func TestDatabaseCollectionWrapper(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("test", func(mt *mtest.T) {
		mt.AddMockResponses(mtest.CreateSuccessResponse())
		dbMock := mongodb.NewMockMongoClientWrappers(mt)
		collection := dbMock.Database("test_db").Collection("test_collection")
		require.NotNil(t, collection)
	})
}

func TestClientStartSessionWrapper(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("test", func(mt *mtest.T) {
		mt.AddMockResponses(mtest.CreateSuccessResponse())
		dbMock := mongodb.NewMockMongoClientWrappers(mt)
		session, err := dbMock.StartSession(nil)
		require.NoError(t, err)
		require.NotNil(t, session)
		defer session.EndSession(context.Background())
	})
}
