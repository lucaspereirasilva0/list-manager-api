package main

import (
	"log"
	"os"
	"strconv"

	"github.com/lucaspereirasilva0/list-manager-api/cmd/api/handlers"
	"github.com/lucaspereirasilva0/list-manager-api/cmd/api/server"
	"github.com/lucaspereirasilva0/list-manager-api/internal/repository/local"
	"github.com/lucaspereirasilva0/list-manager-api/internal/service"
	"go.uber.org/zap"
)

var (
	defaultPort = 8081
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

	//Create repository
	repository := local.NewLocalRepository()

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
