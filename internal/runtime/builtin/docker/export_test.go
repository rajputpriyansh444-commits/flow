package docker

func ExecCommandForTest(shell, cmd []string, opts ExecOptions) []string {
	return execCommand(shell, cmd, opts)
}

func MergeEnvByKeyForTest(layers ...[]string) []string {
	return mergeEnvByKey(layers...)
}
