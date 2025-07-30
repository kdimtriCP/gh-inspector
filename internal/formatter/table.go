package formatter

import (
	"io"

	"github.com/olekukonko/tablewriter"

	"github.com/kdimtriCP/gh-inspector/internal/metrics"
)

type TableFormatter struct{}

func NewTableFormatter() *TableFormatter {
	return &TableFormatter{}
}

func (f *TableFormatter) Format(writer io.Writer, metricsData []*metrics.Repository) error {
	table := tablewriter.NewWriter(writer)
	table.SetHeader(GetRecordHeaders())

	table.SetBorder(false)
	table.SetCenterSeparator("")
	table.SetColumnSeparator("")
	table.SetRowSeparator("")
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_LEFT)

	for _, m := range metricsData {
		raw := MetricsToRecord(m).Strings()
		table.Append(raw)
	}

	table.Render()
	return nil
}
