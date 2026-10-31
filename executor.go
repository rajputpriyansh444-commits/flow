package flow

import (
	"context"
	"strings"

	"github.com/rajputpriyansh444-commits/flow/internal/executor/registry"
	"github.com/rajputpriyansh444-commits/flow/internal/ir"
	runtimeexec "github.com/rajputpriyansh444-commits/flow/internal/runtime/executor"
	"github.com/rajputpriyansh444-commits/flow/internal/spec"
)

// Step is the public alias for a Flow step passed to custom executors.
type Step = ir.Step

// Executor is implemented by custom step executors.
type Executor = runtimeexec.Executor

// ExecutorFactory creates an Executor for a loaded step.
type ExecutorFactory func(context.Context, Step) (Executor, error)

// StepValidator validates custom executor step configuration during DAG loading.
type StepValidator = registry.StepValidator

// ExecutorCapabilities declares which step fields a custom executor supports.
type ExecutorCapabilities = registry.ExecutorCapabilities

// ExecutorOption customizes custom executor registration.
type ExecutorOption func(*executorRegistration)

// executorRegistration holds optional custom executor metadata.
type executorRegistration struct {
	validator StepValidator
	caps      ExecutorCapabilities
}

// RegisterExecutor registers a custom executor type before engine or runtime use.
// It panics when name is empty, invalid, or factory is nil. Registration mutates
// global process state and must be completed before concurrent DAG execution.
func RegisterExecutor(name string, factory ExecutorFactory, opts ...ExecutorOption) {
	name = strings.TrimSpace(name)
	if name == "" {
		panic("flow: RegisterExecutor called with empty name")
	}
	if !spec.IsValidExecutorTypeName(name) {
		panic("flow: RegisterExecutor called with invalid name " + name)
	}
	if factory == nil {
		panic("flow: RegisterExecutor called with nil factory")
	}
	var registration executorRegistration
	for _, opt := range opts {
		opt(&registration)
	}
	runtimeexec.RegisterExecutor(
		name,
		runtimeexec.ExecutorFactory(factory),
		registration.validator,
		registration.caps,
	)
}

// UnregisterExecutor removes a custom executor type registered by RegisterExecutor.
// It is intended for tests and should not run concurrently with engine use.
func UnregisterExecutor(name string) {
	name = strings.TrimSpace(name)
	if name == "" {
		return
	}
	runtimeexec.UnregisterExecutor(name)
}

// WithStepValidator registers a validation function for the custom executor.
func WithStepValidator(validator StepValidator) ExecutorOption {
	return func(r *executorRegistration) {
		r.validator = validator
	}
}

// WithExecutorCapabilities registers supported step fields for the custom executor.
func WithExecutorCapabilities(caps ExecutorCapabilities) ExecutorOption {
	return func(r *executorRegistration) {
		r.caps = caps
	}
}
