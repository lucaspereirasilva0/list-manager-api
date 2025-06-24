package mongodb

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// --- Wrappers para adaptar tipos concretos do driver às nossas interfaces ---

type mongoCollectionWrapper struct {
	collection *mongo.Collection
}

func (mcw *mongoCollectionWrapper) InsertOne(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
	return mcw.collection.InsertOne(ctx, document, opts...)
}

func (mcw *mongoCollectionWrapper) FindOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) *mongo.SingleResult {
	return mcw.collection.FindOne(ctx, filter, opts...)
}

func (mcw *mongoCollectionWrapper) UpdateOne(ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	return mcw.collection.UpdateOne(ctx, filter, update, opts...)
}

func (mcw *mongoCollectionWrapper) DeleteOne(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	return mcw.collection.DeleteOne(ctx, filter, opts...)
}

func (mcw *mongoCollectionWrapper) Find(ctx context.Context, filter interface{}, opts ...*options.FindOptions) (MongoCursorOperations, error) {
	cursor, err := mcw.collection.Find(ctx, filter, opts...)
	if err != nil {
		return nil, err
	}
	return &mongoCursorWrapper{cursor: cursor}, nil
}

func (mcw *mongoCollectionWrapper) DeleteMany(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	return mcw.collection.DeleteMany(ctx, filter, opts...)
}

type mongoCursorWrapper struct {
	cursor *mongo.Cursor
}

func (mcw *mongoCursorWrapper) All(ctx context.Context, results interface{}) error {
	return mcw.cursor.All(ctx, results)
}

func (mcw *mongoCursorWrapper) Close(ctx context.Context) error {
	return mcw.cursor.Close(ctx)
}

type mongoDatabaseWrapper struct {
	database *mongo.Database
}

func (mdw *mongoDatabaseWrapper) Collection(name string, opts ...*options.CollectionOptions) MongoCollectionOperations {
	return &mongoCollectionWrapper{collection: mdw.database.Collection(name, opts...)}
}

type mongoClientWrapper struct {
	client *mongo.Client
}

func (mcw *mongoClientWrapper) StartSession(opts ...*options.SessionOptions) (mongo.Session, error) {
	return mcw.client.StartSession(opts...)
}

func (mcw *mongoClientWrapper) Database(name string, opts ...*options.DatabaseOptions) MongoDatabaseOperations {
	return &mongoDatabaseWrapper{database: mcw.client.Database(name, opts...)}
}
