package formatter

import (
	"io"

	"github.com/kdimtriCP/gh-inspector/internal/metrics"
)

type Formatter interface {
	Format(writer io.Writer, metrics []*metrics.Repository) error
}
