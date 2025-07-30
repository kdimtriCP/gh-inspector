package github

import (
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	
	"github.com/kdimtriCP/gh-inspector/internal/mock/mock_cache"
)

func TestNewClient(t *testing.T) {
	token := "test-token"
	client := NewClient(token)

	require.NotNil(t, client)
	require.NotNil(t, client.graphqlClient)
	require.Equal(t, 1*time.Hour, client.cacheTTL)
	require.Nil(t, client.cache, "cache should be nil by default")
}

func TestClientSetCache(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	client := NewClient("test-token")
	mockCache := mock_cache.NewMockCache(ctrl)

	client.SetCache(mockCache)
	
	require.Equal(t, mockCache, client.cache)
}

func TestClientSetCacheTTL(t *testing.T) {
	client := NewClient("test-token")
	
	newTTL := 30 * time.Minute
	client.SetCacheTTL(newTTL)
	
	require.Equal(t, newTTL, client.cacheTTL)
}