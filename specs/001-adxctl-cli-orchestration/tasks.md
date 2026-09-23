---
description: "Task list for implementing the Adxctl CLI feature."
---

# Tasks: Adxctl - Cloud-Native Agent Task Orchestration CLI

**Input**: Design documents from `/Users/emmanuelrodriguez/git/adxctl/specs/001-adxctl-cli-orchestration/`

**Prerequisites**: plan.md (required), spec.md (required for user stories), research.md, data-model.md, contracts/

**Tests**: Tests are not explicitly requested in the spec, so test tasks are not included in this plan. The focus is on implementing the functional requirements.

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- Include exact file paths in descriptions

## Path Conventions

- The project structure will follow the `plan.md` specification.
- `cmd/` for Cobra commands.
- `internal/` for core logic.
- `pkg/adx` for the reusable library.

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Project initialization and basic structure.

- [X] T001 Create the initial project directory structure as defined in `plan.md`: `cmd/`, `internal/agent`, `internal/pollers`, `internal/events`, `internal/nats`, `pkg/adx`.
- [X] T002 Initialize Go module with `go mod init adxctl` in the project root.
- [X] T003 Add initial Cobra dependency: `go get -u github.com/spf13/cobra@latest`.
- [X] T004 Create the root command file `cmd/root.go` with a basic Cobra command setup.
- [X] T005 [P] Create a placeholder file `main.go` to call the root command.
- [X] T006 [P] Add a basic `.gitignore` file for Go projects.
- [X] T007 [P] Create the configuration file `internal/config/config.go` to manage settings like NATS URL via environment variables.

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core infrastructure that MUST be complete before ANY user story can be implemented.

**⚠️ CRITICAL**: No user story work can begin until this phase is complete.

- [X] T008 [P] Define the common `Task` struct in `internal/events/event.go` as specified in `data-model.md`. Include fields: `ID`, `Source`, `Type`, `Payload`, `Metadata`.
- [X] T009 Implement a NATS client wrapper in `internal/nats/client.go`. This wrapper should handle connection, publishing messages, and graceful shutdown.
- [X] T010 Define the `Poller` and `Worker` data models in `internal/events/event.go` with fields from `data-model.md`.

---

## Phase 3: User Story 1 - Deploy a Poller to Ingest Tasks (Priority: P1) 🎯 MVP

**Goal**: A developer can ingest tasks from GitHub into the event bus using the CLI.

**Independent Test**: Deploy a GitHub poller, create a GitHub issue, and verify a message appears on the `tasks.github` NATS subject.

### Implementation for User Story 1

- [X] T011 [US1] Create the `poller` subcommand structure in `cmd/poller.go`.
- [X] T012 [US1] Implement the `adxctl poller deploy` command in `cmd/poller_deploy.go` with flags `--source`, `--repo`, and `--interval` as per `cli-contracts.md`.
- [X] T013 [P] [US1] Define a `Poller` interface in `internal/pollers/poller.go` with a `Start()` method.
- [X] T014 [US1] Implement the GitHub poller in `internal/pollers/github.go`. It should connect to the GitHub API, fetch issues, and transform them into `Task` messages.
- [X] T015 [US1] The GitHub poller must publish the generated `Task` messages to the `tasks.github` subject using the NATS client.
- [X] T016 [US1] Implement the logic in `cmd/poller_deploy.go` to instantiate and run the correct poller based on the `--source` flag.

**Checkpoint**: At this point, User Story 1 should be functional. A user can deploy a GitHub poller.

---

## Phase 4: User Story 2 - Spin Up a Worker to Process Tasks (Priority: P2)

**Goal**: A developer can process tasks from the event bus using a worker running locally or in Docker.

**Independent Test**: Run a worker for the `tasks.github` subject, and verify it receives and logs tasks when they are published.

### Implementation for User Story 2

- [X] T017 [US2] Create the `worker` subcommand structure in `cmd/worker.go`.
- [X] T018 [US2] Implement the `adxctl worker run` command in `cmd/worker_run.go` with flags `--subject` and `--runtime` as per `cli-contracts.md`.
- [X] T019 [P] [US2] Define a `Runner` interface in `internal/agent/runner.go` with `Start()` and `Stop()` methods.
- [X] T020 [US2] Implement the `local` runner in `internal/agent/runner_local.go`. It should subscribe to the specified NATS subject and log the received `Task` payload.
- [X] T021 [P] [US2] Implement the `docker` runner in `internal/agent/runner_docker.go`. It should start a Docker container with the worker logic inside. A basic Dockerfile will be needed.
- [X] T022 [US2] Implement the logic in `cmd/worker_run.go` to select and use the correct runner based on the `--runtime` flag.

**Checkpoint**: At this point, User Story 2 should be functional. A user can run workers on local and docker runtimes.

---

## Phase 5: User Story 3 - Manage Lifecycle of Workers and Pollers (Priority: P3)

**Goal**: A developer can view the status of, and stop, running workers and pollers.

**Independent Test**: Start a poller and worker, list them, stop them, and verify they are no longer listed.

### Implementation for User Story 3

- [X] T023 [US3] Implement a simple in-memory store or file-based registry to track running processes (pollers and workers) and their metadata (ID, status, runtime).
- [X] T024 [US3] Implement the `adxctl poller list` command in `cmd/poller_list.go` to read from the process registry and display running pollers.
- [X] T025 [US3] Implement the `adxctl poller stop` command in `cmd/poller_stop.go` to terminate the specified poller process gracefully.
- [X] T026 [P] [US3] Implement the `adxctl worker list` command in `cmd/worker_list.go` to display running workers from the process registry.
- [X] T027 [P] [US3] Implement the `adxctl worker stop` command in `cmd/worker_stop.go` to terminate the specified worker process or container.

**Checkpoint**: All user stories should now be independently functional.

---

## Phase 6: Polish & Cross-Cutting Concerns

**Purpose**: Final improvements and fulfilling library requirements.

- [X] T028 Refactor the core logic from `internal/` into the `pkg/adx/` directory to satisfy `FR-010`. Create `pkg/adx/client.go` and `pkg/adx/types.go`.
- [X] T029 Update the `cmd/` implementations to use the new public library from `pkg/adx/`.
- [X] T030 [P] Create a `README.md` file with instructions based on `quickstart.md`.
- [X] T031 [P] Add comments and documentation to all public functions in the `pkg/adx` library.
- [X] T032 Validate the end-to-end scenarios described in `quickstart.md`.

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: Can start immediately.
- **Foundational (Phase 2)**: Depends on Setup. Blocks all user stories.
- **User Stories (Phase 3-5)**: Depend on Foundational. Can be implemented sequentially (US1 -> US2 -> US3).
- **Polish (Phase 6)**: Depends on all user stories being complete.

### User Story Dependencies

- **User Story 1 (P1)**: Depends on Foundational.
- **User Story 2 (P2)**: Depends on Foundational.
- **User Story 3 (P3)**: Depends on Foundational.

### Parallel Opportunities

- Most setup tasks can be done in parallel.
- Within stories, creating command files and interfaces can often be done in parallel with the implementation.
- Different user stories could be implemented in parallel by different developers after the Foundational phase is complete.

---

## Implementation Strategy

### MVP First (User Story 1)

1.  Complete Phase 1: Setup
2.  Complete Phase 2: Foundational
3.  Complete Phase 3: User Story 1
4.  **STOP and VALIDATE**: Test the ability to deploy a poller and see tasks in NATS. This is the first deployable piece of value.

### Incremental Delivery

1.  Deliver MVP (US1).
2.  Add User Story 2 -> Test worker functionality.
3.  Add User Story 3 -> Test lifecycle management.
4.  Complete Polish phase to finalize the library for external consumption.

---
## Phase 7: Convergence
- [X] T033 CRITICAL Refactor core logic into a reusable library in `pkg/adx/` per `FR-010` (`missing`)
- [X] T034 CRITICAL Unify the `Task` struct definition in `internal/events/event.go` to match `data-model.md` and update all usages per `FR-011` (`contradicts`)
- [X] T035 Implement the GitHub poller in `internal/pollers/github.go` per `US1` (`missing`)
- [X] T036 Implement `adxctl poller deploy` command with specified flags per `US1/AC1` (`contradicts`)
- [X] T037 Implement `adxctl worker run` command with specified flags per `US2/AC1` (`contradicts`)
- [X] T038 Implement `list` and `stop` commands for pollers and workers per `US3` (`missing`)
- [X] T039 Implement the Docker runner logic in `internal/agent/runner_docker.go` per `US2` (`partial`)
- [X] T040 Review and justify or remove unrequested files and abstractions per `plan.md` (`unrequested`)

---
## Phase 8: Convergence
- [X] T041 Make the NATS subject in the GitHub poller dynamic per `US1/AC2` (`partial`)
- [X] T042 Improve the Docker runner implementation to be more configurable per `US2/AC3` (`partial`)
- [X] T043 Implement process termination logic in the `stop` commands per `US3/AC2, US3/AC3` (`partial`)

---
## Phase 9: Convergence
- [ ] T044 Implement the Kubernetes runner in `internal/agent/runner_kubernetes.go` per `FR-003, FR-009, US2` (`partial`)
- [ ] T045 Implement the Bitbucket poller per `FR-004` (`partial`)

