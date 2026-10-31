package spec012_step_outputs_test

import (
	"testing"

	"github.com/rajputpriyansh444-commits/flow/conformance/harness"
)

func TestValidateStepOutputs(t *testing.T) {
	t.Parallel()

	validCases := []string{
		"valid_string_output.yaml",
		"valid_string_output_with_heredoc_marker.yaml",
		"valid_json_output.yaml",
	}
	for _, file := range validCases {
		t.Run(file, func(t *testing.T) {
			t.Parallel()

			flow := harness.NewRunner(t)
			result := flow.Run("validate", file)
			result.ExpectExitCode(0)
			result.ExpectStdout("")
			flow.ExpectNoFile("resolved-string.txt")
			flow.ExpectNoFile("resolved-json.txt")
		})
	}

	invalidCases := []struct {
		file        string
		stderrParts []string
	}{
		{file: "invalid_outputs_null.yaml", stderrParts: []string{"outputs", "non-empty sequence"}},
		{file: "invalid_outputs_empty.yaml", stderrParts: []string{"outputs", "non-empty sequence"}},
		{file: "invalid_outputs_missing_id.yaml", stderrParts: []string{"outputs", "must define id"}},
		{file: "invalid_outputs_missing_name.yaml", stderrParts: []string{"outputs", "name is required"}},
		{file: "invalid_outputs_invalid_name.yaml", stderrParts: []string{"outputs", "name must match"}},
		{file: "invalid_outputs_duplicate.yaml", stderrParts: []string{"outputs", "duplicate output name"}},
		{file: "invalid_outputs_unknown_field.yaml", stderrParts: []string{"outputs", "unknown field"}},
		{file: "invalid_outputs_invalid_type.yaml", stderrParts: []string{"outputs", "type must be"}},
	}
	for _, tc := range invalidCases {
		t.Run(tc.file, func(t *testing.T) {
			t.Parallel()

			flow := harness.NewRunner(t)
			result := flow.Run("validate", tc.file)
			result.ExpectExitCode(1)
			result.ExpectStdout("")
			result.ExpectStderrContains(tc.stderrParts...)
			result.ExpectStderrNotContains("Usage:")
			flow.ExpectNoFile("executed.txt")
		})
	}
}

func TestRuntimeStepOutputs(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name    string
		file    string
		output  string
		content string
	}{
		{
			name:    "string output",
			file:    "valid_string_output.yaml",
			output:  "resolved-string.txt",
			content: "v1.2.3\n",
		},
		{
			name:    "string output containing heredoc marker",
			file:    "valid_string_output_with_heredoc_marker.yaml",
			output:  "resolved-string-marker.txt",
			content: "prefix<<suffix\n",
		},
		{
			name:    "json output preserves emitted text",
			file:    "valid_json_output.yaml",
			output:  "resolved-json.txt",
			content: "{\"image\":\"api\",\"tag\":\"v1.2.3\"}\n",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			flow := harness.NewRunner(t)
			result := flow.Run("start", tc.file)
			result.ExpectExitCode(0)
			flow.ExpectFileContent(tc.output, tc.content)
		})
	}
}

func TestRuntimeStepOutputFailures(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name        string
		file        string
		stderrParts []string
		absentFile  string
	}{
		{
			name:        "stdout is not a step output",
			file:        "stdout_is_not_output.yaml",
			stderrParts: []string{"image_tag", "not emitted"},
			absentFile:  "should-not-run.txt",
		},
		{
			name:        "undeclared output write fails",
			file:        "undeclared_output_write.yaml",
			stderrParts: []string{"undeclared"},
			absentFile:  "unexpected-success.txt",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			flow := harness.NewRunner(t)
			result := flow.Run("start", tc.file)
			result.ExpectExitCode(1)
			result.ExpectStderrContains(tc.stderrParts...)
			flow.ExpectNoFile(tc.absentFile)
		})
	}
}
