package cmdutil

import "os/exec"

// StartParentExitWatcher is a no-op on Windows. Windows process cleanup uses
// process-tree handling in the platform-specific lifecycle adapter.
func StartParentExitWatcher(_ *exec.Cmd) (func(), error) {
	return func() {}, nil
}
