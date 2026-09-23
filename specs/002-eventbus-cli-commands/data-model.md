# Data Model: Eventbus

This feature is primarily stateless from a data persistence perspective. The only "data" it manages is the transient state required to link the `start` and `stop` commands.

## 1. Runner State

- **Description**: Represents the state of the running event bus instance. This state is stored temporarily on the local filesystem to be accessible across different `adxctl` command invocations.
- **Location**: A file within the project's temporary directory (e.g., `/tmp/adxctl/`).
- **Attributes**:
  - `runner_type` (string): The type of runner used, either "docker" or "local".
  - `identifier` (string): The identifier for the running instance.
    - For the `docker` runner, this is the container name (`adx-eventbus`).
    - For the `local` runner, this is the process ID (PID) of the `nats-server` process.
- **Format**: A simple JSON or plain text file is sufficient. For example, a `eventbus.state` file could contain:
  ```json
  {
    "runner_type": "local",
    "pid": 12345
  }
  ```
  or
  ```json
  {
    "runner_type": "docker",
    "container_name": "adx-eventbus"
  }
  ```
- **State Transitions**:
  - **`adxctl eventbus start`**: Creates or overwrites the state file with the details of the newly started instance.
  - **`adxctl eventbus stop`**: Reads the state file to identify the instance to terminate, then deletes the state file upon successful termination.
