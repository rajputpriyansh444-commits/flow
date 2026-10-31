package testutil

import (
	"path/filepath"

	"github.com/rajputpriyansh444-commits/flow/internal/persis"
	filedagrun "github.com/rajputpriyansh444-commits/flow/internal/persis/file/dagrun"
)

// NewFileDAGRunRepository constructs a file-backed repository for tests.
func NewFileDAGRunRepository(
	baseDir string,
	options persis.DAGRunRepositoryOptions,
	storeOptions ...filedagrun.StoreOption,
) *persis.DAGRunRepository {
	return persis.NewDAGRunRepository(
		filedagrun.NewStore(baseDir, storeOptions...),
		filedagrun.NewWorkDirStore(filepath.Join(baseDir, ".dag-run-work"), baseDir),
		options,
	)
}

// NewFileDAGRunRepositoryWithEphemeral constructs a file-backed repository
// whose ephemeral runs live in ephemeralDir, outside baseDir.
func NewFileDAGRunRepositoryWithEphemeral(
	baseDir, ephemeralDir string,
	options persis.DAGRunRepositoryOptions,
	storeOptions ...filedagrun.StoreOption,
) *persis.DAGRunRepository {
	options.EphemeralStore = filedagrun.NewStore(ephemeralDir, storeOptions...)
	return NewFileDAGRunRepository(baseDir, options, storeOptions...)
}
