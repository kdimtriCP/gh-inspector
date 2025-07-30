package formatter

import (
	"fmt"
	"io"
	"text/tabwriter"

	"github.com/yourname/gh-inspector/internal/github"
)

type TableFormatter struct{}

func NewTableFormatter() *TableFormatter {
	return &TableFormatter{}
}

func (f *TableFormatter) Format(writer io.Writer, metrics []*github.RepoMetrics) error {
	w := tabwriter.NewWriter(writer, 0, 0, 3, ' ', 0)
	_, err := fmt.Fprintln(w, "REPOSITORY\tSTARS\tFORKS\tOPEN ISSUES\tOPEN PRS\tLAST COMMIT\tLANGUAGE\tCI/CD\tLICENSE")
	if err != nil {
		return err
	}
	_, err = fmt.Fprintln(w, "----------\t-----\t-----\t-----------\t--------\t-----------\t--------\t-----\t-------")
	if err != nil {
		return err
	}

	for _, m := range metrics {
		record := MetricsToRecord(m)
		_, err = fmt.Fprintf(w, "%s\t%d\t%d\t%d\t%d\t%s\t%s\t%s\t%s\n",
			record.Repository,
			record.Stars,
			record.Forks,
			record.OpenIssues,
			record.OpenPRs,
			record.LastCommit,
			record.Language,
			record.CICD,
			record.License,
		)
		if err != nil {
			return err
		}
	}

	return w.Flush()
}
