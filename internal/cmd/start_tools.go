package cmd

import (
	"os"
	"strings"

	"github.com/rajputpriyansh444-commits/flow/internal/ir"
	flowtools "github.com/rajputpriyansh444-commits/flow/internal/tools"
	flowaqua "github.com/rajputpriyansh444-commits/flow/internal/tools/aqua"
)

func prepareDAGTools(ctx *Context, dag *ir.DAG) ([]string, error) {
	workDir := ""
	if dag != nil {
		workDir = dag.WorkingDir
	}
	return flowtools.PrepareDAG(ctx.Context, dag, flowaqua.New(), flowtools.InstallOptions{
		ToolsDir: ctx.Config.Paths.ToolsDir,
		DataDir:  ctx.Config.Paths.DataDir,
		WorkDir:  workDir,
	}, dagToolsBasePath(ctx))
}

func dagToolsBasePath(ctx *Context) string {
	if ctx != nil && ctx.Config != nil {
		for _, env := range ctx.Config.Core.BaseEnv.AsSlice() {
			key, value, ok := strings.Cut(env, "=")
			if ok && strings.EqualFold(key, "PATH") {
				return value
			}
		}
	}
	return os.Getenv("PATH")
}
