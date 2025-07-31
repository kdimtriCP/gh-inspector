package github

import (
	"context"
	"time"

	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"

	"github.com/kdimtriCP/gh-inspector/internal/cache"
	"github.com/kdimtriCP/gh-inspector/internal/metrics"
)

type Client struct {
	graphqlClient   *githubv4.Client
	cache           cache.Cache
	cacheTTL        time.Duration
	metricsRecorder metrics.Recorder
}

func NewClient(token string) *Client {
	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	httpClient := oauth2.NewClient(context.Background(), src)

	return &Client{
		graphqlClient:   githubv4.NewClient(httpClient),
		cacheTTL:        1 * time.Hour,
		metricsRecorder: &metrics.NoOpRecorder{},
	}
}

func (c *Client) SetCache(cache cache.Cache) {
	c.cache = cache
}

func (c *Client) SetCacheTTL(ttl time.Duration) {
	c.cacheTTL = ttl
}

func (c *Client) SetMetricsRecorder(recorder metrics.Recorder) {
	c.metricsRecorder = recorder
}
