# Implementation Plan: Eventbus CLI Commands

**Feature**: `specs/002-eventbus-cli-commands`

## Technical Context

This feature adds a new `eventbus` command to `adxctl` with `start` and `stop` subcommands to manage a local NATS server for development. The implementation will involve creating new Cobra commands and leveraging existing internal packages for Docker and local process execution.

- **Command Structure**: New commands will be added in the `cmd/` directory, following the existing pattern (`cmd/eventbus.go`, `cmd/eventbus_start.go`, `cmd/eventbus_stop.go`). The root `eventbus` command will be registered in `cmd/root.go`.
- **Docker Runner**: The implementation will use the existing Docker client in `internal/agent/runner_docker.go`. The logic will be adapted to start, stop, and remove a `nats:latest` container named `adx-eventbus`.
- **Local Runner**: The `internal/agent/runner_local.go` will be extended or used as a reference to manage a `nats-server` process. It will handle starting the process in the background and storing its PID in a file (e.g., within `/tmp/adxctl/`) to allow the `stop` command to terminate it.
- **Configuration**: No new configuration in `config.yaml` is required, as this is a local development tool. The runner type is specified via a command-line flag.

## Constitution Check

The implementation plan is checked against the project constitution (v1.3.0).

- **I. Independently Testable Features**: The `eventbus` command is a self-contained feature. The runners can be tested independently.
- **II. Explicit & Tested Business Logic**: Business logic is minimal and resides clearly within the `start` and `stop` command handlers.
- **III. Simplicity Over Abstraction**: The solution is simple, using standard libraries (`os/exec`) and existing project dependencies (`moby/client`) without adding unnecessary layers.
- **IV. Containerized CI Testing**: The `docker` runner aligns perfectly with this principle. The `local` runner is a convenience for local dev and would be disabled or mocked in a CI environment.
- **V. Continuous Compilation**: This change will be part of the standard Go build process.
- **VI. Dependency Immutability**: The plan uses official, unmodified `nats:latest` images and relies on an existing `nats-server` binary; it does not modify any dependencies.

**Result**: The plan is fully compliant with the project constitution.

---

*This plan will be broken down into tasks in the next phase (`/speckit.tasks`).*
