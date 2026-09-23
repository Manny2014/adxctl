# Feature Specification: Adxctl - Cloud-Native Agent Task Orchestration CLI

**Feature Branch**: `[001-adxctl-cli-orchestration]`

**Created**: 2026-09-22

**Status**: Draft

**Input**: User description: "Build a cloud native event-driven agent task orchestration cli called "Adxctl"

This tool is intended to be used as a local "code-factory" that allows developers to spin up multiple "workers" on different runtimes such as kubernetes,docker, or locally via gemini-cli or antigravity.

The works pull work from event-bus which we want to use Nats because we want to keep the solution cloud native.

Additionally, the cli should support deploying "pollers" that will pull data/events from sources like jira, github, bitbucket server, ect.. and put those "tasks" into the evenbus. Each subject should start with type based on the source to allow workers to pull from the right subject.

Use golang and leverage the cobra cli tool to create a rich CLI."

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Deploy a Poller to Ingest Tasks (Priority: P1)

A developer wants to ingest tasks from an external source (e.g., GitHub issues) into the event bus. They use the `adxctl` CLI to configure and deploy a poller for that source.

**Why this priority**: This is the entry point for all work in the system. Without pollers, no tasks can be created for workers to consume.

**Independent Test**: Can be tested by deploying a single poller for a mock source. Success is verified if the poller runs, queries the source, and publishes messages to the correct NATS subject.

**Acceptance Scenarios**:

1.  **Given** a running NATS instance, **When** a developer runs `adxctl deploy poller --source github --repo my-org/my-repo`, **Then** a new poller instance is created and starts running.
2.  **Given** a running GitHub poller, **When** a new issue is created in the configured repository, **Then** a new message is published to the `tasks.github` subject in NATS.

---

### User Story 2 - Spin Up a Worker to Process Tasks (Priority: P2)

A developer needs to process tasks from a specific source. They use `adxctl` to spin up a worker that subscribes to the relevant NATS subject and executes the task logic.

**Why this priority**: This is the core "work-performing" functionality. It enables the distributed execution of tasks ingested by pollers.

**Independent Test**: Can be tested by running a single worker that subscribes to a NATS subject. Success is verified if the worker receives a message published to that subject and executes its predefined task logic.

**Acceptance Scenarios**:

1.  **Given** a running NATS instance with tasks on the `tasks.github` subject, **When** a developer runs `adxctl run worker --runtime local --subject tasks.github`, **Then** a new worker process starts locally.
2.  **Given** a running local worker, **When** a message appears on the `tasks.github` subject, **Then** the worker consumes the message and executes its task logic (e.g., logging the issue title).
3.  **Given** the need to run on a different runtime, **When** a developer runs `adxctl run worker --runtime docker --subject tasks.github`, **Then** a new Docker container is started with the worker process inside.

---

### User Story 3 - Manage Lifecycle of Workers and Pollers (Priority: P3)

A developer needs to view the status of, and stop, running workers and pollers.

**Why this priority**: Provides essential control and visibility over the system, preventing resource leaks and allowing for system management.

**Independent Test**: Can be tested by starting, listing, and stopping a single poller or worker.

**Acceptance Scenarios**:

1.  **Given** a running poller and a running worker, **When** a developer runs `adxctl list`, **Then** the output shows both the poller and the worker with their status (e.g., running) and runtime.
2.  **Given** a running poller, **When** a developer runs `adxctl stop poller <poller-id>`, **Then** the poller process is terminated gracefully.
3.  **Given** a running worker, **When** a developer runs `adxctl stop worker <worker-id>`, **Then** the worker process is terminated gracefully.

---

### Edge Cases

-   How does the system handle a NATS connection failure for a running poller or worker?
-   What happens if a worker fails to process a task? Is there a retry mechanism?
-   How are credentials for external sources (GitHub, Jira) managed securely? Credentials will be provided via Environment Variables.

## Requirements *(mandatory)*

### Functional Requirements

-   **FR-001**: System MUST provide a CLI named `adxctl`.
-   **FR-002**: `adxctl` MUST allow users to deploy and manage "pollers".
-   **FR-003**: `adxctl` MUST allow users to run and manage "workers" on different runtimes (local, Docker, Kubernetes).
-   **FR-004**: Pollers MUST fetch data from external sources (Jira, GitHub, Bitbucket Server) and publish them as tasks to a NATS event bus.
-   **FR-005**: Task messages published to NATS MUST use a subject that indicates the source (e.g., `tasks.github`, `tasks.jira`).
-   **FR-006**: Workers MUST subscribe to NATS subjects to consume and process tasks.
-   **FR-007**: The CLI MUST provide commands to list and stop running pollers and workers.
-   **FR-008**: The system MUST be implemented in Golang using the Cobra CLI framework.
-   **FR-009**: The solution MUST be "cloud native", implying container-based deployment and orchestration.
-   **FR-010**: The tool's core logic MUST be exposed as a Go library that `gemini-cli`/`antigravity` can import and use directly.
-   **FR-011**: All pollers MUST transform their source data into a single, common `Task` format.


### Key Entities *(include if feature involves data)*

-   **Task**: A unit of work to be performed. It has a `source` (e.g., GitHub), a `type` (e.g., `issue_created`), and a `payload` (the data from the source event).
-   **Poller**: A long-running process that connects to an external data source, fetches events, and transforms them into `Task` messages.
-   **Worker**: A process that subscribes to `Task` messages from the event bus and executes business logic based on the task's content.
-   **Runtime**: The environment where a worker executes (e.g., local process, Docker container, Kubernetes pod).

## Success Criteria *(mandatory)*

### Measurable Outcomes

-   **SC-001**: A developer can successfully deploy a poller and see tasks appear in NATS within 5 minutes of the initial command.
-   **SC-002**: A developer can spin up a worker for a given task subject and have it start processing messages within 1 minute.
-   **SC-003**: The CLI response time for any `adxctl` command MUST be less than 2 seconds for feedback to the user.
-   **SC-004**: The system can process a throughput of at least 100 tasks per minute per worker instance for a simple processing task (e.g., logging).

## Assumptions

-   Users have a running NATS server and have the necessary connection details.
-   Users have the required permissions to run processes locally, as Docker containers, or to deploy to a Kubernetes cluster.
-   The initial version will focus on the core orchestration logic, and the actual "business logic" within workers will be simple placeholders (e.g., logging the task).
-   The term "code-factory" implies that this tool facilitates the generation and execution of automated development-related tasks.
-   "Cloud native" implies the architecture should be built on principles of containerization, microservices, and dynamic orchestration.
-   Credentials for external sources will be passed via environment variables.
-   A standardized task schema will be used across all sources.
-   The core logic will be exposed as a Go library for other tools to consume.
