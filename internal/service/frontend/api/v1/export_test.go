package api

import (
	"context"
	"time"

	"github.com/rajputpriyansh444-commits/flow/internal/ir"
)

// Test exports expose internal helpers to external-package tests.
var (
	ExtractWebhookToken         = extractWebhookToken
	MarshalWebhookPayload       = marshalWebhookPayload
	MarshalWebhookHeaders       = marshalWebhookHeaders
	RequestedWebhookProfile     = requestedWebhookProfile
	IsWebhookTriggerPath        = isWebhookTriggerPath
	WithRawBody                 = withRawBody
	WithRequestHeaders          = withRequestHeaders
	BuildArtifactPreviewForTest = buildArtifactPreview
	ToDAGRunSummaryForTest      = toDAGRunSummary
)

// ArtifactTextPreviewMaxBytesForTest exposes the artifact preview size limit to external-package tests.
const ArtifactTextPreviewMaxBytesForTest = artifactTextPreviewMaxBytes

// NextRunProjectionForTest returns the API next-run projector for external-package tests.
func NextRunProjectionForTest(ctx context.Context, a *API) func(*ir.DAG, time.Time) time.Time {
	return a.nextRunProjection(ctx)
}
