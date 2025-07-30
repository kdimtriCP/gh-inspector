package formatter

import (
	"io"

	"github.com/yourname/gh-inspector/internal/github"
)

type Formatter interface {
	Format(writer io.Writer, metrics []*github.RepoMetrics) error
}
