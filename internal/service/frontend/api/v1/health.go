package api

import (
	"context"
	"time"

	"github.com/rajputpriyansh444-commits/flow/api/v1"
	"github.com/rajputpriyansh444-commits/flow/internal/cmn/config"
	"github.com/rajputpriyansh444-commits/flow/internal/cmn/stringutil"
	"github.com/rajputpriyansh444-commits/flow/internal/service/frontend/metrics"
)

func (a *API) GetHealthStatus(_ context.Context, _ api.GetHealthStatusRequestObject) (api.GetHealthStatusResponseObject, error) {
	return &api.GetHealthStatus200JSONResponse{
		Status:    api.HealthResponseStatusHealthy,
		Version:   config.Version,
		Uptime:    int(metrics.GetUptime()),
		Timestamp: stringutil.FormatTime(time.Now()),
	}, nil
}
