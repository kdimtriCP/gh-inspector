package formatter

import (
	"encoding/csv"
	"io"

	"github.com/kdimtriCP/gh-inspector/internal/metrics"
)

type CSVFormatter struct{}

func NewCSVFormatter() *CSVFormatter {
	return &CSVFormatter{}
}

func (f *CSVFormatter) Format(writer io.Writer, metricsData []*metrics.Repository) error {
	w := csv.NewWriter(writer)
	defer w.Flush()

	headers := GetRecordHeaders()
	if err := w.Write(headers); err != nil {
		return err
	}

	for _, m := range metricsData {
		row := MetricsToRecord(m).Strings()
		if err := w.Write(row); err != nil {
			return err
		}
	}

	return nil
}
