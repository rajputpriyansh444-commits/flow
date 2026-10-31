package value_test

import "github.com/rajputpriyansh444-commits/flow/internal/cmn/value"

func testEnvScope(entries map[string]string) *value.EnvScope {
	scope := value.NewEnvScope(nil, false)
	for key, val := range entries {
		scope = scope.WithEntry(key, val, value.EnvSourceDAGEnv)
	}
	return scope
}
