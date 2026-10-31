package datapath

import (
	"context"
	"log/slog"

	"github.com/rajputpriyansh444-commits/flow/internal/cmn/logger"
	"github.com/rajputpriyansh444-commits/flow/internal/cmn/logger/tag"
	"github.com/itchyny/gojq"
)

// Select extracts a value from structured data using a jq-style path.
func Select(ctx context.Context, varName string, raw any, path string) (any, bool) {
	query, err := gojq.Parse(path)
	if err != nil {
		logger.Warn(ctx, "Failed to parse path in data",
			tag.Path(path),
			slog.String("var", varName),
			tag.Error(err))
		return nil, false
	}

	iter := query.RunWithContext(ctx, raw)
	v, ok := iter.Next()
	if !ok {
		return nil, false
	}

	if evalErr, ok := v.(error); ok {
		logger.Warn(ctx, "Failed to evaluate path in data",
			tag.Path(path),
			slog.String("var", varName),
			tag.Error(evalErr))
		return nil, false
	}

	return v, true
}
