//go:build windows

package api

import "errors"

func makeNamedPipe(string) error {
	return errors.New("named pipes via mkfifo are not supported on Windows")
}
