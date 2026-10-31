package cli_test

import (
	"strings"
	"testing"

	"github.com/rajputpriyansh444-commits/flow/conformance/harness"
	"github.com/stretchr/testify/require"
)

func TestVersionPrintsNonEmptyVersion(t *testing.T) {
	t.Parallel()

	flow := harness.NewRunner(t)
	result := flow.Run("version")
	result.ExpectExitCode(0)
	require.NotEmpty(t, strings.TrimSpace(result.Stdout()))
}
