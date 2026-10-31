package cli_test

import (
	"testing"

	"github.com/rajputpriyansh444-commits/flow/conformance/harness"
	"github.com/stretchr/testify/require"
)

func TestExampleListsAvailableExamples(t *testing.T) {
	t.Parallel()

	flow := harness.NewRunner(t)
	result := flow.Run("example")
	result.ExpectExitCode(0)
	require.Contains(t, result.Stdout(), "parallel-steps")
}

func TestExampleShowsOneByID(t *testing.T) {
	t.Parallel()

	flow := harness.NewRunner(t)
	result := flow.Run("example", "1")
	result.ExpectExitCode(0)
	require.Contains(t, result.Stdout(), "Parallel Steps")
	require.Contains(t, result.Stdout(), "steps:")
}

func TestExampleRejectsUnknownID(t *testing.T) {
	t.Parallel()

	flow := harness.NewRunner(t)
	result := flow.Run("example", "99999")
	result.ExpectNonZeroExitCode()
}
