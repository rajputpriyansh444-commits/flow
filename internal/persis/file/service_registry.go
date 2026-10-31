package file

import (
	"github.com/rajputpriyansh444-commits/flow/internal/cmn/config"
	fileserviceregistry "github.com/rajputpriyansh444-commits/flow/internal/persis/file/serviceregistry"
	"github.com/rajputpriyansh444-commits/flow/internal/serviceregistry"
)

// NewServiceRegistry wires the file-backed service registry from application config.
func NewServiceRegistry(cfg *config.Config) serviceregistry.ServiceRegistry {
	return fileserviceregistry.New(cfg.Paths.ServiceRegistryDir)
}
