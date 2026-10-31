package api

import (
	"context"
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/rajputpriyansh444-commits/flow/api/v1"
	"github.com/rajputpriyansh444-commits/flow/internal/audit"
	"github.com/rajputpriyansh444-commits/flow/internal/browserhost"
	"github.com/rajputpriyansh444-commits/flow/internal/cmn/config"
	"github.com/rajputpriyansh444-commits/flow/internal/cmn/replaycache"
)

func (a *API) ClearDAGBrowserCache(ctx context.Context, request api.ClearDAGBrowserCacheRequestObject) (api.ClearDAGBrowserCacheResponseObject, error) {
	if err := a.isAllowed(config.PermissionRunDAGs); err != nil {
		return nil, err
	}
	dag, err := a.dagRepository.GetMetadata(ctx, request.FileName)
	if err != nil {
		return nil, &Error{
			HTTPStatus: http.StatusNotFound,
			Code:       api.ErrorCodeNotFound,
			Message:    fmt.Sprintf("DAG %s not found", request.FileName),
		}
	}
	if err := a.requireExecuteForWorkspace(ctx, dagWorkspaceName(dag)); err != nil {
		return nil, err
	}

	// Browser steps key their cache by the DAG name, which can differ from
	// the file name.
	step := valueOf(request.Params.Step)
	cache := replaycache.New(filepath.Join(a.config.Paths.DataDir, browserhost.DataDirName))
	steps, err := cache.Clear(dag.Name, step)
	if err != nil {
		return nil, fmt.Errorf("failed to clear browser replay cache: %w", err)
	}
	if steps == nil {
		steps = []string{}
	}

	a.logAudit(ctx, audit.CategoryDAG, "dag_browser_cache_clear", map[string]any{
		"dag_name": request.FileName,
		"step":     step,
		"steps":    steps,
	})
	return api.ClearDAGBrowserCache200JSONResponse{Steps: steps}, nil
}
