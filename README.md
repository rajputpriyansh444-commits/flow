# Flow

A self-hosted, distributed task scheduler. Define workflows in YAML, schedule them with cron, run them on this machine or on worker machines, and watch every run live in the web UI.

## Quick start (Windows)

```powershell
.\scripts\build.ps1           # builds the UI and .local\bin\flow.exe
.\.local\bin\flow.exe start-all
```

Open http://127.0.0.1:8080, create the admin account, and run one of the example workflows.

## Write a workflow

Save a file in `%USERPROFILE%\.config\flow\dags\`, for example `backup.yaml`:

```yaml
schedule: "0 22 * * *"        # every day at 22:00
steps:
  - name: copy projects
    command: robocopy "D:\my all projects" "E:\backup" /MIR
  - name: done
    command: echo backup finished
```

It appears in the UI immediately. `flow start backup` runs it now, `flow status backup` shows the last run.

## Run jobs on more machines

On the main machine: `flow start-all --coordinator.host=0.0.0.0` (scheduler + web UI + coordinator on port 50055).
On each worker: `flow worker --worker.coordinators=<main-ip>:50055 --worker.labels=os=windows`.
Steps with `worker_selector: {os: windows}` are routed to matching workers.

## Project layout

| Path | What it is |
|---|---|
| `cmd/` | `flow` command entry point |
| `internal/spec`, `internal/ir` | YAML parsing and validated workflow model |
| `internal/runtime` | step execution, retries, lifecycle |
| `internal/service/scheduler` | cron scheduler and queue |
| `internal/service/coordinator`, `worker` | distributed execution over gRPC |
| `internal/service/frontend` | HTTP API, SSE, embedded UI |
| `internal/persis` | file-based storage (no database needed) |
| `ui/` | React + TypeScript web UI |

See `README_SCHEMA.md` for every YAML field.
