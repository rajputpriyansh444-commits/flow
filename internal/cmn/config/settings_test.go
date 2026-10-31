package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func loadSettings(t *testing.T, yamlContent string) map[string]Setting {
	t.Helper()

	configFile := filepath.Join(t.TempDir(), "config.yaml")
	require.NoError(t, os.WriteFile(configFile, []byte(yamlContent), 0600))

	loader := NewConfigLoader(viper.New(), WithAppHomeDir(t.TempDir()), WithConfigFile(configFile))
	cfg, err := loader.Load()
	require.NoError(t, err)

	settings := make(map[string]Setting)
	for _, s := range loader.Settings(cfg) {
		_, dup := settings[s.Key]
		require.False(t, dup, "key %s is listed twice", s.Key)
		settings[s.Key] = s
	}
	return settings
}

func TestSettingsSource(t *testing.T) {
	t.Setenv("FLOW_PORT", "")
	t.Setenv("FLOW_HOST", "0.0.0.0")
	t.Setenv("FLOW_COORDINATOR_PORT", "")

	settings := loadSettings(t, `
coordinator:
  port: 50060
`)

	assert.Equal(t, Setting{Key: "port", Value: "8080", Source: SourceDefault}, settings["port"])
	assert.Equal(t, Setting{Key: "host", Value: "0.0.0.0", Source: "env:FLOW_HOST"}, settings["host"])
	assert.Equal(t, Setting{Key: "coordinator.port", Value: "50060", Source: SourceFile}, settings["coordinator.port"])
}

func TestSettingsEnvOverridesFile(t *testing.T) {
	t.Setenv("FLOW_PORT", "9000")

	settings := loadSettings(t, "port: 8090\n")

	assert.Equal(t, Setting{Key: "port", Value: "9000", Source: "env:FLOW_PORT"}, settings["port"])
}

// A key bound to several variables reports the one Viper reads: the first
// in the binding table that is set.
func TestSettingsEnvAlias(t *testing.T) {
	t.Run("Legacy", func(t *testing.T) {
		t.Setenv("FLOW_UI_NAVBAR_COLOR", "")
		t.Setenv("FLOW_NAVBAR_COLOR", "red")

		settings := loadSettings(t, "")
		assert.Equal(t, Setting{Key: "ui.navbar_color", Value: "red", Source: "env:FLOW_NAVBAR_COLOR"}, settings["ui.navbar_color"])
	})

	t.Run("CanonicalWins", func(t *testing.T) {
		t.Setenv("FLOW_UI_NAVBAR_COLOR", "blue")
		t.Setenv("FLOW_NAVBAR_COLOR", "red")

		settings := loadSettings(t, "")
		assert.Equal(t, Setting{Key: "ui.navbar_color", Value: "blue", Source: "env:FLOW_UI_NAVBAR_COLOR"}, settings["ui.navbar_color"])
	})
}

// Viper's automatic FLOW_<KEY> lookup runs before the bound variables.
func TestSettingsAutomaticEnv(t *testing.T) {
	t.Setenv("FLOW_SERVER_METRICS", "private")
	t.Setenv("FLOW_METRICS", "public")

	settings := loadSettings(t, "")
	assert.Equal(t, Setting{Key: "metrics", Value: "public", Source: "env:FLOW_METRICS"}, settings["metrics"])
}

func TestSettingsListsAPIBasePath(t *testing.T) {
	t.Setenv("FLOW_API_BASE_PATH", "")

	settings := loadSettings(t, "api_base_path: /custom/api\n")
	assert.Equal(t, Setting{Key: "api_base_path", Value: "/custom/api", Source: SourceFile}, settings["api_base_path"])
}

// With no auth.mode set, Load runs the server in builtin mode; the report
// says so rather than showing an empty value.
func TestSettingsAuthMode(t *testing.T) {
	t.Setenv("FLOW_AUTH_MODE", "")

	t.Run("Default", func(t *testing.T) {
		settings := loadSettings(t, "")
		assert.Equal(t, Setting{Key: "auth.mode", Value: "builtin", Source: SourceDefault}, settings["auth.mode"])
	})

	t.Run("File", func(t *testing.T) {
		settings := loadSettings(t, "auth:\n  mode: none\n")
		assert.Equal(t, Setting{Key: "auth.mode", Value: "none", Source: SourceFile}, settings["auth.mode"])
	})
}

func TestSettingsMasksSecrets(t *testing.T) {
	t.Setenv("FLOW_AUTH_TOKEN_SECRET", "")
	t.Setenv("FLOW_AUTH_BASIC_PASSWORD", "")
	t.Setenv("FLOW_LICENSE_KEY", "env-license")

	settings := loadSettings(t, `
auth:
  builtin:
    token:
      secret: file-secret
`)

	assert.Equal(t, Setting{Key: "auth.builtin.token.secret", Value: MaskedValue, Source: SourceFile}, settings["auth.builtin.token.secret"])
	assert.Equal(t, Setting{Key: "license.key", Value: MaskedValue, Source: "env:FLOW_LICENSE_KEY"}, settings["license.key"])
	assert.Equal(t, Setting{Key: "auth.basic.password", Value: "", Source: SourceDefault}, settings["auth.basic.password"])
}

func TestSettingsFormatsLists(t *testing.T) {
	t.Setenv("FLOW_CORS_ALLOWED_ORIGINS", "")

	settings := loadSettings(t, `
cors_allowed_origins:
  - https://a.example
  - https://b.example
`)

	assert.Equal(t, "https://a.example,https://b.example", settings["cors_allowed_origins"].Value)
}
