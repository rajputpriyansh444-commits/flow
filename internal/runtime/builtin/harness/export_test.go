package harness

import (
	"context"
	"os"

	"github.com/rajputpriyansh444-commits/flow/internal/ir"
)

func NewTestExecutorForTest(step ir.Step, prompt string, script string, workDir string) *harnessExecutor {
	return &harnessExecutor{
		step:    step,
		prompt:  prompt,
		script:  script,
		workDir: workDir,
	}
}

func NewTestExecutorWithProviderConfigsForTest(step ir.Step, prompt string, script string, workDir string, configs ...providerConfig) *harnessExecutor {
	return &harnessExecutor{
		step:    step,
		configs: configs,
		prompt:  prompt,
		script:  script,
		workDir: workDir,
	}
}

func NewTestProviderConfigForTest(name string, definition ir.HarnessDefinition, flags map[string]any) providerConfig {
	return providerConfig{
		name:       name,
		definition: &definition,
		flags:      flags,
	}
}

func (e *harnessExecutor) RunOnceForTest(ctx context.Context, cfg providerConfig) (*os.File, error) {
	return e.runOnce(ctx, cfg)
}

func SharedContainerHarnessEnvForTest(userEnv map[string]string) []string {
	return sharedContainerHarnessEnv(userEnv)
}
