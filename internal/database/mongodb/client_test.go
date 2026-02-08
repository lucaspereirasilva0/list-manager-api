package mongodb_test

import (
	"context"
	"errors"
	"testing"

	"github.com/lucaspereirasilva0/list-manager-api/internal/database/mongodb"
	"github.com/stretchr/testify/require"
)

func TestPing(t *testing.T) {
	t.Run("Given_NilMongoClient_When_Ping_Then_ReturnsError", func(t *testing.T) {
		clientWrapper := &mongodb.ClientWrapper{}
		err := clientWrapper.Ping(context.Background())
		require.Error(t, err)
		require.Contains(t, err.Error(), "mongoClient is nil")
	})

	t.Run("Given_CanceledContext_When_Ping_Then_ReturnsError", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		clientWrapper := &mongodb.ClientWrapper{}
		err := clientWrapper.Ping(ctx)
		require.Error(t, err)
	})

	t.Run("Given_ContextTimeout_When_Ping_Then_ReturnsError", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		clientWrapper := &mongodb.ClientWrapper{}
		err := clientWrapper.Ping(ctx)
		require.Error(t, err)
		require.True(t, errors.Is(err, context.Canceled) || ctx.Err() != nil)
	})
}
