package formatter

import (
	"fmt"
	"time"

	"github.com/yourname/gh-inspector/internal/github"
)

const (
	FormatTable       = "table"
	FormatJSON        = "json"
	FormatJSONCompact = "json-compact"
	FormatCSV         = "csv"
)

type Record struct {
	Repository  string  `json:"repository"`
	Score       float64 `json:"score"`
	Stars       int     `json:"stars"`
	Forks       int     `json:"forks"`
	OpenIssues  int     `json:"open_issues"`
	OpenPRs     int     `json:"open_prs"`
	LastCommit  string  `json:"last_commit"`
	Language    string  `json:"language"`
	CICD        string  `json:"ci_cd"`
	License     string  `json:"license"`
	Description string  `json:"description"`
	Archived    string  `json:"archived"`
}

func MetricsToRecord(m *github.RepoMetrics) *Record {
	lastCommit := "N/A"
	if !m.LastCommitDate.IsZero() {
		daysAgo := int(time.Since(m.LastCommitDate).Hours() / 24)
		lastCommit = fmt.Sprintf("%d days ago", daysAgo)
	}

	cicd := "No"
	if m.HasCICD {
		cicd = "Yes"
	}

	license := "No"
	if m.HasLicense {
		license = "Yes"
	}

	lang := m.PrimaryLanguage
	if lang == "" {
		lang = "N/A"
	}

	archived := "No"
	if m.IsArchived {
		archived = "Yes"
	}

	return &Record{
		Repository:  fmt.Sprintf("%s/%s", m.Owner, m.Name),
		Score:       m.Score,
		Stars:       m.Stars,
		Forks:       m.Forks,
		OpenIssues:  m.OpenIssues,
		OpenPRs:     m.OpenPRs,
		LastCommit:  lastCommit,
		Language:    lang,
		CICD:        cicd,
		License:     license,
		Description: m.Description,
		Archived:    archived,
	}
}
