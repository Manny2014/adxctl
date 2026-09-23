# Quickstart Guide

This guide provides runnable validation scenarios to prove the core functionality of the `adxctl` feature works end-to-end.

## Prerequisites

1.  **Go**: Go compiler installed on your local machine.
2.  **Docker**: Docker daemon running locally.
3.  **NATS Server**: A NATS server running and accessible. The simplest way is via Docker:
    ```bash
    docker run --rm -p 4222:4222 -p 8222:8222 nats:latest
    ```
4.  **Configuration**: `adxctl` is configured with the NATS server address (e.g., via a `config.yaml` or environment variables). `NATS_URL=nats://localhost:4222`.
5.  **GitHub Token**: A GitHub Personal Access Token with `repo` scope is available as an environment variable: `GITHUB_TOKEN=...`

## Scenario 1: Poll GitHub Issues and Process Locally

This scenario validates `FR-002`, `FR-004`, `FR-005`, `FR-006`, and `FR-007`.

### 1. Deploy the Poller

Start a poller to watch for new issues in a specific GitHub repository.

**Command**:
```bash
# Replace with a real repository you have access to
export GITHUB_REPO="your-org/your-test-repo"

adxctl poller deploy --source github --repo $GITHUB_REPO --interval 10s
```

**Expected Outcome**:
-   The command returns a poller ID.
-   Running `adxctl poller list` shows the new poller in `RUNNING` status.
    ```
    ID                                   SOURCE   STATUS    DETAILS
    poller-github-your-org-your-test...  github   RUNNING   Repo: your-org/your-test-repo, Interval: 10s
    ```

### 2. Run the Worker

Start a local worker to listen for tasks from the GitHub poller.

**Command**:
```bash
adxctl worker run --subject "tasks.github" --runtime local
```

**Expected Outcome**:
-   The command returns a worker ID.
-   Running `adxctl worker list` shows the new worker in `RUNNING` status.
    ```
    ID                                   SUBJECT          RUNTIME      STATUS
    worker-local-tasks-github-xyz        tasks.github     local        RUNNING
    ```
-   The worker process starts logging messages indicating it is waiting for tasks.

### 3. Trigger and Verify an Event

Create a new issue in the target GitHub repository.

**Action**:
-   Go to `https://github.com/your-org/your-test-repo/issues` and create a "New Issue".

**Expected Outcome**:
-   Within the 10s polling interval, the poller detects the new issue.
-   The poller publishes a `Task` message to the `tasks.github` subject on NATS.
-   The local worker's console output shows a log message indicating it received and processed the task, similar to:
    ```
    [worker-local-tasks-github-xyz] Received task: ID=..., Source=github, Type=issue_created, Payload={... "title": "My New Test Issue" ...}
    ```

### 4. Stop Services

Clean up the running processes.

**Commands**:
```bash
# Get IDs from the 'list' commands
adxctl poller stop <poller-id>
adxctl worker stop <worker-id>
```

**Expected Outcome**:
-   Running `adxctl poller list` and `adxctl worker list` no longer shows the stopped instances.

## Scenario 2: Run Worker in Docker

This scenario validates the Docker runtime (`FR-003`).

1.  **Prerequisite**: A poller is already running from Scenario 1.
2.  **Command**:
    ```bash
    adxctl worker run --subject "tasks.github" --runtime docker
    ```
3.  **Expected Outcome**:
    -   A worker ID is returned.
    -   `adxctl worker list` shows the worker with `RUNTIME: docker`.
    -   `docker ps` shows a new container running with the worker image.
    -   Creating a new GitHub issue triggers a log message, which can be viewed with `docker logs <container-id>`.
4.  **Cleanup**:
    ```bash
    adxctl worker stop <docker-worker-id>
    ```
    - The `docker ps` command no longer shows the worker container.
