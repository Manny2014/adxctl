# adxctl

`adxctl` is a highly scalable Command Line Interface (CLI) harness and worker agent system designed to bridge localized CLI execution with highly distributed, AI-driven task orchestration. 

It provides an architectural foundation for pulling events from external tools (such as Jira), dispatching them to an intelligent, decoupled event bus, assessing task complexity using advanced probabilistic AI inference (Google Gemini), and executing workloads across varied runner topologies (Local scripts/CLIs, containerized Docker, or distributed Kubernetes clusters).

---

## Architecture Overview

`adxctl` is built on a clean, decoupled, event-driven architecture designed to scale seamlessly from a single developer machine to a highly distributed, production-grade cloud environment.

```
+------------------+         +--------------------------+         +-------------------------+
|                  |  Poll   |                          | Publish |                         |
|  External Tools  |-------->|   adxctl poller start    |-------->|   Event Bus (NATS)      |
|  (e.g., Jira)    |         |                          |         |   (Stream: JIRA_EVENTS) |
+------------------+         +--------------------------+         +-------------------------+
                                                                               |
                                                                               | Subscribe
                                                                               v
+------------------+         +--------------------------+         +-------------------------+
|  AI Classifier   |  Infer  |                          |         |                         |
|  (Google Gemini) |<--------|    adxctl agent start    |<--------|   Agent Orchestrator    |
|  (Complexity)    |         |  (Agent Orchestrator)    |         |                         |
+------------------+         +--------------------------+         +-------------------------+
                                           |
                                           | Dispatch based on Complexity
                                           v
                             +---------------------------+
                             |       Task Runners        |
                             +---------------------------+
                             | - Local (CLI / script)    |
                             | - Docker (Containerized)  |
                             | - Kubernetes (Distributed)|
                             +---------------------------+
```

### Core Components

#### 1. Generic Event Bus Abstraction (`eventbus`)
The central communication spine of `adxctl` is the `eventbus` abstraction, enabling fully decoupled publishers and consumers. 
* **Implementation:** Currently backed by **NATS JetStream** for persistent, high-performance messaging.
* **Stream Management:** Employs explicit streams (such as `JIRA_EVENTS`) configured with safe, non-overlapping subjects (e.g., `jira.>`) to completely prevent conflicts with internal NATS system namespaces while keeping publish acknowledgments enabled.

#### 2. Stateless Data Pollers (`poller`)
Pollers are scheduled systems that synchronize external platforms with our state.
* **Jira Poller:** Syncs with Jira's search JQL endpoint, fetching `*all` fields.
* **Dynamic Event Detection:** It queries state dynamically from a Key-Value watermark store (`jira_poller_state`) and compares the issue's `created` timestamp against the watermark:
  - If `createdTime` > `lastPollTime`, it publishes a `jira.issue_created.<Project>` event.
  - Otherwise, it publishes a `jira.issue_updated.<Project>` event.

#### 3. Agent Orchestrator (`agent`)
The worker daemon that subscribes to the event bus stream. Upon consuming a task payload:
* It invokes the **AI Complexity Classifier** to run probabilistic inference.
* Based on the task payload details, the AI model classifies the workload as `simple` or `complex`.
* It routes the task to the corresponding execution environment.

#### 4. Probabilistic AI Router (`complexity`)
Integrates directly with **Google Cloud AI (Gemini)**.
* Evaluates raw task details (e.g. code modifications, payload sizes, fields changed) using intelligent prompts.
* Automatically selects a `local` or `remote` runner depending on the returned complexity assessment.

#### 5. Multi-Environment Runners (`runner`)
Executes tasks across diverse infrastructures:
* **Local Runner:** Ideal for quick, lightweight scripts. Supports three execution types:
  - `gemini-cli`: Leverages local AI capabilities.
  - `antigravity-cli`: Integrates with migrations and external CLI systems.
  - `local-script`: Runs any local executable script, piping the raw task JSON payload directly to the script's `stdin`.
* **Docker Runner:** Isolates and executes medium-heavy tasks inside secure container environments.
* **Kubernetes Runner:** Dispatches intensive, long-running, or highly parallel computations to an elastic Kubernetes cluster.

---

## Quick Start

Follow these steps to spin up the entire `adxctl` platform locally to visualize your agentic workflows executing in real-time.

### Prerequisites
- **Go** (v1.27+)
- **Node.js** (v18+) and `npm`
- **Docker** (for running NATS eventbus)
- A **Google Gemini API key** (for AI complexity routing)

### 1. Clone & Build the Platform
Compile the CLI binary and package the Single Page Application:
```bash
# Compile the adxctl binary
make build

# Install UI dependencies and build the static assets
npm install --prefix web
npm run build --prefix web
```

### 2. Set Up Environment Variables
Export your credentials. These will be dynamically expanded within `config.yaml` upon loading:
```bash
export GOOGLE_API_KEY="your-gemini-api-key"
export JIRA_API_TOKEN="your-jira-token"
```

### 3. Spin Up the Event Bus
Start NATS JetStream in the background inside a Docker container:
```bash
make eventbus-start
```

### 4. Run the API and UI Servers
Launch the API service (responsible for tracking queue and agent status) and the UI dashboard:
```bash
# Start the API server (default port 8080)
./bin/adxctl api --port 8080 &

# Start the static UI file server (default port 3000)
./bin/adxctl ui --port 3000 --dir ./web/build &
```
*Now, open your browser and navigate to **`http://localhost:3000`**.*

### 5. Launch the Agent Orchestrator
Start the AI-powered task consumer. This agent subscribes to the event bus, requests complexity classification from Gemini, and routes tasks to the appropriate local, docker, or Kubernetes runners:
```bash
./bin/adxctl agent start
```

### 6. Trigger Workloads (Start Polling)
Activate your Jira data poller to start ingestion and watch the workflows populate on your live UI dashboard:
```bash
./bin/adxctl poller start jira-default
```

---

## Configuration

`adxctl` is fully configuration-driven and supports environment variable expansion natively within the YAML config (allowing you to protect secret values like tokens).

See `config.example.yaml` for a template:

```yaml
# config.yaml
eventbus:
  type: nats
  nats:
    runner: "local-cli" # local-cli or docker
    port: 4222
    store_dir: ""

pollers:
  jira-prod:
    type: jira
    enabled: true
    interval: 1m
    domain: "your-domain.atlassian.net"
    email: "user@example.com"
    api_token: "${JIRA_API_TOKEN}" # Dynamically expanded
    nats_url: "nats://127.0.0.1:4222"

agent:
  runners:
    default: local
    local:
      type: local-script
      script: "/usr/local/bin/my-task-processor.sh"
    docker: {}
    kubernetes: {}
  ai_provider:
    type: google
    google:
      api_key: "${GOOGLE_API_KEY}" # Dynamically expanded
      model: "gemini-pro"
```

---

## Command Line Interface (CLI)

Build and run `adxctl` using the optimized targets provided in the `Makefile`.

### 1. Build and Setup
```bash
# Build the production-ready optimized binary (saved to bin/adxctl)
make build

# Display all available targets
make help
```

### 2. Event Bus Management
```bash
# Start the event bus container (Docker)
make eventbus-start

# Inspect event bus stream statistics
make nats-stream-info

# View live message streams (optionally filter with SUBJECT=...)
make nats-stream-view SUBJECT="jira.issue_created.*"
```

### 3. Running Data Pollers
```bash
# Start a defined poller by its name from config.yaml
./bin/adxctl poller start jira-prod
```

### 4. Running Agent Workers
```bash
# Launch the agent worker orchestrator to consume events and execute tasks
./bin/adxctl agent start
```

### 5. API & UI
```bash
# Run the API server
./bin/adxctl api --port 8080

# Build and run the UI server
npm run build --prefix web
./bin/adxctl ui --port 3000 --dir ./web/build
```
