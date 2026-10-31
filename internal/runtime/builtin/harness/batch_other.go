//go:build !windows

package harness

import (
	"errors"
	"os/exec"
)

func setBatchCommandLine(*exec.Cmd, string) error {
	return errors.New("harness: batch files run only on Windows")
}
