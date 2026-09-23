# Implementation Plan: Adxctl - Cloud-Native Agent Task Orchestration CLI

**Branch**: `001-adxctl-cli-orchestration` | **Date**: 2026-09-22 | **Spec**: [./spec.md](./spec.md)

**Input**: Feature specification from `/Users/emmanuelrodriguez/git/adxctl/specs/001-adxctl-cli-orchestration/spec.md`

**Note**: This template is filled in by the `/speckit.plan` command; its definition describes the execution workflow.

## Summary

The feature is to build a cloud-native, event-driven agent task orchestration CLI named "adxctl". It will be written in Golang using the Cobra framework. The tool will manage "pollers" to ingest tasks from sources like Jira and GitHub into a NATS event bus, and "workers" to process these tasks on local, Docker, or Kubernetes runtimes. The core logic will be exposed as a Go library.

## Technical Context

**Language/Version**: Golang (version NEEDS CLARIFICATION)

**Primary Dependencies**: Cobra, NATS Go Client, Docker API Client, Kubernetes API Client (all NEED CLARIFICATION on versions and best practices)

**Storage**: NATS (for task queueing), no long-term persistence defined.

**Testing**: Go's standard testing library. (NEEDS CLARIFICATION: Are mocking frameworks like testify/mockery or specific integration testing patterns preferred?)

**Target Platform**: Linux/macOS/Windows for the CLI, with worker runtimes on local host, Docker, and Kubernetes.

**Project Type**: CLI with a reusable core library.

**Performance Goals**: Process at least 100 tasks/minute/worker. CLI response < 2s.

**Constraints**: Must be implemented as a Go library consumable by other tools (`gemini-cli`/`antigravity`). Secure credential management via environment variables.

**Scale/Scope**: Support for multiple, concurrent pollers and workers across different runtimes. (NEEDS CLARIFICATION: Expected number of concurrent tasks, pollers, and workers for initial design.)

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

- **I. Independently Testable Features**: **PASS**. The feature spec breaks down functionality into User Stories (Pollers, Workers, Management) that are independently testable. The requirement for a core Go library (FR-010) further supports this principle by enforcing a clear separation of concerns.
- **II. Explicit & Tested Business Logic**: **PASS**. The core logic (polling, task transformation, worker execution) will be placed in dedicated internal modules. The `Task` entity is explicitly defined (FR-011). The testing strategy requires clarification but the structure allows for it.
- **III. Simplicity Over Abstraction**: **PASS**. The initial approach uses well-defined, existing technologies (Go, Cobra, NATS, Docker/K8s APIs) without introducing speculative layers. The design maps directly to the requirements.

**Evaluation**: No violations detected.

## Project Structure

### Documentation (this feature)

```text
specs/001-adxctl-cli-orchestration/
├── plan.md              # This file (/speckit.plan command output)
├── research.md          # Phase 0 output (/speckit.plan command)
├── data-model.md        # Phase 1 output (/speckit.plan command)
├── quickstart.md        # Phase 1 output (/speckit.plan command)
├── contracts/           # Phase 1 output (/speckit.plan command)
└── tasks.md             # Phase 2 output (/speckit.tasks command - NOT created by /speckit.plan)
```

### Source Code (repository root)

```text
cmd/
├── agent_start.go
├── poller_start.go
└── root.go
internal/
├── agent/
│   ├── runner.go
│   ├── runner_local.go
│   ├── runner_docker.go
│   └── runner_kubernetes.go
├── pollers/
│   ├── jira.go
│   └── github.go
├── events/
│   └── event.go
└── nats/
    └── client.go
pkg/
└── adx/
    ├── client.go
    └── types.go
```

**Structure Decision**: The project follows a standard Go CLI application structure. 
- `cmd/` will contain the Cobra command definitions.
- `internal/` will house the core application logic for agents (workers), pollers, and NATS communication. Specific runner implementations will be in `internal/agent/`.
- `pkg/adx/` will be created to satisfy `FR-010`, exposing the core functionality as a reusable library for external tools.

## Complexity Tracking

> **Fill ONLY if Constitution Check has violations that must be justified**

| Violation | Why Needed | Simpler Alternative Rejected Because |
|-----------|------------|-------------------------------------|
| N/A       | -          | -                                   |
