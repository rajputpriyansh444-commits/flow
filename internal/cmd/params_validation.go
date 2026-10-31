package cmd

import (
	"github.com/rajputpriyansh444-commits/flow/internal/cmn/stringutil"
	"github.com/rajputpriyansh444-commits/flow/internal/ir"
	"github.com/rajputpriyansh444-commits/flow/internal/spec"
)

func validateStartArgumentSeparator(ctx *Context, args []string) error {
	return spec.ValidateStartArgs(ctx.Command.ArgsLenAtDash() != -1, args)
}

// rawParams is the params text resolved by loadDAGWithParams ("--params" flag
// value or piped stdin content); it is ignored when params come after "--".
func validateStartPositionalParamCount(ctx *Context, args []string, dag *ir.DAG, rawParams string) error {
	return spec.ValidateStartParams(dag.DefaultParams, buildStartValidationInput(ctx, args, rawParams))
}

func buildStartValidationInput(ctx *Context, args []string, rawParams string) spec.StartParamInput {
	if argsLenAtDash := ctx.Command.ArgsLenAtDash(); argsLenAtDash != -1 {
		if argsLenAtDash >= len(args) {
			return spec.StartParamInput{}
		}
		return spec.StartParamInput{DashArgs: quoteStartDashArgs(args[argsLenAtDash:])}
	}

	if ctx.Command.Flags().Changed("params") {
		rawParams = stringutil.RemoveQuotes(rawParams)
	}
	return spec.StartParamInput{RawParams: rawParams}
}
