# Phase 0: Research Plan

This document outlines the research tasks required to resolve the "NEEDS CLARIFICATION" items identified in the implementation plan.

## Research Tasks

### 1. Project Setup & Versioning

-   **Task**: Determine the optimal and compatible versions for the core dependencies.
-   **Questions**:
    -   What is the current stable and recommended version of Golang for new projects?
    -   What are the latest stable versions for `cobra`, `nats.go`, `docker/client`, and `kubernetes/client-go`?
    -   Are there any known compatibility issues between these libraries?

### 2. Dependency Best Practices

-   **Task**: Research best practices for using the primary dependencies to ensure a robust and maintainable implementation.
-   **Questions**:
    -   **Cobra**: What is the recommended structure for a complex CLI with nested subcommands for managing different resources (pollers, workers)?
    -   **NATS (`nats.go`)**: What are the best practices for connection pooling, handling reconnects, and managing subscriptions gracefully in a long-running agent? How should we structure our subjects for discoverability and scalability?
    -   **Docker API (`docker/client`)**: What is the standard pattern for managing the lifecycle of a container (create, run, inspect, stop, remove) from a Go application? How should we handle authentication with the Docker daemon?
    -   **Kubernetes API (`kubernetes/client-go`)**: What is the idiomatic way to create and manage Kubernetes Jobs or Pods for the worker runtime? How should the application handle authentication (in-cluster vs. out-of-cluster)?

### 3. Testing Strategy

-   **Task**: Define a comprehensive testing strategy for the project.
-   **Questions**:
    -   What are the community-standard mocking frameworks in the Go ecosystem (e.g., `testify/mock`, `gomock`)?
    -   What are effective patterns for writing integration tests that involve NATS, Docker, and Kubernetes? (e.g., using test containers, embedded NATS server).
    -   How can we structure tests to align with the "Independently Testable Features" principle?

### 4. Scalability and Performance

-   **Task**: Clarify the expected scale to inform the design of the orchestration logic.
-   **Questions**:
    -   What is a reasonable starting assumption for the number of concurrent pollers and workers the system should support (e.g., 10s, 100s, 1000s)?
    -   What are potential bottlenecks for task ingestion and processing at scale? (e.g., NATS throughput, Docker daemon limits, K8s API rate limits).
    -   What concurrency patterns in Go (e.g., worker pools, fan-out/fan-in) would be appropriate for the agent's task processing loop?
