package testutil

import (
	"github.com/rajputpriyansh444-commits/flow/internal/persis"
	filedag "github.com/rajputpriyansh444-commits/flow/internal/persis/file/dag"
)

// NewFileDAGRepository constructs a file-backed DAG repository for tests.
func NewFileDAGRepository(baseDir string, options ...filedag.Option) *persis.DAGRepository {
	return persis.NewDAGRepository(filedag.NewStore(baseDir, options...), persis.DAGRepositoryOptions{})
}
