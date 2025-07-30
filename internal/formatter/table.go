package formatter

import (
	"io"

	"github.com/olekukonko/tablewriter"
	"github.com/yourname/gh-inspector/internal/github"
)

type TableFormatter struct{}

func NewTableFormatter() *TableFormatter {
	return &TableFormatter{}
}

func (f *TableFormatter) Format(writer io.Writer, metrics []*github.RepoMetrics) error {
	table := tablewriter.NewWriter(writer)
	table.SetHeader(GetRecordHeaders())

	table.SetBorder(false)
	table.SetCenterSeparator("")
	table.SetColumnSeparator("")
	table.SetRowSeparator("")
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_LEFT)

	for _, m := range metrics {
		raw := MetricsToRecord(m).Strings()
		table.Append(raw)
	}

	table.Render()
	return nil
}
