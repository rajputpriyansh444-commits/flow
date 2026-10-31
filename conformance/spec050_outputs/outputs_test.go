package spec050_outputs_test

import (
	"testing"

	"github.com/rajputpriyansh444-commits/flow/conformance/harness"
)

func TestOutputsWrite(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct{ file, want string }{
		{"write_and_read.yaml", "greeting=hello count=3 strict=hello"},
		{"write_dynamic_value.yaml", "from-env"},
		{"write_list_value.yaml", `["a","b"]`},
	} {
		t.Run(tc.file, func(t *testing.T) {
			t.Parallel()

			flow := harness.NewRunner(t)
			result := flow.Run("start", tc.file)
			result.ExpectExitCode(0)
			flow.ExpectFileContent("result.out", tc.want)
		})
	}
}

func TestOutputsValidation(t *testing.T) {
	t.Parallel()

	for _, file := range []string{
		"missing_values.yaml",
		"empty_values.yaml",
		"values_not_object.yaml",
		"unsupported_field.yaml",
		"empty_key.yaml",
	} {
		t.Run(file, func(t *testing.T) {
			t.Parallel()

			flow := harness.NewRunner(t)
			result := flow.Run("validate", file)
			result.ExpectNonZeroExitCode()
		})
	}
}
