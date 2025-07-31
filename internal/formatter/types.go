package formatter

import (
	"fmt"
	"time"

	"github.com/kdimtriCP/gh-inspector/internal/metrics"
)

const (
	FormatTable       = "table"
	FormatJSON        = "json"
	FormatJSONCompact = "json-compact"
	FormatCSV         = "csv"
)

type Record struct {
	Repository    string  `json:"repository"`
	Score         float64 `json:"score"`
	Stars         int     `json:"stars"`
	Forks         int     `json:"forks"`
	Watchers      int     `json:"watchers"`
	OpenIssues    int     `json:"open_issues"`
	OpenPRs       int     `json:"open_prs"`
	LastCommit    string  `json:"last_commit"`
	Releases      int     `json:"releases"`
	LastRelease   string  `json:"last_release"`
	Language      string  `json:"language"`
	CICD          string  `json:"ci_cd"`
	License       string  `json:"license"`
	Contributing  string  `json:"contributing"`
	Readme        string  `json:"readme"`
	CodeOfConduct string  `json:"code_of_conduct"`
	Security      string  `json:"security"`
	Description   string  `json:"description"`
	Archived      string  `json:"archived"`
}

func MetricsToRecord(m *metrics.Repository) *Record {
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

	lastRelease := "Never"
	if !m.LastReleaseDate.IsZero() {
		daysAgo := int(time.Since(m.LastReleaseDate).Hours() / 24)
		lastRelease = fmt.Sprintf("%d days ago", daysAgo)
	}

	contributing := "No"
	if m.HasContributing {
		contributing = "Yes"
	}

	readme := "No"
	if m.HasReadme {
		readme = "Yes"
	}

	codeOfConduct := "No"
	if m.HasCodeOfConduct {
		codeOfConduct = "Yes"
	}

	security := "No"
	if m.HasSecurity {
		security = "Yes"
	}

	return &Record{
		Repository:    fmt.Sprintf("%s/%s", m.Owner, m.Name),
		Score:         m.Score,
		Stars:         m.Stars,
		Forks:         m.Forks,
		Watchers:      m.Watchers,
		OpenIssues:    m.OpenIssues,
		OpenPRs:       m.OpenPRs,
		LastCommit:    lastCommit,
		Releases:      m.ReleaseCount,
		LastRelease:   lastRelease,
		Language:      lang,
		CICD:          cicd,
		License:       license,
		Contributing:  contributing,
		Readme:        readme,
		CodeOfConduct: codeOfConduct,
		Security:      security,
		Description:   m.Description,
		Archived:      archived,
	}
}

func (r *Record) String() string {
	return fmt.Sprintf(
		"Repository: %s, Score: %.1f, Stars: %d, Forks: %d, Open Issues: %d, Open PRs: %d, Last Commit: %s, Releases: %d, Last Release: %s, Language: %s, CI/CD: %s, License: %s, Description: %s, Archived: %s",
		r.Repository,
		r.Score,
		r.Stars,
		r.Forks,
		r.OpenIssues,
		r.OpenPRs,
		r.LastCommit,
		r.Releases,
		r.LastRelease,
		r.Language,
		r.CICD,
		r.License,
		r.Description,
		r.Archived,
	)
}

func (r *Record) Strings() []string {
	return []string{
		r.Repository,
		fmt.Sprintf("%.1f", r.Score),
		fmt.Sprintf("%d", r.Stars),
		fmt.Sprintf("%d", r.Forks),
		fmt.Sprintf("%d", r.Watchers),
		fmt.Sprintf("%d", r.OpenIssues),
		fmt.Sprintf("%d", r.OpenPRs),
		r.LastCommit,
		fmt.Sprintf("%d", r.Releases),
		r.LastRelease,
		r.Language,
		r.CICD,
		r.License,
		r.Contributing,
		r.Description,
		r.Archived,
	}
}

func GetRecordHeaders() []string {
	return []string{
		"Repository",
		"Score",
		"Stars",
		"Forks",
		"Watchers",
		"Open Issues",
		"Open PRs",
		"Last Commit",
		"Releases",
		"Last Release",
		"Language",
		"CI/CD",
		"License",
		"Contributing",
		"Description",
		"Archived",
	}
}
