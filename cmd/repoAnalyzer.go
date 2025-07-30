package cmd

import (
	"context"
	"fmt"

	"github.com/yourname/gh-inspector/internal/github"
)

type RepoAnalyzer struct {
	client *github.Client
}

func NewRepoAnalyzer(token string) *RepoAnalyzer {
	return &RepoAnalyzer{
		client: github.NewClient(token),
	}
}

func (ra *RepoAnalyzer) Analyze(ctx context.Context, repo string) (*github.RepoMetrics, error) {
	metrics, err := ra.client.CollectBasicMetrics(ctx, repo)
	if err != nil {
		return nil, fmt.Errorf("failed to collect metrics for %s: %w", repo, err)
	}
	return metrics, nil
}
