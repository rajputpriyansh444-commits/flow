package command

import (
	"context"
	"path/filepath"
	"strings"

	"github.com/rajputpriyansh444-commits/flow/internal/runtime"
	flowtools "github.com/rajputpriyansh444-commits/flow/internal/tools"
)

func resolveRuntimeToolCommand(ctx context.Context, command string) string {
	if strings.TrimSpace(command) == "" || filepath.IsAbs(command) || strings.ContainsAny(command, `/\`) {
		return command
	}
	manifestPath := runtime.GetEnv(ctx).UserEnvsMap()[flowtools.EnvManifest]
	resolved, ok := resolveDeclaredToolCommand(manifestPath, command)
	if !ok {
		return command
	}
	return resolved
}

func resolveDeclaredToolCommand(manifestPath, command string) (string, bool) {
	if manifestPath == "" || command == "" {
		return "", false
	}
	manifest, err := flowtools.ReadManifest(manifestPath)
	if err != nil {
		return "", false
	}
	cmd, ok := manifest.Commands[command]
	if !ok || cmd.Path == "" {
		return "", false
	}
	return cmd.Path, true
}
