package github

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type RepoMetrics struct {
	Owner           string
	Name            string
	Stars           int
	Forks           int
	OpenIssues      int
	OpenPRs         int
	LastCommitDate  time.Time
	Description     string
	PrimaryLanguage string
	IsArchived      bool
	HasCICD         bool
	HasLicense      bool
	HasContributing bool
	Score           float64
}

func (m *RepoMetrics) GetStars() int                { return m.Stars }
func (m *RepoMetrics) GetForks() int                { return m.Forks }
func (m *RepoMetrics) GetOpenIssues() int           { return m.OpenIssues }
func (m *RepoMetrics) GetOpenPRs() int              { return m.OpenPRs }
func (m *RepoMetrics) GetLastCommitDate() time.Time { return m.LastCommitDate }
func (m *RepoMetrics) GetIsArchived() bool          { return m.IsArchived }
func (m *RepoMetrics) GetHasLicense() bool          { return m.HasLicense }
func (m *RepoMetrics) GetHasCICD() bool             { return m.HasCICD }
func (m *RepoMetrics) GetHasContributing() bool     { return m.HasContributing }

type repositoryData struct {
	Repository struct {
		Owner struct {
			Login string `json:"login"`
		} `json:"owner"`
		Name            string `json:"name"`
		Description     string `json:"description"`
		StargazerCount  int    `json:"stargazerCount"`
		ForkCount       int    `json:"forkCount"`
		IsArchived      bool   `json:"isArchived"`
		PrimaryLanguage *struct {
			Name string `json:"name"`
		} `json:"primaryLanguage"`
		Issues struct {
			TotalCount int `json:"totalCount"`
		} `json:"issues"`
		PullRequests struct {
			TotalCount int `json:"totalCount"`
		} `json:"pullRequests"`
		DefaultBranchRef struct {
			Target struct {
				History struct {
					Edges []struct {
						Node struct {
							CommittedDate time.Time `json:"committedDate"`
						} `json:"node"`
					} `json:"edges"`
				} `json:"history"`
			} `json:"target"`
		} `json:"defaultBranchRef"`
		LicenseInfo *struct {
			Key string `json:"key"`
		} `json:"licenseInfo"`
		Object struct {
			Entries []struct {
				Name string `json:"name"`
			} `json:"entries"`
		} `json:"object"`
	} `json:"repository"`
}

func (c *Client) CollectBasicMetrics(ctx context.Context, repoFullName string) (*RepoMetrics, error) {
	parts := strings.Split(repoFullName, "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid repository format, expected owner/name")
	}
	owner, name := parts[0], parts[1]

	query := `
		query($owner: String!, $name: String!) {
			repository(owner: $owner, name: $name) {
				owner {
					login
				}
				name
				description
				stargazerCount
				forkCount
				isArchived
				primaryLanguage {
					name
				}
				issues(states: OPEN) {
					totalCount
				}
				pullRequests(states: OPEN) {
					totalCount
				}
				defaultBranchRef {
					target {
						... on Commit {
							history(first: 1) {
								edges {
									node {
										committedDate
									}
								}
							}
						}
					}
				}
				licenseInfo {
					key
				}
				object(expression: "HEAD:") {
					... on Tree {
						entries {
							name
						}
					}
				}
			}
		}`

	variables := map[string]interface{}{
		"owner": owner,
		"name":  name,
	}

	data, err := c.executeGraphQL(ctx, query, variables)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch repository data: %w", err)
	}

	var result repositoryData
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to parse repository data: %w", err)
	}

	metrics := &RepoMetrics{
		Owner:       result.Repository.Owner.Login,
		Name:        result.Repository.Name,
		Stars:       result.Repository.StargazerCount,
		Forks:       result.Repository.ForkCount,
		OpenIssues:  result.Repository.Issues.TotalCount,
		OpenPRs:     result.Repository.PullRequests.TotalCount,
		Description: result.Repository.Description,
		IsArchived:  result.Repository.IsArchived,
		HasLicense:  result.Repository.LicenseInfo != nil,
	}

	if result.Repository.PrimaryLanguage != nil {
		metrics.PrimaryLanguage = result.Repository.PrimaryLanguage.Name
	}

	if len(result.Repository.DefaultBranchRef.Target.History.Edges) > 0 {
		metrics.LastCommitDate = result.Repository.DefaultBranchRef.Target.History.Edges[0].Node.CommittedDate
	}

	for _, entry := range result.Repository.Object.Entries {
		entryLower := strings.ToLower(entry.Name)
		if strings.HasPrefix(entryLower, ".github") ||
			strings.HasPrefix(entryLower, ".gitlab") ||
			strings.HasPrefix(entryLower, ".circleci") ||
			entryLower == ".travis.yml" ||
			entryLower == "jenkinsfile" {
			metrics.HasCICD = true
		}
		if entryLower == "contributing.md" || entryLower == "contributing" {
			metrics.HasContributing = true
		}
	}

	return metrics, nil
}
