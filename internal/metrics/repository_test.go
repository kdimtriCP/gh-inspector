package metrics

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestRepositoryMethods(t *testing.T) {
	now := time.Now()
	repo := &Repository{
		Owner:           "test",
		Name:            "repo",
		Stars:           100,
		Forks:           10,
		OpenIssues:      5,
		OpenPRs:         2,
		LastCommitDate:  now,
		Score:           75.5,
		PrimaryLanguage: "Go",
		HasLicense:      true,
		HasCICD:         true,
		HasContributing: false,
		IsArchived:      false,
		Description:     "Test repository",
	}

	tests := []struct {
		name   string
		method string
		want   interface{}
	}{
		{name: "GetStars", method: "stars", want: 100},
		{name: "GetForks", method: "forks", want: 10},
		{name: "GetOpenIssues", method: "issues", want: 5},
		{name: "GetOpenPRs", method: "prs", want: 2},
		{name: "GetLastCommitDate", method: "commit", want: now},
		{name: "GetIsArchived", method: "archived", want: false},
		{name: "GetHasLicense", method: "license", want: true},
		{name: "GetHasCICD", method: "cicd", want: true},
		{name: "GetHasContributing", method: "contributing", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got interface{}
			switch tt.method {
			case "stars":
				got = repo.GetStars()
			case "forks":
				got = repo.GetForks()
			case "issues":
				got = repo.GetOpenIssues()
			case "prs":
				got = repo.GetOpenPRs()
			case "commit":
				got = repo.GetLastCommitDate()
			case "archived":
				got = repo.GetIsArchived()
			case "license":
				got = repo.GetHasLicense()
			case "cicd":
				got = repo.GetHasCICD()
			case "contributing":
				got = repo.GetHasContributing()
			}

			require.Equal(t, tt.want, got, "%s() should return correct value", tt.name)
		})
	}
}

func TestDaysSinceLastCommit(t *testing.T) {
	tests := []struct {
		name           string
		lastCommitDate time.Time
		wantDays       int
	}{
		{
			name:           "commit today",
			lastCommitDate: time.Now(),
			wantDays:       0,
		},
		{
			name:           "commit yesterday",
			lastCommitDate: time.Now().AddDate(0, 0, -1),
			wantDays:       1,
		},
		{
			name:           "commit last week",
			lastCommitDate: time.Now().AddDate(0, 0, -7),
			wantDays:       7,
		},
		{
			name:           "commit last month",
			lastCommitDate: time.Now().AddDate(0, -1, 0),
			wantDays:       30,
		},
		{
			name:           "zero time",
			lastCommitDate: time.Time{},
			wantDays:       -1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &Repository{LastCommitDate: tt.lastCommitDate}
			got := repo.DaysSinceLastCommit()

			if tt.lastCommitDate.IsZero() {
				require.Equal(t, -1, got, "DaysSinceLastCommit() should return -1 for zero time")
			} else {
				require.InDelta(t, tt.wantDays, got, 1, "DaysSinceLastCommit() = %v, want approximately %v", got, tt.wantDays)
			}
		})
	}
}

func TestFullName(t *testing.T) {
	tests := []struct {
		name     string
		owner    string
		repoName string
		want     string
	}{
		{
			name:     "standard repository",
			owner:    "golang",
			repoName: "go",
			want:     "golang/go",
		},
		{
			name:     "repository with hyphen",
			owner:    "kubernetes",
			repoName: "kubernetes",
			want:     "kubernetes/kubernetes",
		},
		{
			name:     "empty owner",
			owner:    "",
			repoName: "repo",
			want:     "/repo",
		},
		{
			name:     "empty name",
			owner:    "owner",
			repoName: "",
			want:     "owner/",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &Repository{
				Owner: tt.owner,
				Name:  tt.repoName,
			}
			got := repo.FullName()
			require.Equal(t, tt.want, got)
		})
	}
}
