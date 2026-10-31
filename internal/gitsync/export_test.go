package gitsync

func SafeReadFileWithinBaseForTest(baseDir, targetPath string) ([]byte, error) {
	return safeReadFileWithinBase(baseDir, targetPath)
}

func SafeWriteFileWithinBaseForTest(baseDir, targetPath string, content []byte) error {
	return safeWriteFileWithinBase(baseDir, targetPath, content, 0600)
}
