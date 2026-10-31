package cmd

import (
	"github.com/rajputpriyansh444-commits/flow/internal/cmn/config"
	"github.com/rajputpriyansh444-commits/flow/internal/persis"
	"github.com/rajputpriyansh444-commits/flow/internal/persis/file"
)

func newDAGRepository(cfg *config.Config, storeCfg dagRepositoryConfig) (*persis.DAGRepository, error) {
	searchPaths := append([]string{}, storeCfg.SearchPaths...)
	if cfg.Paths.AltDAGsDir != "" {
		searchPaths = append(searchPaths, cfg.Paths.AltDAGsDir)
	}
	return file.NewDAGRepository(
		cfg,
		file.WithDAGFileCache(storeCfg.Cache),
		file.WithDAGSearchPaths(searchPaths),
		file.WithDAGSkipDirectoryCreation(storeCfg.SkipDirectoryCreation),
	)
}
