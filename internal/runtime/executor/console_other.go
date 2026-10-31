//go:build !windows

package executor

import "golang.org/x/text/encoding"

// consoleEncoding returns nil: child processes have no console code page
// outside Windows.
func consoleEncoding() encoding.Encoding {
	return nil
}
