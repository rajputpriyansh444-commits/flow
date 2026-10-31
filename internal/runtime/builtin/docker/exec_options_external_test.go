package docker_test

import (
	"testing"

	"github.com/rajputpriyansh444-commits/flow/internal/runtime/builtin/docker"
	"github.com/stretchr/testify/assert"
)

func TestExecCommandForTest_DirectBypassesShell(t *testing.T) {
	cmd := []string{"agent", "--prompt", "a && b"}

	got := docker.ExecCommandForTest(
		[]string{"/bin/sh", "-c"},
		cmd,
		docker.ExecOptions{Direct: true},
	)

	assert.Equal(t, cmd, got)
}

func TestExecCommandForTest_DefaultHonorsShell(t *testing.T) {
	got := docker.ExecCommandForTest(
		[]string{"/bin/sh"},
		[]string{"echo", "hello"},
		docker.ExecOptions{},
	)

	assert.Equal(t, []string{"/bin/sh", "-c", "echo hello"}, got)
}

func TestMergeEnvByKeyForTest_DeterministicOverride(t *testing.T) {
	got := docker.MergeEnvByKeyForTest(
		[]string{"PATH=/bin", "TOKEN=old", "KEEP=1"},
		[]string{"TOKEN=new", "BAD", "=ignored"},
		[]string{"PATH=/usr/bin"},
	)

	assert.Equal(t, []string{"PATH=/usr/bin", "TOKEN=new", "KEEP=1"}, got)
}
