package formatter

import (
	"io"

	"github.com/kdimtriCP/gh-inspector/internal/metrics"
)

//go:generate mockgen -source=$GOFILE -destination=../mock/mock_formatter/mock_$GOFILE -package=mock_formatter

type Formatter interface {
	Format(writer io.Writer, metrics []*metrics.Repository) error
}
