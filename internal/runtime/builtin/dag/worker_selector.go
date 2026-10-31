package dag

import (
	"context"
	"fmt"
	"strings"

	"github.com/rajputpriyansh444-commits/flow/internal/cmn/config"
	"github.com/rajputpriyansh444-commits/flow/internal/ir"
	"github.com/rajputpriyansh444-commits/flow/internal/runtime/executor"
	"github.com/rajputpriyansh444-commits/flow/internal/spec"
)

func resolveChildRunParams(
	ctx context.Context,
	childDAG *ir.DAG,
	runParams executor.RunParams,
) (executor.RunParams, error) {
	if len(runParams.WorkerSelector) > 0 {
		return runParams, nil
	}
	if len(childDAG.WorkerSelector) == 0 {
		return runParams, nil
	}

	resolved, err := spec.ResolveRuntimeParams(ctx, childDAG, runParams.Params, spec.ResolveRuntimeParamsOptions{
		BaseConfig: config.GetConfig(ctx).Paths.BaseConfig,
	})
	if err != nil {
		return executor.RunParams{}, fmt.Errorf("resolve sub-DAG routing parameters: %w", err)
	}
	runParams.WorkerSelector = resolved.WorkerSelector
	runParams.Params = strings.Join(spec.QuoteRuntimeParams(resolved.Params, resolved.ParamDefs), " ")
	return runParams, nil
}
