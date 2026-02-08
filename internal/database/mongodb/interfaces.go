package mongodb

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoCursorOperations define as operações aplicáveis a um cursor do MongoDB.
type MongoCursorOperations interface {
	All(ctx context.Context, results interface{}) error
	Close(ctx context.Context) error
}

// MongoCollectionOperations define as operações aplicáveis a uma coleção do MongoDB.
type MongoCollectionOperations interface {
	InsertOne(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error)
	FindOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) *mongo.SingleResult
	UpdateOne(ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error)
	DeleteOne(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error)
	Find(ctx context.Context, filter interface{}, opts ...*options.FindOptions) (MongoCursorOperations, error)
	DeleteMany(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error)
}

// MongoDatabaseOperations define as operações aplicáveis a um database MongoDB.
type MongoDatabaseOperations interface {
	Collection(name string, opts ...*options.CollectionOptions) MongoCollectionOperations
}

// MongoClientOperations define as operações aplicáveis a um cliente MongoDB.
type MongoClientOperations interface {
	StartSession(opts ...*options.SessionOptions) (mongo.Session, error)
	Database(name string, opts ...*options.DatabaseOptions) MongoDatabaseOperations
}

// ClientOperations define as operações expostas pelo client wrapper usado nos repositórios.
type ClientOperations interface {
	GetCollection(collectionName string) MongoCollectionOperations
	Disconnect(ctx context.Context) error
	Client() MongoClientOperations // expõe o client subjacente para controle de sessão/transaction
}

// PingClientOperations define as operações de ping para health check.
type PingClientOperations interface {
	Ping(ctx context.Context) error
}
