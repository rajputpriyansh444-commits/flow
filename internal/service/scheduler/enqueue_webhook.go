package scheduler

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/rajputpriyansh444-commits/flow/internal/cmn/logger"
	"github.com/rajputpriyansh444-commits/flow/internal/cmn/logger/tag"
	"github.com/rajputpriyansh444-commits/flow/internal/dagrun"
	"github.com/rajputpriyansh444-commits/flow/internal/intake"
	"github.com/rajputpriyansh444-commits/flow/internal/ir"
	"github.com/rajputpriyansh444-commits/flow/internal/persis"
	"github.com/rajputpriyansh444-commits/flow/internal/publicwebhook"
	queuedomain "github.com/rajputpriyansh444-commits/flow/internal/queue"
	"github.com/rajputpriyansh444-commits/flow/internal/spec"
)

// EnqueueWebhookRun enqueues a webhook-triggered run while preserving the same
// runtime-param semantics as direct webhook execution.
//
// The function is idempotent: if a run with the same ID already exists, it
// returns nil. An error wrapping [publicwebhook.ErrInvalidRun] means the DAG
// or the params cannot start a run.
func EnqueueWebhookRun(
	ctx context.Context,
	dagRunRepository *persis.DAGRunRepository,
	queueStore queuedomain.QueueStore,
	baseLogDir string,
	baseArtifactDir string,
	baseConfig string,
	workspaceBaseConfigDir string,
	definitionID string,
	dag *ir.DAG,
	runID string,
	params string,
	profileName string,
	now time.Time,
) error {
	return enqueueTriggeredRun(ctx, dagRunRepository, queueStore, baseLogDir, baseArtifactDir, baseConfig, workspaceBaseConfigDir,
		definitionID, dag, runID, params, profileName, now, ir.TriggerTypeWebhook, publicwebhook.ErrInvalidRun)
}

// enqueueTriggeredRun enqueues a run an outside event asks for, a webhook
// delivery or a file's arrival, with the params the event gives it. It is
// idempotent: a run with the same ID already existing is left alone. An
// error wrapping invalid means the DAG or the params cannot start a run.
func enqueueTriggeredRun(
	ctx context.Context,
	dagRunRepository *persis.DAGRunRepository,
	queueStore queuedomain.QueueStore,
	baseLogDir string,
	baseArtifactDir string,
	baseConfig string,
	workspaceBaseConfigDir string,
	definitionID string,
	dag *ir.DAG,
	runID string,
	params string,
	profileName string,
	now time.Time,
	triggerType ir.TriggerType,
	invalid error,
) error {
	dagRun := ir.NewDAGRunRef(dag.Name, runID)

	if _, err := dagRunRepository.FindAttempt(ctx, dagRun); err == nil {
		logger.Info(ctx, "Run already exists; skipping",
			tag.DAG(dag.Name),
			tag.RunID(runID),
			slog.String("trigger", triggerType.String()),
		)
		return nil
	} else if !errors.Is(err, dagrun.ErrDAGRunIDNotFound) {
		return fmt.Errorf("failed to check existing %s run: %w", triggerType, err)
	}

	fullDAG, err := rehydrateExecutionDAG(ctx, dag, params, baseConfig, workspaceBaseConfigDir)
	if err != nil {
		return fmt.Errorf("%w: %w", invalid, err)
	}
	if fullDAG == nil {
		return fmt.Errorf("failed to load full DAG for %s enqueue: DAG is nil", triggerType)
	}
	if len(fullDAG.BuildErrors) > 0 {
		return fmt.Errorf("%w: %w", invalid, errors.Join(fullDAG.BuildErrors...))
	}
	if err := spec.ValidateStartParams(fullDAG.DefaultParams, spec.StartParamInput{RawParams: params}); err != nil {
		return fmt.Errorf("%w: %w", invalid, err)
	}

	dagCopy := fullDAG.Clone()
	dagCopy.Location = ""

	_, err = intake.EnqueueRun(ctx, intake.QueueRequest{
		DAGRunRepository: dagRunRepository,
		QueueStore:       queueStore,
		DAG:              dagCopy,
		DAGRunID:         runID,
		LogBaseDir:       baseLogDir,
		ArtifactBaseDir:  baseArtifactDir,
		TriggerType:      triggerType,
		ProfileName:      profileName,
		DefinitionID:     definitionID,
		Now:              func() time.Time { return now },
	})
	if err != nil {
		return fmt.Errorf("failed to enqueue %s run: %w", triggerType, err)
	}

	logger.Info(ctx, "Run enqueued",
		tag.DAG(dag.Name),
		tag.RunID(runID),
		slog.String("trigger", triggerType.String()),
	)

	return nil
}
