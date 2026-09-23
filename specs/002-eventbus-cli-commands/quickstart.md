# Quickstart: Eventbus CLI Validation

This guide provides the steps to perform an end-to-end validation of the `eventbus` CLI feature.

**Contracts**: For detailed command usage, see [CLI Contracts](./contracts/cli-contracts.md).

## Prerequisites

- `adxctl` binary is built and available.
- For the `docker` runner: Docker is installed and the daemon is running.
- For the `local` runner: The `nats-server` executable is in the system's PATH.
- A NATS client (e.g., `nats-pub`) is available for verification.

## Scenario 1: Docker Runner

This scenario validates the ability to manage the event bus using a Docker container.

### 1. Start the Eventbus

```sh
adxctl eventbus start --runner docker
```

- **Expected Outcome**:
  - The CLI prints a success message indicating the event bus has started.
  - The `docker ps` command shows a container named `adx-eventbus` running the `nats:latest` image with port `4222` exposed.

### 2. Verify Connection

```sh
nats-pub -s nats://localhost:4222 test.subject "hello world"
```

- **Expected Outcome**: The command executes successfully without connection errors.

### 3. Stop the Eventbus

```sh
adxctl eventbus stop
```

- **Expected Outcome**:
  - The CLI prints a success message indicating the event bus has stopped.
  - The `docker ps` command no longer shows the `adx-eventbus` container.
  - Subsequent attempts to connect with `nats-pub` fail.

## Scenario 2: Local Runner

This scenario validates the ability to manage the event bus as a local process.

### 1. Start the Eventbus

```sh
adxctl eventbus start --runner local
```

- **Expected Outcome**:
  - The CLI prints a success message including the PID of the `nats-server` process.
  - A `nats-server` process is running in the background.

### 2. Verify Connection

```sh
nats-pub -s nats://localhost:4222 test.subject "hello local"
```

- **Expected Outcome**: The command executes successfully without connection errors.

### 3. Stop the Eventbus

```sh
adxctl eventbus stop
```

- **Expected Outcome**:
  - The CLI prints a success message.
  - The `nats-server` process started earlier is no longer running.
  - Subsequent attempts to connect with `nats-pub` fail.
