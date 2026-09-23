# CLI Contracts

This document specifies the command-line interface for `adxctl`.

## Top-Level Commands

```
adxctl <command>
```

## 1. Poller Management

### `adxctl poller deploy`

Deploys and starts a new poller.

**Usage**:
```
adxctl poller deploy --source <type> [flags]
```

**Arguments & Flags**:

-   `--source` (string, required): The type of poller to deploy (e.g., `github`, `jira`).
-   `--interval` (string, optional, default="1m"): The polling interval (e.g., "30s", "5m", "1h").
-   **GitHub Source Flags**:
    -   `--repo` (string, required): The GitHub repository in `owner/name` format.
-   **Jira Source Flags**:
    -   `--project` (string, required): The Jira project key.
    -   `--jql` (string, optional): A custom JQL query to filter issues.

**Example**:
```bash
# Deploy a poller for a GitHub repository
adxctl poller deploy --source github --repo my-org/my-repo --interval 5m
```

### `adxctl poller list`

Lists all running pollers.

**Usage**:
```
adxctl poller list
```

**Output Format**:
```
ID                                   SOURCE   STATUS    DETAILS
poller-github-my-org-my-repo-xyz     github   RUNNING   Repo: my-org/my-repo, Interval: 5m
```

### `adxctl poller stop`

Stops a running poller.

**Usage**:
```
adxctl poller stop <poller-id>
```

**Arguments**:
-   `<poller-id>` (string, required): The ID of the poller to stop.

## 2. Worker Management

### `adxctl worker run`

Runs a new worker on a specified runtime.

**Usage**:
```
adxctl worker run --subject <subject> --runtime <runtime> [flags]
```

**Arguments & Flags**:
-   `--subject` (string, required): The NATS subject for tasks (e.g., `tasks.github`, `tasks.*`).
-   `--runtime` (string, required): The execution environment (`local`, `docker`, `kubernetes`).
-   `--name` (string, optional): A custom name for the worker instance.

**Example**:
```bash
# Run a worker locally to process GitHub tasks
adxctl worker run --subject tasks.github --runtime local
```

### `adxctl worker list`

Lists all running workers.

**Usage**:
```
adxctl worker list
```

**Output Format**:
```
ID                                   SUBJECT          RUNTIME      STATUS
worker-local-tasks-github-abc        tasks.github     local        RUNNING
worker-docker-tasks-jira-def         tasks.jira       docker       RUNNING
```

### `adxctl worker stop`

Stops a running worker.

**Usage**:
```
adxctl worker stop <worker-id>
```

**Arguments**:
-   `<worker-id>` (string, required): The ID of the worker to stop.
