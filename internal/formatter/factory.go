package formatter

import (
	"fmt"
	"strings"
)

func New(format string) (Formatter, error) {
	switch strings.ToLower(format) {
	case FormatTable:
		return NewTableFormatter(), nil
	case FormatJSON:
		return NewJSONFormatter(true), nil
	case FormatJSONCompact:
		return NewJSONFormatter(false), nil
	case FormatCSV:
		return NewCSVFormatter(), nil
	default:
		return nil, fmt.Errorf("unsupported format: %s", format)
	}
}
