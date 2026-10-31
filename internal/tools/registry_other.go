//go:build !unix

package tools

import (
	"fmt"
	"os"
)

func openLocalRegistry(path string) (*os.File, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	if !info.Mode().IsRegular() {
		return nil, fmt.Errorf("local registry %q is not a regular file", path)
	}
	return os.Open(path) //nolint:gosec
}
