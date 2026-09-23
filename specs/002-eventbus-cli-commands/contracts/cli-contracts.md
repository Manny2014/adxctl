# CLI Contracts: `eventbus`

This document defines the command-line interface for the `eventbus` feature.

## 1. `eventbus`

The top-level command for managing the local development event bus.

### Subcommands

- `start`: Starts the event bus.
- `stop`: Stops the event bus.

---

## 2. `eventbus start`

Starts a NATS server using the specified runner.

### Usage

```sh
adxctl eventbus start [flags]
```

### Flags

| Flag       | Type   | Shorthand | Default  | Description                                        |
|------------|--------|-----------|----------|----------------------------------------------------|
| `--runner` | string | `-r`      | `docker` | The runner to use. Options are `docker` or `local`. |

### Output

- **Success**:
  ```
  Starting event bus with docker runner...
  Event bus started successfully.
  ```
  or
  ```
  Starting event bus with local runner...
  Event bus started successfully (PID: 12345).
  ```
- **Failure**:
  ```
  Error: Docker is not running.
  ```
  ```
  Error: 'nats-server' command not found in PATH.
  ```
  ```
  Error: Port 4222 is already in use.
  ```

---

## 3. `eventbus stop`

Stops the NATS server instance that was previously started by `adxctl`.

### Usage

```sh
adxctl eventbus stop
```

### Flags

None.

### Output

- **Success**:
  ```
  Stopping event bus (docker runner)...
  Event bus stopped successfully.
  ```
  or
  ```
  Stopping event bus (local runner, PID: 12345)...
  Event bus stopped successfully.
  ```
- **Failure**:
  ```
  Error: No adxctl-managed event bus is currently running.
  ```
