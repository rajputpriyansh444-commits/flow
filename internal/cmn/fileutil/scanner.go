package fileutil

import "bufio"

const (
	ScannerInitialBufferSize = 64 * 1024
	ScannerMaxTokenSize      = 10 * 1024 * 1024
)

func ConfigureScanner(scanner *bufio.Scanner) {
	if scanner == nil {
		return
	}
	scanner.Buffer(make([]byte, ScannerInitialBufferSize), ScannerMaxTokenSize)
}
