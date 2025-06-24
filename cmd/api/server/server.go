package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/lucaspereirasilva0/list-manager-api/cmd/api/handlers"
	"go.uber.org/zap"
)

// Server encapsulates the HTTP server configuration
type Server struct {
	handler handlers.ItemHandler
	logger  *zap.Logger
	server  *http.Server
}

// NewServer creates a new server instance
func NewServer(handler handlers.ItemHandler, logger *zap.Logger, port int) *Server {
	return &Server{
		handler: handler,
		logger:  logger,
		server: &http.Server{
			Addr:         fmt.Sprintf(":%d", port),
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
		},
	}
}

// setupRoutes configures the server routes
func (s *Server) setupRoutes() {
	router := mux.NewRouter()

	// Routes for item operations
	router.Handle("/item", handlers.ErrorHandlingMiddleware(s.handler.CreateItem)).Methods("POST")
	router.Handle("/item", handlers.ErrorHandlingMiddleware(s.handler.GetItem)).Methods("GET")
	router.Handle("/item", handlers.ErrorHandlingMiddleware(s.handler.UpdateItem)).Methods("PUT")
	router.Handle("/item", handlers.ErrorHandlingMiddleware(s.handler.DeleteItem)).Methods("DELETE")
	router.Handle("/items", handlers.ErrorHandlingMiddleware(s.handler.ListItems)).Methods("GET")

	// Route for application version (for PWA auto-update)
	router.HandleFunc("/_app/version.json", handlers.GetVersion).Methods("GET")

	// Middleware for logging
	loggingMiddleware := handlers.LoggingMiddleware(s.logger)

	// Middleware for CORS (allow all origins for now; adjust as needed)
	corsMiddleware := handlers.CORSMiddleware([]string{"*"})

	// Apply middlewares: CORS first, then logging, then router
	s.server.Handler = corsMiddleware(loggingMiddleware(router))
}

// Start initializes the HTTP server
func (s *Server) Start() error {
	s.setupRoutes()

	// Channel for termination signals
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	// Start server in a goroutine
	go func() {

		s.logger.Info("starting server", zap.String("addr", s.server.Addr))
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.logger.Fatal("server error", zap.Error(err))
		}
	}()

	// Wait for termination signal
	<-stop

	s.logger.Info("shutting down server")

	// Create context with timeout for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Attempt graceful shutdown
	if err := s.server.Shutdown(ctx); err != nil {
		s.logger.Error("server shutdown error", zap.Error(err))
		return err
	}

	s.logger.Info("server stopped gracefully")
	return nil
}
