package github

import (
	"context"
	"time"

	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"

	"github.com/kdimtriCP/gh-inspector/internal/cache"
)

type Client struct {
	graphqlClient *githubv4.Client
	cache         cache.Cache
	cacheTTL      time.Duration
}

func NewClient(token string) *Client {
	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	httpClient := oauth2.NewClient(context.Background(), src)

	return &Client{
		graphqlClient: githubv4.NewClient(httpClient),
		cacheTTL:      1 * time.Hour,
	}
}

func (c *Client) SetCache(cache cache.Cache) {
	c.cache = cache
}

func (c *Client) SetCacheTTL(ttl time.Duration) {
	c.cacheTTL = ttl
}
