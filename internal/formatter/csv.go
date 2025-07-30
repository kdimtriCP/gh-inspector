package formatter

import (
	"encoding/csv"
	"fmt"
	"io"

	"github.com/yourname/gh-inspector/internal/github"
)

type CSVFormatter struct{}

func NewCSVFormatter() *CSVFormatter {
	return &CSVFormatter{}
}

func (f *CSVFormatter) Format(writer io.Writer, metrics []*github.RepoMetrics) error {
	w := csv.NewWriter(writer)
	defer w.Flush()

	headers := []string{
		"Repository",
		"Score",
		"Stars",
		"Forks",
		"Open Issues",
		"Open PRs",
		"Last Commit",
		"Language",
		"CI/CD",
		"License",
		"Description",
		"Archived",
	}
	if err := w.Write(headers); err != nil {
		return err
	}

	for _, m := range metrics {
		record := MetricsToRecord(m)
		row := []string{
			record.Repository,
			fmt.Sprintf("%.1f", record.Score),
			fmt.Sprintf("%d", record.Stars),
			fmt.Sprintf("%d", record.Forks),
			fmt.Sprintf("%d", record.OpenIssues),
			fmt.Sprintf("%d", record.OpenPRs),
			record.LastCommit,
			record.Language,
			record.CICD,
			record.License,
			record.Description,
			record.Archived,
		}
		if err := w.Write(row); err != nil {
			return err
		}
	}

	return nil
}
