# Research: Eventbus CLI Implementation

This document outlines the decisions made for implementing the `eventbus` CLI commands.

## 1. Local Process Management (`local` runner)

- **Decision**: Use the standard `os/exec` package to start the `nats-server` process. The `Cmd.Start()` function will be used to run the process in the background. A PID file will be created in a temporary directory (e.g., `/tmp/adxctl/eventbus.pid`) to store the process ID.
- **Rationale**: This is the standard Go approach for managing child processes. A PID file is a simple and reliable mechanism for the `stop` command to identify and terminate the correct process, ensuring `adxctl` doesn't interfere with other `nats-server` instances not started by it.
- **Alternatives Considered**:
  - **In-memory PID storage**: Storing the PID in a variable within the `adxctl` process is not viable because `start` and `stop` are separate command invocations.
  - **System-level process management (e.g., systemd)**: Overly complex for a local development tool and not portable across operating systems.

## 2. Docker Container Management (`docker` runner)

- **Decision**: Utilize the existing `github.com/moby/moby/client` library, which is already a project dependency. The implementation will mirror the patterns found in `internal/agent/runner_docker.go` for creating a Docker client, pulling the `nats:latest` image, creating a container with the name `adx-eventbus`, and starting it.
- **Rationale**: Reusing the existing, familiar Docker client library is consistent with project conventions and avoids adding a new dependency. It provides all the necessary functionality for container lifecycle management.
- **Alternatives Considered**:
  - **Shelling out to `docker` CLI**: Directly calling `docker run`, `docker stop`, etc., using `os/exec` would work but is less robust, harder to manage, and couples the implementation to the presence of the `docker` CLI executable, whereas the Moby client interacts with the Docker daemon API directly.

## 3. State Management

- **Decision**: The state (PID for `local` runner, container name for `docker` runner) will be managed via files in a temporary project directory (`/tmp/adxctl/`).
- **Rationale**: This provides a simple, stateless approach for the CLI, where each command execution can read the state left by the previous one. This is suitable for a local development utility.
- **Alternatives Considered**: None, as this is the most straightforward approach that meets the requirements.
