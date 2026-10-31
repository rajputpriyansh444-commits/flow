package noop

import (
	"context"
	"io"
	"os"

	"github.com/rajputpriyansh444-commits/flow/internal/executor/registry"
	"github.com/rajputpriyansh444-commits/flow/internal/ir"
	"github.com/rajputpriyansh444-commits/flow/internal/runtime/executor"
)

type noopExecutor struct{}

func newNoop(_ context.Context, _ ir.Step) (executor.Executor, error) {
	return &noopExecutor{}, nil
}

func (*noopExecutor) SetStdout(_ io.Writer) {}

func (*noopExecutor) SetStderr(_ io.Writer) {}

func (*noopExecutor) Kill(_ os.Signal) error { return nil }

func (*noopExecutor) Run(_ context.Context) error { return nil }

func init() {
	executor.RegisterExecutor("noop", newNoop, nil, registry.ExecutorCapabilities{})
}
