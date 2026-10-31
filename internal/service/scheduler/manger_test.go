package scheduler_test

import (
	"context"
	filedag "github.com/rajputpriyansh444-commits/flow/internal/persis/file/dag"
	"testing"

	"github.com/rajputpriyansh444-commits/flow/internal/test"
	"github.com/stretchr/testify/require"
)

func TestReadEntries(t *testing.T) {
	t.Run("InvalidDirectory", func(t *testing.T) {
		manager := filedag.NewFileEntryReader("invalid_directory", nil, false, "", "")
		err := manager.Init(context.Background())
		require.Error(t, err)
	})
	t.Run("InitAndDAGs", func(t *testing.T) {
		th := test.SetupScheduler(t)
		ctx := context.Background()

		err := th.EntryReader.Init(ctx)
		require.NoError(t, err)

		entries := th.EntryReader.Entries()
		require.NotEmpty(t, entries, "DAGs should not be empty")
	})
}
