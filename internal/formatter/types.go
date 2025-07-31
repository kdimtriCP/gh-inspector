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

// Record represents a scored repository
// @Description Repository scoring results
type Record struct {
	// Repository name in owner/name format
	Repository string `json:"repository" example:"kubernetes/kubernetes"`
	// Repository score (0-100)
	Score float64 `json:"score" example:"95.5"`
	// Number of stars
	Stars int `json:"stars" example:"108000"`
	// Number of forks
	Forks int `json:"forks" example:"39000"`
	// Number of watchers
	Watchers int `json:"watchers" example:"3500"`
	// Number of open issues
	OpenIssues int `json:"open_issues" example:"1500"`
	// Number of open pull requests
	OpenPRs int `json:"open_prs" example:"300"`
	// Last commit relative time
	LastCommit string `json:"last_commit" example:"1 days ago"`
	// Number of releases
	Releases int `json:"releases" example:"350"`
	// Last release relative time
	LastRelease string `json:"last_release" example:"7 days ago"`
	// Primary programming language
	Language string `json:"language" example:"Go"`
	// CI/CD presence
	CICD string `json:"ci_cd" example:"Yes" enums:"Yes,No"`
	// License presence
	License string `json:"license" example:"Yes" enums:"Yes,No"`
	// Contributing guide presence
	Contributing string `json:"contributing" example:"Yes" enums:"Yes,No"`
	// README presence
	Readme string `json:"readme" example:"Yes" enums:"Yes,No"`
	// Code of conduct presence
	CodeOfConduct string `json:"code_of_conduct" example:"Yes" enums:"Yes,No"`
	// Security policy presence
	Security string `json:"security" example:"Yes" enums:"Yes,No"`
	// Repository description
	Description string `json:"description" example:"Production-Grade Container Scheduling and Management"`
	// Archive status
	Archived string `json:"archived" example:"No" enums:"Yes,No"`
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
