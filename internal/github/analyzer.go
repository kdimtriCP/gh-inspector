package github

import (
	"context"
	"fmt"
	"time"

	"github.com/kdimtriCP/gh-inspector/internal/cache"
	"github.com/kdimtriCP/gh-inspector/internal/metrics"
	"github.com/kdimtriCP/gh-inspector/internal/scoring"
)

type RepoAnalyzer struct {
	client *Client
	scorer *scoring.Scorer
}

func NewRepoAnalyzer(token string, scoringConfig *scoring.Config) *RepoAnalyzer {
	return &RepoAnalyzer{
		client: NewClient(token),
		scorer: scoring.NewScorer(scoringConfig),
	}
}

func (ra *RepoAnalyzer) SetCache(c cache.Cache) {
	ra.client.SetCache(c)
}

func (ra *RepoAnalyzer) SetCacheTTL(ttl time.Duration) {
	ra.client.SetCacheTTL(ttl)
}

func (ra *RepoAnalyzer) Analyze(ctx context.Context, url string) (*metrics.Repository, error) {
	repo, err := ra.client.CollectBasicMetrics(ctx, url)
	if err != nil {
		return nil, fmt.Errorf("failed to collect metrics for %s: %w", url, err)
	}

	repo.Score = ra.scorer.Score(repo)

	return repo, nil
}
