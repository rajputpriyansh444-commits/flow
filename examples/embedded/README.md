# Embedded Flow Examples

These examples show how another Go application can import Flow as a library:

```go
import "github.com/rajputpriyansh444-commits/flow"
```

The embedded API is experimental. It is intended for applications that want to
run Flow DAGs in-process or dispatch them to existing Flow coordinators while
keeping Flow's CLI and server behavior unchanged.

## Local Execution

Run a DAG in the current process with file-backed state:

```sh
go run ./examples/embedded/local
```

The example creates a temporary Flow home directory, loads
`examples/embedded/local/workflow.yaml`, passes parameters, waits for completion,
and prints the final status.

## Custom Executor

Register a process-local executor and use it from DAG YAML:

```sh
go run ./examples/embedded/custom-executor
```

This is useful when the embedding application already has domain-specific Go
code and wants DAG steps to call that code without shelling out.

## Distributed Worker Execution

Run against an existing Flow coordinator:

```sh
FLOW_COORDINATORS=127.0.0.1:50055 go run ./examples/embedded/distributed
```

The example starts an embedded worker and dispatches a DAG through the
coordinator. Set `FLOW_COORDINATORS` to a comma-separated list when using more
than one coordinator. The sample opts into plaintext coordinator transport for
local development; production embedders should provide TLS files instead.
