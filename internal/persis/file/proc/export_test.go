package proc

import "os"

// WriteProcFileAtomicWithCreateTempForTest exposes the proc writer retry path to external tests.
func WriteProcFileAtomicWithCreateTempForTest(path string, data []byte, createTemp func(string, string) (*os.File, error)) error {
	return writeProcFileAtomicWithCreateTemp(path, data, createTemp)
}
