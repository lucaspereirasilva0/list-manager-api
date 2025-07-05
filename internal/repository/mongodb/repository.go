package mongodb

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	dbmongo "github.com/lucaspereirasilva0/list-manager-api/internal/database/mongodb"
	"github.com/lucaspereirasilva0/list-manager-api/internal/repository"
)

const (
	CollectionItems = "items"
	CollectionUsers = "users"
)

// MongoDBItemRepository implements repository.ItemRepository for MongoDB
type MongoDBItemRepository struct {
	client dbmongo.ClientOperations
}

// NewMongoDBItemRepository creates a new instance of MongoDBItemRepository
func NewMongoDBItemRepository(client dbmongo.ClientOperations) repository.ItemRepository {
	return &MongoDBItemRepository{
		client: client,
	}
}

// Create inserts a new item in the MongoDB repository
func (r *MongoDBItemRepository) Create(ctx context.Context, item repository.Item) (repository.Item, error) {
	collection := r.client.GetCollection(CollectionItems)

	// Converte o ID fornecido (agora um hexadecimal de 24 caracteres) para um ObjectID
	objectID, err := primitive.ObjectIDFromHex(item.ID)
	if err != nil {
		return repository.Item{}, repository.ErrInvalidHexID
	}

	// Use o ObjectID para a inserção no MongoDB
	_, err = collection.InsertOne(ctx, bson.M{
		"_id":    objectID,
		"name":   item.Name,
		"active": item.Active,
	})
	if err != nil {
		return repository.Item{}, fmt.Errorf("failed to create item: %w", err)
	}

	return item, nil
}

// Update modifies an existing item in the MongoDB repository
func (r *MongoDBItemRepository) Update(ctx context.Context, item repository.Item) (repository.Item, error) {
	collection := r.client.GetCollection(CollectionItems)

	id, err := primitive.ObjectIDFromHex(item.ID)
	if err != nil {
		return repository.Item{}, repository.ErrInvalidHexID
	}

	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{
		"name":   item.Name,
		"active": item.Active,
	}}

	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return repository.Item{}, fmt.Errorf("failed to update item: %w", err)
	}

	if result.MatchedCount == 0 {
		return repository.Item{}, repository.ErrNotFound
	}

	return item, nil
}

// Delete removes an item from the MongoDB repository
func (r *MongoDBItemRepository) Delete(ctx context.Context, id string) error {
	collection := r.client.GetCollection(CollectionItems)

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return repository.ErrInvalidHexID
	}

	filter := bson.M{"_id": objID}
	result, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("failed to delete item: %w", err)
	}

	if result.DeletedCount == 0 {
		return repository.ErrNotFound
	}

	return nil
}

// GetByID retrieves an item by its ID from the MongoDB repository
func (r *MongoDBItemRepository) GetByID(ctx context.Context, id string) (repository.Item, error) {
	collection := r.client.GetCollection(CollectionItems)

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return repository.Item{}, repository.ErrInvalidHexID
	}

	filter := bson.M{"_id": objID}
	var item repository.Item

	err = collection.FindOne(ctx, filter).Decode(&item)
	if err == mongo.ErrNoDocuments {
		return repository.Item{}, repository.ErrNotFound
	} else if err != nil {
		return repository.Item{}, fmt.Errorf("failed to get item by ID: %w", err)
	}

	return item, nil
}

// List retrieves all items from the MongoDB repository
func (r *MongoDBItemRepository) List(ctx context.Context) ([]repository.Item, error) {
	collection := r.client.GetCollection(CollectionItems)

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("failed to list items: %w", err)
	}
	defer func() {
		if cursor != nil {
			if err := cursor.Close(ctx); err != nil {
				log.Printf("Error closing MongoDB cursor: %v", err)
			}
		}
	}()

	var items []repository.Item
	if err = cursor.All(ctx, &items); err != nil {
		return nil, fmt.Errorf("failed to decode items: %w", err)
	}

	return items, nil
}

//TODO adicionar quando implementar autenticacao de usuario
// // CreateItemWithUser inserts a new item and an associated user in a single transaction
// func (r *MongoDBItemRepository) CreateItemWithUser(ctx context.Context, item repository.Item, user repository.User) (repository.Item, repository.User, error) {
// 	session, err := r.client.Client().StartSession()
// 	if err != nil {
// 		return repository.Item{}, repository.User{}, fmt.Errorf("failed to start session: %w", err)
// 	}
// 	defer session.EndSession(ctx)

// 	_, err = session.WithTransaction(ctx, func(sessCtx mongo.SessionContext) (interface{}, error) {
// 		itemsCollection := r.client.GetCollection(CollectionItems)
// 		usersCollection := r.client.GetCollection(CollectionUsers)

// 		// Use the ID provided from the domain layer (already in hex format)
// 		objectID, err := primitive.ObjectIDFromHex(item.ID)
// 		if err != nil {
// 			return nil, fmt.Errorf("provided item ID is not a valid ObjectID: %w", err)
// 		}

// 		// Insert item
// 		_, err = itemsCollection.InsertOne(sessCtx, bson.M{
// 			"_id":    objectID,
// 			"name":   item.Name,
// 			"active": item.Active,
// 		})
// 		if err != nil {
// 			return nil, fmt.Errorf("failed to insert item: %w", err)
// 		}

// 		// Generate a new ObjectID for the user if it's not set
// 		if user.ID == "" {
// 			user.ID = primitive.NewObjectID().Hex()
// 		}

// 		// Set CreatedBy to the new item's ID for association
// 		user.CreatedBy = item.ID

// 		// Insert user
// 		_, err = usersCollection.InsertOne(sessCtx, user)
// 		if err != nil {
// 			return nil, fmt.Errorf("failed to insert user: %w", err)
// 		}

// 		return nil, nil
// 	})

// 	if err != nil {
// 		return repository.Item{}, repository.User{}, fmt.Errorf("transaction failed: %w", err)
// 	}

// 	return item, user, nil
// }
