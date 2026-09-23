# Tasks: Eventbus CLI Commands

**Input**: Design documents from `/specs/002-eventbus-cli-commands/`

**Prerequisites**: plan.md (required), spec.md (required for user stories), research.md, data-model.md, contracts/

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2)
- Include exact file paths in descriptions

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Project initialization and directory preparation

- [X] T001 Create `internal/eventbus` package directory structure for eventbus runners and state management

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core infrastructure and state management that MUST be complete before ANY user story can be implemented

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

- [X] T002 Implement state management in `internal/eventbus/state.go` storing `runner_type` ("docker" or "local") and `identifier` (container name `adx-eventbus` or process PID) in `/tmp/adxctl/eventbus.state`
- [X] T003 Create base `eventbus` command in `cmd/eventbus.go` and register it with the root command in `cmd/root.go`

**Checkpoint**: Foundation ready - user story implementation can now begin

---

## Phase 3: User Story 1 - Run eventbus in Docker (Priority: P1) 🎯 MVP

**Goal**: Allow developers to start and stop a local NATS event bus using Docker with `adxctl eventbus start --runner docker` and `adxctl eventbus stop`.

**Independent Test**: Run `adxctl eventbus start --runner docker`, verify container `adx-eventbus` is running via `docker ps` and responding on port 4222, then run `adxctl eventbus stop` and verify container is removed.

### Implementation for User Story 1

- [X] T004 [P] [US1] Implement Docker runner lifecycle (start container `nats:latest` mapping port 4222, stop and remove container) using `moby/moby/client` in `internal/eventbus/docker.go`
- [X] T005 [US1] Implement `start` subcommand in `cmd/eventbus_start.go` with `--runner` flag defaulting to `docker`, integrating Docker runner and persisting state
- [X] T006 [US1] Implement `stop` subcommand in `cmd/eventbus_stop.go` with Docker container shutdown and state cleanup
- [X] T007 [P] [US1] Add unit tests for Docker runner invocation and state tracking in `internal/eventbus/docker_test.go`

**Checkpoint**: At this point, User Story 1 should be fully functional and testable independently (MVP)

---

## Phase 4: User Story 2 - Run eventbus locally (Priority: P2)

**Goal**: Allow developers with `nats-server` installed to run the event bus locally as a background process with `adxctl eventbus start --runner local` and `adxctl eventbus stop`.

**Independent Test**: Run `adxctl eventbus start --runner local`, verify `nats-server` process is running and responding on port 4222, then run `adxctl eventbus stop` and verify the process is terminated.

### Implementation for User Story 2

- [X] T008 [P] [US2] Implement local runner lifecycle (exec `nats-server`, track PID, terminate process gracefully) in `internal/eventbus/local.go`
- [X] T009 [US2] Extend `cmd/eventbus_start.go` to support `--runner local` by invoking the local runner and persisting process PID state
- [X] T010 [US2] Extend `cmd/eventbus_stop.go` to terminate local process by PID when saved `runner_type` is `local`
- [X] T011 [P] [US2] Add unit tests for local runner execution and error handling in `internal/eventbus/local_test.go`

**Checkpoint**: User Stories 1 AND 2 are both functional and independently testable

---

## Phase 5: Polish & Cross-Cutting Concerns

**Purpose**: Validation, error handling refinement, and end-to-end verification

- [X] T012 [P] Refine error messages and diagnostics (Docker daemon unreachable, missing `nats-server` binary, port 4222 conflict) in `cmd/eventbus_start.go` and `cmd/eventbus_stop.go`
- [X] T013 Verify end-to-end scenarios per `specs/002-eventbus-cli-commands/quickstart.md` for both docker and local runners

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies - can start immediately
- **Foundational (Phase 2)**: Depends on Phase 1 completion - BLOCKS all user stories
- **User Story 1 (Phase 3)**: Depends on Phase 2 completion
- **User Story 2 (Phase 4)**: Depends on Phase 2 completion (and builds on start/stop commands from Phase 3)
- **Polish (Phase 5)**: Depends on Phase 3 and Phase 4 completion

### User Story Dependencies

- **User Story 1 (P1)**: Can start after Foundational (Phase 2) - Delivers the core MVP
- **User Story 2 (P2)**: Extends start/stop commands to add the local runner alternative; can be implemented after or in parallel with US1 command wiring

### Within Each User Story

- Runner logic implemented before CLI command binding
- State persistence verified before stop command wiring
- Unit tests written alongside or immediately following implementation

### Parallel Opportunities

- T004 ([US1] Docker runner) and T008 ([US2] Local runner) can be developed in parallel once Phase 2 is complete
- T007 and T011 (unit tests) can run in parallel with polish tasks

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Setup (`internal/eventbus/` directory)
2. Complete Phase 2: Foundational (state management and base command)
3. Complete Phase 3: User Story 1 (Docker runner + `start` and `stop` CLI subcommands)
4. **STOP and VALIDATE**: Test User Story 1 independently using `adxctl eventbus start --runner docker` and `adxctl eventbus stop`

### Incremental Delivery

1. Foundation ready (T001 - T003)
2. Add Docker runner (T004 - T007) -> MVP Complete!
3. Add Local runner (T008 - T011) -> Full feature capabilities
4. Polish & E2E Validation (T012 - T013)
