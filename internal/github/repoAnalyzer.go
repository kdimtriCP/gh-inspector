package github

import (
	"context"
	"fmt"
)

type RepoAnalyzer struct {
	client *Client
}

func NewRepoAnalyzer(token string) *RepoAnalyzer {
	return &RepoAnalyzer{
		client: NewClient(token),
	}
}

func (ra *RepoAnalyzer) Analyze(ctx context.Context, repo string) (*RepoMetrics, error) {
	metrics, err := ra.client.CollectBasicMetrics(ctx, repo)
	if err != nil {
		return nil, fmt.Errorf("failed to collect metrics for %s: %w", repo, err)
	}
	return metrics, nil
}
