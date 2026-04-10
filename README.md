# Dominus

<p align="center">
  <a href="https://go.dev/" target="_blank" rel="noopener noreferrer">
    <img
      src="https://go.dev/blog/go-brand/Go-Logo/PNG/Go-Logo_LightBlue.png"
      width="120"
      alt="Go Logo"
    />
  </a>
</p>

<p align="center">
  <img src="https://img.shields.io/badge/Go-1.26.1-00ADD8?logo=go&logoColor=white" alt="Go 1.26.1" />
  <img src="https://img.shields.io/badge/gRPC-1.79.3-4285F4?logo=grpc&logoColor=white" alt="gRPC" />
  <a href="https://github.com/MBI-88/dominus-proto-definition" target="_blank" rel="noopener noreferrer">
    <img src="https://img.shields.io/badge/protos-v1.3.7-blue" alt="dominus-proto-definition v1.3.7" />
  </a>
  <img src="https://img.shields.io/badge/license-Proprietary-lightgrey" alt="License" />
</p>

<p align="center">
  A <strong>Go</strong> service with <strong>gRPC</strong> (broker + SQS-like queue), <strong>HTTP</strong> (fasthttp), <strong>Redis</strong>, <strong>Prometheus</strong>, and <strong>OpenTelemetry</strong>.
</p>

| | |
|---|---|
| **Module** | `dominus-broker` |
| **Go** | 1.26.1 (see `go.mod`) |
| **Contracts** | [dominus-proto-definition](https://github.com/MBI-88/dominus-proto-definition) v1.3.7 |

---

## What it does

- **gRPC broker** (`BrokerAPI`): client, server, and bidirectional streaming between nodes; the use case fans out to subscriber URLs via the gRPC client in `internal/infrastructure/grpc/outbound`.
- **SQS-like gRPC** (`SqsAPI`): producer, consumer, and ack on Redis Streams (`cmemory`) with duplicate control via `cchecker`.
- **REST (fasthttp)**: monitor/health endpoints behind API token and allowed-host middleware.
- **Security / ops**: optional TLS for gRPC and REST when `cert_config` files exist; gRPC API key and unary idempotency middleware.

---

## Documentation

Detailed design and file-level notes live under **[`doc/`](doc/)**:

| Guide | Description |
|-------|-------------|
| [doc/README.md](doc/README.md) | Index of all technical docs |
| [doc/broker-streaming.md](doc/broker-streaming.md) | Broker gRPC streams, outbound fan-out, bidi shutdown, channel lifecycle, payload sharing |
| [doc/concurrency.md](doc/concurrency.md) | Race detector, bidi integration teardown, channel rules, idempotency |
| [doc/configuration.md](doc/configuration.md) | Viper, JSON schema, dev vs prod |
| [doc/grpc-security.md](doc/grpc-security.md) | Middleware, API key, idempotency, TLS, metrics |
| [doc/redis.md](doc/redis.md) | Checker (idempotency) and cmemory (streams) |
| [doc/sqs-use-cases.md](doc/sqs-use-cases.md) | Producer / consumer / ack over gRPC |
| [doc/http-monitor.md](doc/http-monitor.md) | fasthttp monitor API and middleware |
| [doc/testing.md](doc/testing.md) | `tests/cases` vs `tests/integration`, commands |
| [doc/tradeoffs.md](doc/tradeoffs.md) | Trade-offs: architecture, broker concurrency, Redis, security, ops |
| [doc/coverage.md](doc/coverage.md) | **`Makefile.ps1 -Target test-cover`**, baseline **91.8%**, per-function breakdown |
| [doc/terraform.md](doc/terraform.md) | Local stack: Terraform + Docker provider, modules, **`Makefile.ps1`** Terraform targets |

---

## Repository layout

```
cmd/api/              # Entrypoint (flags -prod, -banner)
config/               # Configuration types
doc/                  # Technical documentation (see Documentation above)
env/                  # env.dev.json, env.prod.json, entrypoint
internal/
  application/        # Use cases (broker, sqs)
  domain/             # Entities and ports (repositories)
  infrastructure/     # gRPC in/out, Redis, fasthttp, event, enum
  boostrap/           # Server bootstrap and wiring
mocks/                # go.uber.org/mock (generated)
Makefile.ps1          # Windows PowerShell: build, test, coverage, lint, audit (see Tests / Development)
tests/
  cases/              # Layered tests (unit / component)
  integration/        # broker_flow_test, sqs_flow_test, etc.
terraform/            # Terraform: Docker provider, dev stack (Dominus, Redis, sidecar, Prometheus, Grafana)
```

---

## Requirements

- Go **1.26.1** or compatible with `go.mod`.
- **Redis** reachable per `redis_config` (streams + checker DB).
- Certificates at the paths in `cert_config` for HTTPS/TLS on gRPC and REST; if missing, startup uses **HTTP / insecure** in development.

---

## Configuration

- **Development**: Viper reads `env/env.dev.json` (path relative to the process; often run from `cmd/api`). For containers, see **[doc/terraform.md](doc/terraform.md)** (`APP_CONFIG` or extend mounts).
- **Production**: `APP_CONFIG` environment variable with inline JSON (see `config.NewConfig`).

Main root JSON fields:

| Block | Purpose |
|--------|---------|
| `grpc_config` | gRPC port, `connection_key` (API key middleware) |
| `rest_config` | REST port, `api_token`, allowed origins |
| `cert_config` | Paths to cert, key, and CA |
| `redis_config` | Host, port, DBs, stream/group, idempotency TTL |
| `infra_config` | Log mode (`cmd`, etc.) |

> Keep JSON aligned with tags in `config/config.go` (`redis_config` at root, not nested under another key).

---

## Running locally

From the module root:

```bash
go run ./cmd/api
```

Flags:

| Flag | Description |
|------|-------------|
| `-prod` | Production mode (config via `APP_CONFIG`) |
| `-banner` | Show or hide banner (default `true`) |

Development example without banner:

```bash
go run ./cmd/api -banner=false
```

Default ports come from `env.dev.json` (e.g. gRPC **5000**, REST **8000**).

---

## Infrastructure (Terraform)

A full **local** stack (broker, Redis, Nginx sidecar, Prometheus, Grafana) is defined under **`terraform/`** using the Docker provider. See **[doc/terraform.md](doc/terraform.md)** for prerequisites, raw CLI usage, variable defaults, and **`Makefile.ps1`** wrappers (**`-Target terraform-init`**, **`terraform-plan`**, **`terraform-apply`**, **`terraform-destroy`**, etc.).

---

## Tests

Test packages live under `tests/...` (separate from `internal` packages):

```bash
go test -race -count=1 ./tests/...
```

Use **`-race`** (race detector) and **`-count=1`** (disable test result cache) to match **`Makefile.ps1 -Target test`** / **`test-cover`** and avoid cached passes hiding races. For the full module: **`go test -race -count=1 ./...`**.

- **`tests/cases/`**: domain-focused unit tests (broker, sqs, gRPC, Redis, fasthttp, middlewares).
- **`tests/integration/broker_flow_test/`**: ClientStream / ServerStream / BidirectionalStream flows with bufconn + TCP peers; mocks limited to logging/checker and simulated downstream peers.
- **`tests/integration/sqs_flow_test/`**: SqsAPI Producer / Consumer / Ack with miniredis and real stream adapter (see **[doc/sqs-use-cases.md](doc/sqs-use-cases.md)**).

### Coverage (`Makefile.ps1`)

Meaningful totals for **`internal/...`** need **`-coverpkg`** (see **[doc/coverage.md](doc/coverage.md)**). From the repository root:

```powershell
.\Makefile.ps1 -Target test-cover      # writes coverage.out + go tool cover -func
.\Makefile.ps1 -Target check-coverage  # same, then fails if total is below -CoverageMinPct (see script default)
```

- **test-cover** builds **`-coverpkg`** from **`go list ./internal/...`** (minus **`internal/boostrap`**) + **`go list ./tests/...`**, so **`cmd/`**, **`config/`**, **`mocks/`**, and bootstrap are **not** in the profile (see **[doc/coverage.md](doc/coverage.md)**).
- **Latest recorded total: 91.8%** statements in **[doc/coverage.md](doc/coverage.md)** (refresh that doc after large changes).

`coverage.out` is gitignored (`*.out`).

---

## Development

- **Mocks**: `go generate` / `mockgen` per code comments; output in `mocks/`.
- **Lint**: `.golangci.yaml`, optional hooks in `.pre-commit-config.yaml`.
- **Windows automation**: **`Makefile.ps1`** at repo root (`-Target build`, `test`, `test-cover`, `fmt`, `lint`, `vuln`, `audit`, `deploy-check`, and Terraform: **`terraform-init`**, **`terraform-fmt`**, **`terraform-validate`**, **`terraform-plan`**, **`terraform-apply`**, **`terraform-destroy`**, **`terraform-output`** — see **[doc/terraform.md](doc/terraform.md)**). **`go`**, **`golangci-lint`**, **`govulncheck`**, and **`terraform`** (for those targets) must already be on **`PATH`**; the script does not modify `PATH`.

- **Dependency sync**: If module or workspace sums are inconsistent (for example `go.mod` requires `github.com/MBI-88/dominus-proto-definition v1.3.7` but `go.work.sum` lists v1.3.6), sync dependencies with:

```bash
go get github.com/MBI-88/dominus-proto-definition@v1.3.7
go work sync
```

Run the commands from the repository root and commit any changes to `go.mod`, `go.sum`, or `go.work.sum` as appropriate.

---

## License

See the `LICENSE` file in the repository.

For commercial licensing, redistribution, or enterprise terms, contact:
`licensing@dominus-broker.example` (see `LICENSE` for details).
