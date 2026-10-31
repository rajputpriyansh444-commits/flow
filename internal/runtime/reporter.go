package runtime

import (
	"context"
	"io"
	"os"

	"github.com/rajputpriyansh444-commits/flow/internal/ir"
	"github.com/rajputpriyansh444-commits/flow/internal/runctx"
)

// StatusPusher reports DAG run status outside the current execution process.
type StatusPusher interface {
	Push(ctx context.Context, status ir.DAGRunStatus) error
}

// AttemptRejected marks a status push failure caused by a non-authoritative attempt.
type AttemptRejected interface {
	error
	AttemptRejectedReason() string
}

// SchedulerLogStreamer streams a completed scheduler log.
type SchedulerLogStreamer interface {
	runctx.LogWriterFactory
	NewSchedulerLogWriter(ctx context.Context, localFile *os.File) io.WriteCloser
	StreamSchedulerLog(ctx context.Context, logFilePath string) error
}

// ArtifactFinalizer persists artifacts before terminal status is reported.
type ArtifactFinalizer interface {
	Finalize(ctx context.Context, attemptID, dir string) error
}
