package github

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/shurcooL/githubv4"

	"github.com/kdimtriCP/gh-inspector/internal/cache"
	"github.com/kdimtriCP/gh-inspector/internal/metrics"
)

func (c *Client) CollectBasicMetrics(ctx context.Context, repoFullName string) (*metrics.Repository, error) {
	if c.cache != nil {
		cacheKey := cache.GenerateKey("repo", repoFullName)
		if data, found, err := c.cache.Get(cacheKey); err == nil && found {
			var result metrics.Repository
			if err := json.Unmarshal(data, &result); err == nil {
				return &result, nil
			}
		}
	}

	parts := strings.Split(repoFullName, "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid repository format, expected owner/name")
	}
	owner, name := parts[0], parts[1]

	var query metrics.RepositoryQuery
	variables := map[string]interface{}{
		metrics.VarOwner: githubv4.String(owner),
		metrics.VarName:  githubv4.String(name),
	}

	err := c.graphqlClient.Query(ctx, &query, variables)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch repository data: %w", err)
	}

	repo := query.Repository
	result := &metrics.Repository{
		Owner:       string(repo.Owner.Login),
		Name:        string(repo.Name),
		Stars:       int(repo.StargazerCount),
		Forks:       int(repo.ForkCount),
		OpenIssues:  int(repo.Issues.TotalCount),
		OpenPRs:     int(repo.PullRequests.TotalCount),
		Description: string(repo.Description),
		IsArchived:  bool(repo.IsArchived),
		HasLicense:  repo.LicenseInfo != nil,
	}

	if repo.PrimaryLanguage != nil {
		result.PrimaryLanguage = string(repo.PrimaryLanguage.Name)
	}

	if repo.DefaultBranchRef != nil && len(repo.DefaultBranchRef.Target.Commit.History.Edges) > 0 {
		result.LastCommitDate = repo.DefaultBranchRef.Target.Commit.History.Edges[0].Node.CommittedDate.Time
	}

	for _, entry := range repo.Object.Tree.Entries {
		entryLower := strings.ToLower(string(entry.Name))
		if strings.HasPrefix(entryLower, metrics.CIGitHub) ||
			strings.HasPrefix(entryLower, metrics.CIGitLab) ||
			strings.HasPrefix(entryLower, metrics.CICircleCI) ||
			entryLower == metrics.CITravis ||
			entryLower == metrics.CIJenkins {
			result.HasCICD = true
		}
		if strings.HasPrefix(entryLower, metrics.FileContributingAlt) {
			result.HasContributing = true
		}
	}

	result.ReleaseCount = int(repo.Releases.TotalCount)
	if len(repo.Releases.Edges) > 0 {
		result.LastReleaseDate = repo.Releases.Edges[0].Node.PublishedAt.Time
	}

	if c.cache != nil {
		cacheKey := cache.GenerateKey("repo", repoFullName)
		if data, err := json.Marshal(result); err == nil {
			_ = c.cache.Set(cacheKey, data, c.cacheTTL)
		}
	}

	return result, nil
}
