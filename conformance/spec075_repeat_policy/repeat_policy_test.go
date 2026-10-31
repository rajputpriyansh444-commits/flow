package spec075_repeat_policy_test

import (
	"runtime"
	"testing"

	"github.com/rajputpriyansh444-commits/flow/conformance/harness"
)

// An unevaluable repeat condition stops repetition and follows continue_on.
func TestRuntimeRepeatConditionEvaluationErrorUnix(t *testing.T) {
	t.Parallel()
	if runtime.GOOS == "windows" {
		t.Skip("fixtures use POSIX shell snippets")
	}

	cases := []struct {
		file     string
		exitCode int
		ticks    string
	}{
		{"until_eval_error_fails.yaml", 1, "tick\n"},
		{"while_eval_error_fails.yaml", 1, "tick\n"},
		{"until_eval_error_mark_success.yaml", 0, "tick\ndependent\n"},
		{"while_eval_error_mark_success.yaml", 0, "tick\ndependent\n"},
	}
	for _, tc := range cases {
		t.Run(tc.file, func(t *testing.T) {
			t.Parallel()

			flow := harness.NewRunner(t)
			result := flow.Run("start", tc.file)
			result.ExpectExitCode(tc.exitCode)
			flow.ExpectFileContent("ticks.txt", tc.ticks)
		})
	}
}

func TestRuntimeRepeatDecisionUnix(t *testing.T) {
	t.Parallel()
	if runtime.GOOS == "windows" {
		t.Skip("fixtures use POSIX shell snippets")
	}

	cases := []struct {
		name  string
		file  string
		ticks string
	}{
		{
			name:  "until not-met repeats until the limit",
			file:  "until_not_met_respects_limit.yaml",
			ticks: "tick\ntick\ntick\n",
		},
		{
			name:  "until met stops after the first attempt",
			file:  "until_met_stops.yaml",
			ticks: "tick\n",
		},
		{
			name:  "while not-met stops after the first attempt",
			file:  "while_not_met_stops.yaml",
			ticks: "tick\n",
		},
		{
			name:  "until without a condition repeats until the attempt succeeds",
			file:  "until_unconditional_succeeds.yaml",
			ticks: "tick\ntick\n",
		},
		{
			name:  "while without a condition repeats while the attempt succeeds",
			file:  "while_unconditional_succeeds.yaml",
			ticks: "tick\ntick\n",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			flow := harness.NewRunner(t)
			result := flow.Run("start", tc.file)
			result.ExpectExitCode(0)
			flow.ExpectFileContent("ticks.txt", tc.ticks)
		})
	}
}
