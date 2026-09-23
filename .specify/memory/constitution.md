<!--
    Sync Impact Report
    - Version: 1.2.0 -> 1.3.0
    - Modified Principles: None
    - Added sections: Core Principles/VI. Dependency Immutability
    - Removed sections: None
    - Follow-up TODOs: None
-->
# Example-1 Constitution

## Core Principles

### I. Independently Testable Features
Features MUST be designed and implemented to be independently testable in isolation from unrelated components. Each feature must define clear interfaces, dependencies, and test boundaries so that unit and integration verification can run autonomously and reliably.

*Rationale*: Independent testability prevents regressions, reduces debugging cognitive load, accelerates development feedback loops, and simplifies automated continuous integration.

### II. Explicit & Tested Business Logic
Business logic MUST be represented explicitly in dedicated modules and domain models rather than hidden within UI components, side-effect hooks, or framework glue. All business logic rules, domain calculations, state transitions, and edge cases MUST be covered by automated tests.

*Rationale*: Clear representation and rigorous testing of core business rules ensure correctness, prevent unintended behavioral side effects, and serve as verifiable documentation of system requirements.

### III. Simplicity Over Abstraction
Architectural and design decisions MUST prefer simple, readable, and maintainable solutions over speculative abstractions, premature generalization, and unnecessary indirection (YAGNI). Abstractions are permitted only when justified by concrete, immediate requirements.

*Rationale*: Premature abstractions increase maintenance overhead, obscure intent, and slow down future iteration without providing immediate practical value.

### IV. Containerized CI Testing
All features MUST include automated tests that can be executed within a containerized Continuous Integration (CI) environment. Test suites SHOULD be self-contained, requiring no external dependencies beyond what is defined in the container, to ensure portability and reliable execution.

*Rationale*: Containerized testing guarantees a consistent and reproducible test environment, eliminates "works on my machine" issues, and enables reliable, automated quality gates in the CI/CD pipeline.

### V. Continuous Compilation
After any update to dependencies or significant code changes, local compilation commands MUST be executed to ensure the project builds successfully. This practice MUST be followed to prevent integration issues and maintain a stable codebase.

*Rationale*: Frequent compilation catches errors early, reduces the risk of breaking the build, and ensures that all changes are integrated correctly, leading to a more stable and reliable development process.

### VI. Dependency Immutability
Code that is a dependency or is not directly owned by the project MUST NEVER be modified. All external libraries, packages, vendored modules, and upstream dependencies are strictly immutable. Any necessary behavioral adjustments, bug workarounds, or custom integrations MUST be achieved through external adapters, composition, configuration, or upstream contributions rather than in-place changes to unowned code.

*Rationale*: Directly altering dependency or unowned code compromises build reproducibility, creates untracked drift, breaks upgrade and patching paths, and introduces hidden maintenance liabilities.

## Quality & Architectural Standards
- **Test Boundaries**: Tests must target isolated units and explicit interfaces, avoiding fragile coupling to external or shared mutable state.
- **Complexity Review**: Any proposed abstraction or architectural layer must be justified by demonstrable necessity and measurable maintainability benefits.
- **Maintainability First**: Code readability, deterministic execution, and ease of onboarding take precedence over clever or overly condensed constructs.
- **Dependency Integrity**: External code and dependencies must remain untouched; customizations must use supported extension points, wrapper abstractions, or official upstream mechanisms.

## Development Workflow & Verification Gates
- **Pre-Merge Verification**: Automated test suites validating business logic and feature isolation must pass before changes are accepted.
- **Review Adherence**: Pull requests and design reviews must verify compliance with core principles prior to merging.

## Governance
This constitution defines the fundamental engineering standards and practices for Example-1.

- **Supremacy**: This constitution supersedes informal practices and conflicting conventions.
- **Amendments**: Amendments require documented rationale, an evaluation of project impact, and team approval.
- **Versioning Policy**: Constitution versioning adheres strictly to semantic versioning:
  - **MAJOR**: Structural changes, removal of rules, or fundamental redefinition of core principles.
  - **MINOR**: Addition of new principles, new sections, or materially expanded governance rules.
  - **PATCH**: Non-semantic clarifications, wording refinements, and typo fixes.
- **Compliance**: Ongoing development, pull requests, and architecture reviews must continuously verify compliance with these principles.

**Version**: 1.3.0 | **Ratified**: 2026-09-22 | **Last Amended**: 2026-09-22
