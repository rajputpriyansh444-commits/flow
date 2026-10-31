package spec035_file_dependencies_test

import (
	"runtime"
	"testing"

	"github.com/rajputpriyansh444-commits/flow/conformance/harness"
)

func TestLocalFileDependencies(t *testing.T) {
	t.Parallel()
	if runtime.GOOS == "windows" {
		t.Skip("fixture uses a POSIX shell command")
	}

	flow := harness.NewRunner(t)
	result := flow.Run("start", "local.yaml")
	result.ExpectExitCode(0)
}

func TestFileDependenciesRejectValueReferences(t *testing.T) {
	t.Parallel()

	flow := harness.NewRunner(t)
	result := flow.Run("validate", "value-reference.yaml")
	result.ExpectExitCode(1)
	result.ExpectStderrContains("dependencies", "literal")
}
