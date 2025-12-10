## 🚀 DIMAG IngestList API Wrapper – Documentation and Quickstart

This README describes the available endpoints of the DIMAG IngestList API Wrapper and shows how to quickly try them out.

> Tip: The examples are designed so you can copy‑paste and run them directly in your terminal.

---

## ✨ Overview / Highlights

A quick, skimmable summary of what this project provides.

- 🔌 API for digital preservation jobs:
  - Validate or Identify files via a single endpoint (`POST /api/create`).
  - Works with existing server files (remote) and direct uploads (multipart).
  - Inspect processing results via Jobs (`GET /api/jobs`, `GET /api/job/{id}`).
- 🔐 Authentication with JWT: simple login (`POST /api/login`) returns a token for subsequent calls.
- 🧩 Two operation types: Validate and Identify, unified request model.
- 🧭 Developer‑friendly docs: copy‑pasteable curl snippets and end‑to‑end workflows.
- ⚙️ Configuration via YAML with sensible defaults and auto‑creation of `config.yaml` if missing.
- 🗃️ Database via GORM: SQLite for local/dev; PostgreSQL for production/Kubernetes.
- 🧵 Background processing:
  - robfig/cron‑based scheduler (configurable interval) to process queued tasks.
  - Worker pool with configurable parallelism and automatic cleanup of finished work.
- 🗄️ File storage management: local storage directory for uploads and processing artifacts.
- 📈 Prometheus metrics endpoint `/metrics` with custom business and storage metrics, plus HTTP counters/histograms.
- 📦 Containerized application: multi‑stage Docker/Podman build and ready‑to‑run images.
- ☸️ Kubernetes‑ready: manifests for PostgreSQL and guidance for deploying the app (Deployment/Service/Ingress).
- 📚 Clear project structure and major dependencies (Gin, GORM, robfig/cron, JWT, Prometheus, YAML, UUID).
- 💬 Feedback & contributions welcome: open issues/PRs; MIT‑licensed for hassle‑free adoption.

---

### 📚 Table of Contents

- [Project Structure](#project structure)
- [Dependencies](#dependencies)
- [Notes](#notes)
- [Configuration](#configuration)
- [Metrics (Prometheus)](#metrics-prometheus)
- [API](#api)
  - [Quickstart (TL;DR)](#quickstart-tldr)
  - [Base URL](#base-url)
  - [Prerequisites](#prerequisites)
  - [Authentication](#authentication)
  - [Endpoints](#endpoints)
    - [Login](#login)
    - [Remote-Validate](#remote-validate)
    - [Remote-Identify](#remote-identify)
    - [Local-Validate](#local-validate)
    - [Local-Identify](#local-identify)
    - [Jobs](#jobs)
    - [Job by ID](#job-by-id)
  - [Operation Types](#operation-types)
  - [Workflows](#workflows)
  - [Differences between Remote and Local](#differences-between-remote-and-local)
- [Container Image (Docker/Podman)](#container-image-dockerpodman)
- [Kubernetes](#kubernetes)
- [Feedback & Contributions](#feedback--contributions)
- [License](#license)

---

## 🗂️ Project Structure

High-level overview of the repository layout:

```
.
├── backend                  # Go backend service
│   ├── cmd/server           # Main entry point
│   ├── internal             # Application layers (services, infra, worker, scheduler)
│   ├── pkg/config           # YAML config loader and defaults
│   └── Dockerfile           # Multi-stage image build
├── conf                     # Example or local configuration (mount into container)
├── frontend                 # (Reserved for UI, if applicable)
├── k8s                      # Kubernetes manifests (Postgres)
├── scripts                  # Helper scripts
└── README.md                # This documentation
```

## 📦 Dependencies

Major libraries used by the backend (see `backend/go.mod` for exact versions):

- gin-gonic/gin — HTTP web framework (routing, middleware)
- gin-contrib/cors — CORS middleware for Gin
- gorm.io/gorm — ORM used for persistence
  - gorm.io/driver/sqlite — SQLite driver (default dev mode)
  - gorm.io/driver/postgres — PostgreSQL driver (for production/K8s option)
- robfig/cron/v3 — Cron-like scheduler used for background workers (`task_scheduler.interval`)
- golang-jwt/jwt/v5 — JWT creation and validation for auth
- prometheus/client_golang — Prometheus metrics instrumentation
- gopkg.in/yaml.v3 — YAML configuration parsing
- google/uuid — UUID generation for entities/jobs

Notes:
- Go toolchain is defined via the container image (`golang:1.25-alpine`) and `backend/go.mod` (`go 1.24.x`).
- To update or add a dependency locally: `cd backend && go get <module>@latest && go mod tidy`.

## 📝 Notes

- ⏱️ JWT tokens have an expiration (`exp` claim). After expiry, log in again to obtain a new token.
- 🔁 Both operation types (Validate and Identify) use the same endpoint `/api/create`; the operation type is controlled
  via the `type` parameter.
- 🧪 In all examples, replace variables like `$BASE_URL` and `$TOKEN` with your values if you don't use environment variables.

### 🧯 Troubleshooting

- Login returns `401 Unauthorized`:
  - Check credentials (`email`, `password`).
  - Ensure `Content-Type: application/json` is set.
- Upload fails (`400 Bad Request`):
  - Make sure you use `multipart/form-data` and the field `file=@/path/to/file` exists.
  - Verify that `type` is exactly `Validate` or `Identify` (case‑sensitive!).
- Remote operations return `404`:
  - Is the server path in `filename` correct? Does the file exist on the server?

## ⚙️ Configuration

This application reads configuration from a YAML file. If no file exists, a default file is created with sane defaults.

- Search order for `config.yaml`:
  1. `./config.yaml`
  2. `./config/config.yaml`
  3. `/etc/IngestListApiWrapper/config.yaml`

Default values are used when the file is created:

```yaml
server:
  host: "localhost"
  port: "8080"
database:
  driver: "sqlite"      # use "postgres" to connect to Postgres
  dsn: "tasks.db"       # for Postgres: e.g. "host=... user=... password=... dbname=... sslmode=disable"
  max_open: 10
  max_idle: 5
task_scheduler:
  interval: "*/30 * * * * *"  # every 30 seconds
  file_storage_path: "data"
  max_workers: 3
security:
  secret_key: "your-secret-key"
```

Tips:
- For PostgreSQL, set `database.driver: postgres` and provide a proper `dsn`.
- In containers, mount your config to `/app/config/config.yaml` or `/etc/IngestListApiWrapper/config.yaml`.
- The HTTP server listens on `server.host:server.port` (default `localhost:8080`).

## 📈 Metrics (Prometheus)

The service exposes Prometheus metrics at the endpoint `/metrics`. Besides standard Go/HTTP metrics, it publishes custom business and storage metrics.

- Endpoint: `GET /metrics`
- Behavior: right before responding, the server refreshes business metrics from the database and filesystem via the `MetricsService`, so values are up-to-date at scrape time.

Custom metrics exposed:

- `ingestlist_tasks_total` (gauge) — Total number of tasks in the database
- `ingestlist_tasks_by_status{status}` (gauge) — Number of tasks grouped by status
- `ingestlist_storage_files_count` (gauge) — Count of files in the storage directory
- `ingestlist_storage_size_bytes` (gauge) — Total size of files in the storage directory (bytes)
- `ingestlist_http_requests_total{method,endpoint,status}` (counter) — Total HTTP requests
- `ingestlist_http_request_duration_seconds{method,endpoint}` (histogram) — HTTP request duration; exposed as `_bucket`, `_sum`, `_count`

Quick check with curl:

```bash
curl -sS http://localhost:8080/metrics | head -n 40
```

Example Prometheus scrape configuration:

```yaml
scrape_configs:
  - job_name: "ilwrapper"
    static_configs:
      - targets: ["ilwrapper:8080"]  # or "localhost:8080" when running locally
    metrics_path: /metrics
    scrape_interval: 15s
```

Example Grafana/PromQL queries:

- Requests per second by endpoint (last 5m):
  ```
  sum by (endpoint) (rate(ingestlist_http_requests_total[5m]))
  ```
- Tasks by status:
  ```
  sum by (status) (ingestlist_tasks_by_status)
  ```

## 🔌 API

### ⚡️ Quickstart (TL;DR)

```bash
# 1) Set the base URL (adapt to your environment if needed)
export BASE_URL="http://dimagapps-ilwrapper"

# 2) Log in and store the token in $TOKEN
TOKEN=$(curl -sS -X POST "$BASE_URL/api/login" \
  -H "Content-Type: application/json" \
  -d '{"email":"user","password":"password"}' | jq -r '.token')

# 3) Example: run Remote‑Validate
curl -sS -X POST "$BASE_URL/api/create" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"filename":"data/example.pdf","type":"Validate"}' | jq .
```

### 🔗 Base URL

It’s convenient to set an environment variable for the base URL:

```bash
export BASE_URL="http://dimagapps-ilwrapper"
```

### 🧰 Prerequisites

- curl (for example requests)
- optional: jq (for pretty JSON output and easy token parsing)

### 🔐 Authentication

The API uses Bearer token authentication (JWT). After logging in, you’ll receive a JWT token that must be sent with subsequent
requests in the `Authorization` header: `Authorization: Bearer <TOKEN>`.

### 🧭 Endpoints

#### 🔑 Login

Authenticate a user and receive a JWT token.

```text
Endpoint: POST /api/login
Headers:
  Content-Type: application/json
Request Body:
  {
    "email": "user",
    "password": "password"
  }
Response (example 200):
  {
    "token": "<JWT>"
  }
Status codes: 200, 400, 401
```

```bash
# cURL example:
curl -sS -X POST "$BASE_URL/api/login" \
  -H "Content-Type: application/json" \
  -d '{"email":"user","password":"password"}' | jq .

# Store token in a variable (if jq is available):
TOKEN=$(curl -sS -X POST "$BASE_URL/api/login" \
  -H "Content-Type: application/json" \
  -d '{"email":"user","password":"password"}' | jq -r '.token')
```

#### ✅ Remote-Validate

Validate a file that already exists on the server.

```text
Endpoint: POST /api/create
Headers:
  Authorization: Bearer <TOKEN>
  Content-Type: application/json
Request Body:
  {
    "filename": "data/path/to/file.pdf",
    "type": "Validate"
  }
Status codes: 202 (Accepted), 400, 401, 404
```

```bash
# cURL example:
curl -sS -X POST "$BASE_URL/api/create" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"filename":"data/example.pdf","type":"Validate"}' | jq .
```

#### 🧪 Remote-Identify

Identify the file format of a file that already exists on the server.

```text
Endpoint: POST /api/create
Headers:
  Authorization: Bearer <TOKEN>
  Content-Type: application/json
Request Body:
  {
    "filename": "data/path/to/file.any",
    "type": "Identify"
  }
Status codes: 202 (Accepted), 400, 401, 404
```

```bash
# cURL example:
curl -sS -X POST "$BASE_URL/api/create" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"filename":"data/example.pdf","type":"Identify"}' | jq .
```

#### 📤 Local-Validate
Validate an uploaded file.
```text
Endpoint: POST /api/create
Headers:
  Authorization: Bearer <TOKEN>
  Content-Type: multipart/form-data
Form Data:
  type=Validate
  file=@/path/to/file.pdf
Status codes: 202 (Accepted), 400, 401
```
```bash
# cURL example:
curl -sS -X POST "$BASE_URL/api/create" \
  -H "Authorization: Bearer $TOKEN" \
  -F "type=Validate" \
  -F "file=@/path/to/your/file.pdf" | jq .
```

#### 🔎 Local-Identify
Identify the file format of an uploaded file.
```text
Endpoint: POST /api/create
Headers:
  Authorization: Bearer <TOKEN>
  Content-Type: multipart/form-data
Form Data:
  type=Identify
  file=@/path/to/file.any
Status codes: 202 (Accepted), 400, 401
```

```bash
# cURL example:
curl -sS -X POST "$BASE_URL/api/create" \
  -H "Authorization: Bearer $TOKEN" \
  -F "type=Identify" \
  -F "file=@/path/to/your/file.html" | jq .
```
#### 📜 Jobs
Fetch all jobs.

```text
Endpoint: GET /api/jobs
Headers:
  Authorization: Bearer <TOKEN>
Status codes: 200, 401
```
```bash
# cURL example:
curl -sS -X GET "$BASE_URL/api/jobs" \
  -H "Authorization: Bearer $TOKEN" | jq .
```

#### 🎯 Job by ID
Fetch a specific job by its ID.
```text
Endpoint: GET /api/job/{id}
Headers:
  Authorization: Bearer <TOKEN>
Path Parameter:
  id (integer): Job ID
Status codes: 200, 401, 404
```
```bash
# cURL example:
curl -sS -X GET "$BASE_URL/api/job/5" \
  -H "Authorization: Bearer $TOKEN" | jq .
```
### 🧩 Operation Types

The API supports the following operation types:

- Validate: Validate a file
- Identify: Identify the file format

### 🔄 Workflows

### Typical workflow with remote files

1. Log in (get token)
```bash
TOKEN=$(curl -sS -X POST "$BASE_URL/api/login" \
  -H "Content-Type: application/json" \
  -d '{"email":"user","password":"password"}' | jq -r '.token')
```

2. Validate file (remote)
```bash
curl -sS -X POST "$BASE_URL/api/create" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"filename":"data/example.pdf","type":"Validate"}' | jq .
```

3. Fetch jobs
```bash
curl -sS -X GET "$BASE_URL/api/jobs" \
  -H "Authorization: Bearer $TOKEN" | jq .
```

### Typical workflow with file upload

1. Log in
```bash
TOKEN=$(curl -sS -X POST "$BASE_URL/api/login" \
  -H "Content-Type: application/json" \
  -d '{"email":"user","password":"password"}' | jq -r '.token')
```

2. Upload and validate file
```bash
curl -sS -X POST "$BASE_URL/api/create" \
  -H "Authorization: Bearer $TOKEN" \
  -F "type=Validate" \
  -F "file=@/path/to/file.pdf" | jq .
```

3. Fetch a specific job
```bash
curl -sS -X GET "$BASE_URL/api/job/5" \
  -H "Authorization: Bearer $TOKEN" | jq .
```

### ⚖️ Differences between Remote and Local

#### Remote operations (JSON)

- The file must already exist on the server
- Use `Content-Type: application/json`
- File path is passed in the `filename` parameter

#### Local operations (multipart)

- The file is uploaded with the request
- Use `Content-Type: multipart/form-data`
- File is provided as the `file` form field


## 📦 Container Image (Docker/Podman)

This project is a containerized application. A multi-stage build is provided at `backend/Dockerfile`.

Build (Docker):

```bash
# from repository root
docker build -f backend/Dockerfile -t ilwrapper:latest .

# Optional: Configure HTTP(S) proxy for bundled tools at build-time
docker build -f backend/Dockerfile \
  --build-arg PROXY_HOST="your.proxy" \
  --build-arg PROXY_PORT="3128" \
  -t ilwrapper:latest .
```

Build (Podman):

```bash
podman build -f backend/Dockerfile -t ilwrapper:latest .
```

Run:

```bash
# Minimal run; maps 8080
docker run --rm -p 8080:8080 --name ilwrapper ilwrapper:latest

# With mounted config and data directory
docker run --rm -p 8080:8080 \
  -v "$(pwd)/conf/config.yaml:/app/config/config.yaml:ro" \
  -v "$(pwd)/data:/app/data" \
  --name ilwrapper ilwrapper:latest
```

Run with Podman (analogous):

```bash
podman run --rm -p 8080:8080 --name ilwrapper ilwrapper:latest
```

Notes:
- The image exposes port `8080` and runs the server binary.
- The container working directory is `/app`.

## ☸️ Kubernetes

The `k8s/` directory currently contains manifests to provision a PostgreSQL database:

- `k8s/1-database-secret.yml` (credentials as a Kubernetes Secret)
- `k8s/1-database-statefulset.yml` (PostgreSQL StatefulSet and Service)

Quick start:

```bash
# Create namespace (once)
kubectl create namespace ilwrapper

# Apply database manifests
kubectl apply -n ilwrapper -f k8s/1-database-secret.yml
kubectl apply -n ilwrapper -f k8s/1-database-statefulset.yml

# Check
kubectl get all -n ilwrapper
```

Next steps:
- Create and apply a Deployment and Service for ILWrapper pointing to your built image `ilwrapper:latest` (or a pushed registry tag).
- Configure the application to use Postgres by setting `database.driver: postgres` and a Postgres `dsn` via config (see Configuration above). Consider mounting the config with a ConfigMap or Secret.
- Optionally add an Ingress or LoadBalancer Service for external access.

## 💬 Feedback & Contributions

Feedback, ideas, and contributions are very welcome!

- Open an issue for bugs or feature requests.
- Fork the repo and submit a pull request for improvements.
- If unsure where to start, feel free to propose changes first in an issue.

Please keep changes small and focused where possible. Add context to PRs (what, why) and include examples or screenshots if relevant.

## 📄 License

This project is licensed under the MIT License. See the [LICENSE](./LICENSE) file for details.