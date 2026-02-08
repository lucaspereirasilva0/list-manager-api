package mongodb

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// ClientWrapper holds the MongoDB client instance and helper methods.
type ClientWrapper struct {
	mongoClient  *mongo.Client
	DatabaseName string
}

// NewClient creates a new MongoDB client, connects and pings the server.
func NewClient(ctx context.Context, uri, dbName string) (*ClientWrapper, error) {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)

	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	if err = client.Ping(ctx, readpref.Primary()); err != nil {
		if disconnectErr := client.Disconnect(ctx); disconnectErr != nil {
			log.Printf("error disconnecting MongoDB client after failed ping: %v", disconnectErr)
		}
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	fmt.Println("Successfully connected and pinged MongoDB!")

	return &ClientWrapper{
		mongoClient:  client,
		DatabaseName: dbName,
	}, nil
}

// Disconnect closes the MongoDB client connection.
func (cw *ClientWrapper) Disconnect(ctx context.Context) error {
	if cw.mongoClient == nil {
		return nil
	}

	disconnectCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := cw.mongoClient.Disconnect(disconnectCtx); err != nil {
		return fmt.Errorf("failed to disconnect from MongoDB: %w", err)
	}
	fmt.Println("Disconnected from MongoDB.")
	return nil
}

// GetCollection returns a collection wrapper for the configured database.
func (cw *ClientWrapper) GetCollection(name string) MongoCollectionOperations {
	return &mongoCollectionWrapper{collection: cw.mongoClient.Database(cw.DatabaseName).Collection(name)}
}

// Client exposes the underlying *mongo.Client for advanced use cases.
func (cw *ClientWrapper) Client() MongoClientOperations {
	return &mongoClientWrapper{client: cw.mongoClient}
}

// Ping verifies the connection to MongoDB by sending a ping command.
func (cw *ClientWrapper) Ping(ctx context.Context) error {
	if cw.mongoClient == nil {
		return fmt.Errorf("mongoClient is nil")
	}
	return cw.mongoClient.Ping(ctx, readpref.Primary())
}
