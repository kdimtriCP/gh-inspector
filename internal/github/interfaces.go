package github

import (
	"context"

	"github.com/kdimtriCP/gh-inspector/internal/metrics"
)

//go:generate mockgen -source=$GOFILE -destination=../mock/mock_github/mock_$GOFILE -package=mock_github

type MetricsCollector interface {
	CollectBasicMetrics(ctx context.Context, repoFullName string) (*metrics.Repository, error)
}

type Analyzer interface {
	Analyze(ctx context.Context, repo string) (*metrics.Repository, error)
}
