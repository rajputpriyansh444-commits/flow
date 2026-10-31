package config_test

import (
	"runtime"
	"strings"
	"testing"

	"github.com/rajputpriyansh444-commits/flow/internal/cmn/config"
	"github.com/stretchr/testify/require"
)

func TestLoadBaseEnv(t *testing.T) {
	windowsExpected := runtime.GOOS == "windows"
	unixExpected := runtime.GOOS != "windows"

	testCases := []struct {
		name     string
		expected bool
	}{
		{"TEST_VAR_BASE_ENV", false},
		{"FLOW_TEST_BASE_ENV", true},

		// XDG base-dir variables should pass through on Unix/macOS only.
		{"XDG_CONFIG_HOME", unixExpected},
		{"XDG_DATA_HOME", unixExpected},
		{"XDG_CACHE_HOME", unixExpected},
		{"XDG_STATE_HOME", unixExpected},

		// Windows-specific user profile and install roots should only pass
		// through on Windows builds.
		{"APPDATA", windowsExpected},
		{"LOCALAPPDATA", windowsExpected},
		{"USERNAME", windowsExpected},
		{"USERDOMAIN", windowsExpected},
		{"PROGRAMFILES", windowsExpected},
		{"PROGRAMFILES(X86)", windowsExpected},
		{"PROGRAMDATA", windowsExpected},

		// Docker daemon connection vars must pass through.
		{"DOCKER_HOST", true},
		{"DOCKER_TLS_VERIFY", true},
		{"DOCKER_CERT_PATH", true},
		{"DOCKER_API_VERSION", true},

		// Docker credentials must NOT leak through the global env.
		{"DOCKER_AUTH_CONFIG", false},

		// Kubernetes in-cluster discovery vars must pass through so kubectl and
		// client-go based tools do not fall back to localhost:8080.
		{"KUBERNETES_SERVICE_HOST", true},
		{"KUBERNETES_SERVICE_PORT", true},
		{"KUBERNETES_SERVICE_PORT_HTTPS", true},
	}

	for _, tc := range testCases {
		t.Setenv(tc.name, "value")
	}

	baseEnv := config.LoadBaseEnv()
	envSlice := baseEnv.AsSlice()
	envMap := parseEnvSlice(envSlice)

	for _, tc := range testCases {
		_, found := envMap[tc.name]
		require.Equal(t, tc.expected, found, "expected %s: %v", tc.name, tc.expected)
	}
}

func TestBaseEnv_AsSlice(t *testing.T) {
	t.Parallel()

	baseEnv := config.NewBaseEnv([]string{"A=1", "B=2"})
	require.Equal(t, []string{"A=1", "B=2"}, baseEnv.AsSlice())
}

func TestLoadBaseEnvWithExtras_ExactNamesAndPrefixes(t *testing.T) {
	t.Setenv("EXTRA_ALLOWED_NAME", "name-value")
	t.Setenv("EXTRA_ALLOWED_PREFIX_ONE", "prefix-value")
	t.Setenv("EXTRA_BLOCKED", "blocked-value")

	baseEnv := config.LoadBaseEnvWithExtras(
		[]string{" EXTRA_ALLOWED_NAME ", "EXTRA_ALLOWED_NAME", ""},
		[]string{" EXTRA_ALLOWED_PREFIX_", "EXTRA_ALLOWED_PREFIX_", ""},
	)
	envMap := parseEnvSlice(baseEnv.AsSlice())

	require.Equal(t, "name-value", envMap["EXTRA_ALLOWED_NAME"])
	require.Equal(t, "prefix-value", envMap["EXTRA_ALLOWED_PREFIX_ONE"])
	_, found := envMap["EXTRA_BLOCKED"]
	require.False(t, found)
}

func TestLoadBaseEnvWithExtras_NeverIncludesInternalTransport(t *testing.T) {
	t.Setenv("_FLOW_INTERNAL_OPENCODE_PASSWORD", "secret")
	baseEnv := config.LoadBaseEnvWithExtras(
		[]string{"_FLOW_INTERNAL_OPENCODE_PASSWORD"},
		[]string{"_FLOW_INTERNAL_"},
	)
	_, found := parseEnvSlice(baseEnv.AsSlice())["_FLOW_INTERNAL_OPENCODE_PASSWORD"]
	require.False(t, found)
}

func TestLoadBaseEnvWithExtras_PrefixMatchingHonorsPlatformCaseRules(t *testing.T) {
	t.Setenv("CASE_MATCH_ENV", "matched")

	baseEnv := config.LoadBaseEnvWithExtras(nil, []string{"case_"})
	envMap := parseEnvSlice(baseEnv.AsSlice())

	_, found := envMap["CASE_MATCH_ENV"]
	if runtime.GOOS == "windows" {
		require.True(t, found)
	} else {
		require.False(t, found)
	}
}

// Flow's child processes load their settings from the environment, so the
// base environment keeps Flow's credentials; the environment a workflow's
// processes start from leaves them out and keeps Flow's other settings.
func TestBaseEnvWithoutCredentials(t *testing.T) {
	credentials := []string{
		"FLOW_AUTH_BASIC_PASSWORD",
		"FLOW_AUTH_OIDC_CLIENT_SECRET",
		"FLOW_AUTH_TOKEN_SECRET",
		"FLOW_AUTH_BUILTIN_INITIAL_ADMIN_PASSWORD",
		"FLOW_SECRETS_VAULT_TOKEN",
		"FLOW_TUNNEL_TAILSCALE_AUTH_KEY",
		"FLOW_LICENSE",
		"FLOW_LICENSE_KEY",
		"FLOW_GITSYNC_AUTH_TOKEN",
		"FLOW_GITSYNC_AUTH_SSH_PASSPHRASE",
		"FLOW_ENCRYPTION_KEY",
		"FLOW_XLSX_PASSWORD",
	}
	for _, name := range credentials {
		t.Setenv(name, "credential")
	}
	t.Setenv("FLOW_HOME", "/srv/flow")

	baseEnv := config.LoadBaseEnv()
	inherited := parseEnvSlice(baseEnv.AsSlice())
	forSteps := parseEnvSlice(baseEnv.WithoutCredentials())
	for _, name := range credentials {
		require.Contains(t, inherited, name, "Flow's child processes need %s", name)
		require.NotContains(t, forSteps, name, "a workflow must not see %s", name)
		require.True(t, config.IsCredentialEnv(name), name)
	}
	require.Equal(t, "/srv/flow", forSteps["FLOW_HOME"])
	require.False(t, config.IsCredentialEnv("FLOW_HOME"))
}

func TestIsCredentialEnvHonorsPlatformCaseRules(t *testing.T) {
	t.Parallel()

	require.Equal(t, runtime.GOOS == "windows", config.IsCredentialEnv("flow_auth_basic_password"))
}

func parseEnvSlice(envSlice []string) map[string]string {
	envMap := make(map[string]string)
	for _, kv := range envSlice {
		key, value, found := strings.Cut(kv, "=")
		if found {
			if runtime.GOOS == "windows" {
				key = strings.ToUpper(key)
			}
			envMap[key] = value
		}
	}
	return envMap
}
