package formatter

import (
	"encoding/json"
	"io"

	"github.com/yourname/gh-inspector/internal/github"
)

type JSONFormatter struct {
	indent bool
}

func NewJSONFormatter(indent bool) *JSONFormatter {
	return &JSONFormatter{indent: indent}
}

func (f *JSONFormatter) Format(writer io.Writer, metrics []*github.RepoMetrics) error {
	var records []*Record
	for _, m := range metrics {
		records = append(records, MetricsToRecord(m))
	}

	encoder := json.NewEncoder(writer)
	if f.indent {
		encoder.SetIndent("", "  ")
	}
	return encoder.Encode(records)
}
