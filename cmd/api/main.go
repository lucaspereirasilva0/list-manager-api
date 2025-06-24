package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/lucaspereirasilva0/list-manager-api/cmd/api/handlers"
	"github.com/lucaspereirasilva0/list-manager-api/cmd/api/server"
	dbmongo "github.com/lucaspereirasilva0/list-manager-api/internal/database/mongodb"
	repositorymongo "github.com/lucaspereirasilva0/list-manager-api/internal/repository/mongodb"
	"github.com/lucaspereirasilva0/list-manager-api/internal/service"
	"go.uber.org/zap"
)

var (
	defaultPort = 8085
	mongoURI    = "mongodb://localhost:27017"
	mongoDBName = "listmanager"
)

func main() {
	//Setup logger
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer func() {
		if err := logger.Sync(); err != nil {
			log.Printf("Failed to sync logger: %v", err)
		}
	}()

	//Get port from environment variable
	if portStr := os.Getenv("PORT"); portStr != "" {
		if p, err := strconv.Atoi(portStr); err == nil {
			defaultPort = p
		}
	}

	// Get MongoDB URI from environment variable
	if uri := os.Getenv("MONGO_URI"); uri != "" {
		mongoURI = uri
	}
	// Get MongoDB DB Name from environment variable
	if dbName := os.Getenv("MONGO_DB_NAME"); dbName != "" {
		mongoDBName = dbName
	}

	// Create context for MongoDB connection
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Create MongoDB client
	mongoClient, err := createMongoClient(ctx, mongoURI, mongoDBName, logger)
	if err != nil {
		logger.Fatal("Failed to create MongoDB client", zap.Error(err))
	}
	defer func() {
		if err := mongoClient.Disconnect(context.Background()); err != nil {
			logger.Error("Failed to disconnect MongoDB client", zap.Error(err))
		}
	}()

	//Create repository
	repository := repositorymongo.NewMongoDBItemRepository(mongoClient)

	//Create item service
	itemService := service.NewItemService(repository)
	//Create handler
	handler := handlers.NewHandler(itemService)

	//Create server
	srv := server.NewServer(handler, logger, defaultPort)
	if err := srv.Start(); err != nil {
		logger.Fatal("server error", zap.Error(err))
	}

	logger.Info("service started successfully",
		zap.Int("port", defaultPort),
	)
}

func createMongoClient(ctx context.Context, mongoURI, mongoDBName string, logger *zap.Logger) (*dbmongo.ClientWrapper, error) {
	// Create MongoDB client
	mongoClient, err := dbmongo.NewClient(ctx, mongoURI, mongoDBName)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	return mongoClient, nil
}
