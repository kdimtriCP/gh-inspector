package formatter

import (
	"encoding/json"
	"io"

	"github.com/kdimtriCP/gh-inspector/internal/metrics"
)

type JSONFormatter struct {
	indent bool
}

func NewJSONFormatter(indent bool) *JSONFormatter {
	return &JSONFormatter{indent: indent}
}

func (f *JSONFormatter) Format(writer io.Writer, metricsData []*metrics.Repository) error {
	var records []*Record
	for _, m := range metricsData {
		records = append(records, MetricsToRecord(m))
	}

	encoder := json.NewEncoder(writer)
	if f.indent {
		encoder.SetIndent("", "  ")
	}
	return encoder.Encode(records)
}
