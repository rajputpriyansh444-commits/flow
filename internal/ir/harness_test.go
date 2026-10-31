package ir

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuiltinCLIHarnessProviderNamesSorted(t *testing.T) {
	assert.Equal(t, []string{
		"aider",
		"amp",
		"claude",
		"cline",
		"codex",
		"copilot",
		"cursor",
		"deepseek",
		"droid",
		"gemini",
		"goose",
		"kilo",
		"kiro",
		"opencode",
		"pi",
		"qwen",
	}, BuiltinCLIHarnessProviderNames())
}

func TestIsBuiltinCLIHarnessProvider(t *testing.T) {
	assert.True(t, IsBuiltinCLIHarnessProvider("codex"))
	assert.False(t, IsBuiltinCLIHarnessProvider("builtin"))
}
