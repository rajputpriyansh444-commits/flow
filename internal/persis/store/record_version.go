package store

import (
	"context"

	"github.com/rajputpriyansh444-commits/flow/internal/persis"
)

type recordVersionCollection interface {
	RecordVersion(ctx context.Context, id string) (string, error)
}

func collectionRecordVersion(ctx context.Context, col persis.Collection, id string) (string, bool, error) {
	versioned, ok := col.(recordVersionCollection)
	if !ok {
		return "", false, nil
	}
	version, err := versioned.RecordVersion(ctx, id)
	return version, true, err
}
