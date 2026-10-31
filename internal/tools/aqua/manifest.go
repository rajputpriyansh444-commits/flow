package aqua

import (
	"encoding/json"
	"fmt"

	"github.com/rajputpriyansh444-commits/flow/internal/cmn/fileutil"
	"github.com/rajputpriyansh444-commits/flow/internal/tools"
)

func writeManifest(path string, manifest *tools.Manifest) error {
	data, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal tools manifest: %w", err)
	}
	if err := fileutil.WriteFileAtomic(path, data, 0o600); err != nil {
		return fmt.Errorf("write tools manifest: %w", err)
	}
	return nil
}
