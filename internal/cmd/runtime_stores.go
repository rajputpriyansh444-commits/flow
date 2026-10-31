package cmd

import (
	"path/filepath"

	"github.com/rajputpriyansh444-commits/flow/internal/build"
	"github.com/rajputpriyansh444-commits/flow/internal/persis"
	"github.com/rajputpriyansh444-commits/flow/internal/persis/file"
	filematerialization "github.com/rajputpriyansh444-commits/flow/internal/persis/file/materialization"
	"github.com/rajputpriyansh444-commits/flow/internal/profile"
	"github.com/rajputpriyansh444-commits/flow/internal/secret"
)

type executionStores struct {
	SecretStore          secret.Store
	ProfileStore         profile.Store
	MaterializationStore build.MaterializationStore
}

// runtimeStores creates the runtime store bundle for this command context.
func (c *Context) runtimeStores() executionStores {
	return executionStores{
		SecretStore:          file.NewSecretStore(c.Context, c.Config, c.backend.Collection(persis.CollectionSecrets)),
		ProfileStore:         file.NewProfileStore(c.Context, c.Config, c.backend.Collection(persis.CollectionProfiles)),
		MaterializationStore: filematerialization.New(filepath.Join(c.Config.Paths.DataDir, "materializations")),
	}
}
