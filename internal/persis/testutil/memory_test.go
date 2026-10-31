package testutil_test

import (
	"testing"

	"github.com/rajputpriyansh444-commits/flow/internal/persis"
	"github.com/rajputpriyansh444-commits/flow/internal/persis/testutil"
)

func TestMemoryCollection(t *testing.T) {
	t.Parallel()

	factory := func(t *testing.T) (persis.Collection, persis.Collection) {
		t.Helper()
		backend := testutil.NewMemoryBackend()
		return backend.Collection("test"), backend.Collection("test")
	}
	testutil.RunCollectionContract(t, factory)
}
