# Data Model

This document defines the key entities for the Adxctl system based on the feature specification.

## 1. Task

A `Task` represents a single, standardized unit of work that can be processed by a worker. All data from external sources is transformed into this common format.

**Fields**:

| Field      | Type     | Description                                                                 | Required | Example                               |
|------------|----------|-----------------------------------------------------------------------------|----------|---------------------------------------|
| `ID`       | `string` | A unique identifier for the task, typically a UUID generated upon creation. | Yes      | `"a3a7f1a9-3dc4-4b95-a228-7f985ca546b5"` |
| `Source`   | `string` | The origin system of the task.                                              | Yes      | `"github"`, `"jira"`                  |
| `Type`     | `string` | The specific event type from the source system.                             | Yes      | `"issue_created"`, `"comment_added"`    |
| `Payload`  | `JSON`   | The original, unmodified data payload from the source event.                | Yes      | `{"issue": {"title": "Fix bug...", ...}}` |
| `Metadata` | `JSON`   | Additional metadata for routing or processing.                              | No       | `{"priority": "high"}`                |

**State Transitions**:

A `Task` is immutable once created and published to the event bus. Its state is transient and exists only within the message queue until it is consumed and processed by a worker.

1.  `CREATED`: A poller transforms a source event into a `Task` and publishes it.
2.  `ACKNOWLEDGED`: A worker successfully receives and processes the task.
3.  `FAILED` (Implicit): A worker fails to process the task. The handling of this state (e.g., dead-letter queue, retries) is a key part of the worker implementation.

**Validation Rules**:

-   `FR-011`: All pollers MUST transform their source data into this `Task` format.
-   `ID`, `Source`, `Type`, and `Payload` are mandatory fields.

## 2. Poller

A `Poller` is a configurable, long-running process responsible for fetching data from an external system and creating `Task` messages.

**Fields / Configuration**:

| Field      | Type     | Description                                               | Required | Example                                    |
|------------|----------|-----------------------------------------------------------|----------|--------------------------------------------|
| `ID`       | `string` | A unique identifier for the running poller instance.      | Yes      | `"poller-github-my-repo-1663882800"`        |
| `Source`   | `string` | The type of source to poll (e.g., `github`, `jira`).      | Yes      | `"github"`                                 |
| `Interval` | `string` | The frequency at which to poll the source.                | Yes      | `"5m"` (5 minutes)                         |
| `Config`   | `JSON`   | Source-specific configuration (e.g., repository, project key). | Yes      | `{"repo": "my-org/my-repo"}`               |
| `Status`   | `string` | The current runtime status of the poller.                 | Yes      | `"RUNNING"`, `"STOPPED"`, `"ERROR"` |

## 3. Worker

A `Worker` is a process that subscribes to one or more task subjects on the event bus and executes logic based on the received `Task`.

**Fields / Configuration**:

| Field       | Type     | Description                                                  | Required | Example                                      |
|-------------|----------|--------------------------------------------------------------|----------|----------------------------------------------|
| `ID`        | `string` | A unique identifier for the running worker instance.         | Yes      | `"worker-local-tasks-github-1663882900"`     |
| `Subject`   | `string` | The NATS subject(s) to subscribe to.                         | Yes      | `"tasks.github"`                             |
| `Runtime`   | `string` | The environment where the worker is executing.               | Yes      | `"local"`, `"docker"`, `"kubernetes"`        |
| `Status`    | `string` | The current runtime status of the worker.                    | Yes      | `"RUNNING"`, `"STOPPED"`, `"ERROR"`      |

## 4. Runtime

`Runtime` is an enumeration that defines the execution environment for a worker. It is not a data entity in itself but a critical configuration parameter for a `Worker`.

**Possible Values**:

-   `local`: The worker runs as a standard process on the user's local machine.
-   `docker`: The worker runs as a container managed by a local Docker daemon.
-   `kubernetes`: The worker runs as a Pod or Job in a Kubernetes cluster.
