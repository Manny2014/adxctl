# Feature Specification: Eventbus CLI Commands

**Feature Branch**: `[002-eventbus-cli-commands]`

**Created**: 2026-09-22

**Status**: Draft

**Input**: User description: "Implement CLI commants for "eventbus" that allows that user to start a nats eventbus either locally or running in docker. This is purely to facilitate local development and not meant to be a scalable production solution. Example Commands: adxclt eventbus start|stop --runner docker|local if docker, run standard docker nats image and if local, assume that the nats cli is available"

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Run eventbus in Docker (Priority: P1)

As a developer, I want to start a local NATS event bus in a Docker container so that I can develop and test features that rely on an event bus without having to install NATS directly on my machine.

**Why this priority**: This is the most common use case, providing a clean, isolated environment that works for any developer with Docker installed, minimizing local setup friction.

**Independent Test**: The `adxctl eventbus start --runner docker` command can be executed, and a NATS server will be verifiable at `nats://localhost:4222`. This delivers the core value of having a running event bus.

**Acceptance Scenarios**:

1. **Given** Docker is running and no container named `adx-eventbus` exists, **When** the user runs `adxctl eventbus start --runner docker`, **Then** a NATS Docker container is started, and the CLI confirms that the event bus is running.
2. **Given** a running `adx-eventbus` container, **When** the user runs `adxctl eventbus stop`, **Then** the container is stopped and removed, and the CLI confirms the shutdown.
3. **Given** no Docker daemon is running, **When** the user runs `adxctl eventbus start --runner docker`, **Then** the CLI returns an informative error that Docker is not available.

### User Story 2 - Run eventbus locally (Priority: P2)

As a developer with the NATS CLI already installed, I want to start a local NATS event bus using a direct process so that I can leverage my existing setup for speed and convenience.

**Why this priority**: This provides a lightweight alternative for developers who already have the necessary tools installed, but it's secondary because it relies on a pre-existing local dependency.

**Independent Test**: The `adxctl eventbus start --runner local` command can be executed, and a NATS server process will be started and verifiable at `nats://localhost:4222`.

**Acceptance Scenarios**:

1. **Given** the `nats-server` executable is in the system's PATH, **When** the user runs `adxctl eventbus start --runner local`, **Then** a local `nats-server` process is started in the background, and the CLI confirms it is running.
2. **Given** a `nats-server` process was started by the CLI, **When** the user runs `adxctl eventbus stop`, **Then** the background process is terminated, and the CLI confirms the shutdown.
3. **Given** the `nats-server` executable is not in the system's PATH, **When** the user runs `adxctl eventbus start --runner local`, **Then** the CLI returns an informative error that the `nats-server` command could not be found.

### Edge Cases

- What happens if a process is already listening on port 4222 when the `start` command is run? The command should fail with a clear error message.
- How does the system handle a `stop` command when no event bus was started by `adxctl`? It should inform the user that no active `adxctl`-managed event bus was found.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: The CLI MUST provide a top-level `eventbus` command.
- **FR-002**: The `eventbus` command MUST have `start` and `stop` subcommands.
- **FR-003**: The `start` command MUST include a `--runner` flag that accepts `docker` or `local` as values.
- **FR-004**: The default value for the `--runner` flag MUST be `docker`.
- **FR-005**: When run with `--runner docker`, the `start` command MUST attempt to run the official `nats:latest` Docker image, mapping port 4222.
- **FR-006**: When run with `--runner local`, the `start` command MUST attempt to execute the `nats-server` command as a background process.
- **FR-007**: The `stop` command MUST gracefully terminate the corresponding NATS instance (Docker container or local process) that was initiated by the `start` command.
- **FR-008**: The CLI MUST provide clear, human-readable feedback to the user confirming the success or failure of `start` and `stop` operations.
- **FR-009**: The state of the runner (e.g., the process ID of the local server or the name of the docker container) MUST be tracked to enable the `stop` command.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: A developer can successfully start and stop the NATS event bus using both `docker` and `local` runners within 10 seconds of command execution.
- **SC-002**: When the event bus is confirmed as "started" by the CLI, a standard NATS client can successfully connect to `nats://localhost:4222`.
- **SC-003**: When the event bus is confirmed as "stopped" by the CLI, connections to `nats://localhost:4222` are refused.
- **SC-004**: 100% of errors (e.g., Docker not running, `nats-server` not found, port in use) are reported to the user with a descriptive message.

## Assumptions

- The user has Docker installed and the Docker daemon is running when using the `docker` runner.
- The user has the `nats-server` executable in their system's PATH when using the `local` runner.
- The default NATS port `4222` is assumed to be the target port and is expected to be available on the host machine.
- This feature is for local development only and is not intended for production deployments. No persistence, clustering, or security configurations are included.
