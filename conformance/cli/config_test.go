package cli_test

import (
	"encoding/json"
	"path/filepath"
	"strings"
	"testing"

	"github.com/rajputpriyansh444-commits/flow/conformance/harness"
)

// TestConfigShowsResolvedPaths proves `flow config` resolves FLOW_HOME and
// FLOW_DAGS_DIR into the actual paths in effect, not just that it prints
// labeled lines: it asserts the DAGs-directory line names the project
// directory itself, and the data-directory line names FLOW_HOME's own
// "data" subdirectory, so a build that resolved either to the wrong path
// (or to another DAG-run's home) would fail this test.
func TestConfigShowsResolvedPaths(t *testing.T) {
	t.Parallel()

	flow := harness.NewRunner(t)
	home := filepath.Join(t.TempDir(), "flow")
	env := []string{
		"FLOW_HOME=" + home,
		"FLOW_DAGS_DIR=.",
	}

	result := flow.RunWithEnv(env, "config")
	result.ExpectExitCode(0)

	requireLineContains(t, result.Stdout(), "DAGs directory:", flow.ProjectPath("."))
	requireLineContains(t, result.Stdout(), "Data directory:", filepath.Join(home, "data"))
}

// TestConfigLoadsNamedTimeZone proves the tz setting accepts an IANA name on a
// host without a time zone database, such as Windows without Go installed.
// Pointing GOROOT at an empty directory removes the fallback to the build
// machine's Go installation.
func TestConfigLoadsNamedTimeZone(t *testing.T) {
	t.Parallel()

	flow := harness.NewRunner(t)
	env := []string{"FLOW_TZ=Asia/Tokyo", "GOROOT=" + t.TempDir(), "ZONEINFO="}

	flow.RunWithEnv(env, "config").ExpectExitCode(0)
}

// TestConfigValuesReportsSources proves `flow config --values` names the
// source of each value and never prints a secret: the port comes from the
// file, the host from the environment, the token secret is masked, and
// auth.mode reports the builtin mode Flow falls back to when none is set.
func TestConfigValuesReportsSources(t *testing.T) {
	t.Parallel()

	flow := harness.NewRunner(t)
	flow.WriteFile("config.yaml", "port: 9090\n")
	env := []string{
		"FLOW_HOME=" + filepath.Join(t.TempDir(), "flow"),
		"FLOW_HOST=0.0.0.0",
		"FLOW_AUTH_TOKEN_SECRET=do-not-print",
		// The harness sets auth.mode to none; clear it to see the fallback.
		"FLOW_AUTH_MODE=",
	}

	result := flow.RunWithEnv(env, "config", "--config", flow.ProjectPath("config.yaml"), "--values", "--format", "json")
	result.ExpectExitCode(0)

	if strings.Contains(result.Stdout(), "do-not-print") {
		t.Fatalf("secret printed; output:\n%s", result.Stdout())
	}
	var settings []struct {
		Key    string `json:"key"`
		Value  string `json:"value"`
		Source string `json:"source"`
	}
	if err := json.Unmarshal([]byte(result.Stdout()), &settings); err != nil {
		t.Fatalf("decoding output: %v; output:\n%s", err, result.Stdout())
	}
	got := make(map[string][2]string, len(settings))
	for _, s := range settings {
		got[s.Key] = [2]string{s.Value, s.Source}
	}

	want := map[string][2]string{
		"port":                      {"9090", "file"},
		"host":                      {"0.0.0.0", "env:FLOW_HOST"},
		"auth.builtin.token.secret": {"********", "env:FLOW_AUTH_TOKEN_SECRET"},
		"auth.mode":                 {"builtin", "default"},
	}
	for key, w := range want {
		if got[key] != w {
			t.Errorf("%s = %q, want %q", key, got[key], w)
		}
	}
}

// requireLineContains fails the test unless some line starts with label and
// its remainder, trimmed, equals value exactly. Matching the resolved path
// by exact equality (rather than substring) rejects a path that merely
// shares value as a prefix, such as a sibling directory named value+"-old".
func requireLineContains(t *testing.T, output, label, value string) {
	t.Helper()

	for line := range strings.SplitSeq(output, "\n") {
		rest, ok := strings.CutPrefix(line, label)
		if ok && strings.TrimSpace(rest) == value {
			return
		}
	}
	t.Fatalf("no line labeled %q has value %q; output:\n%s", label, value, output)
}
