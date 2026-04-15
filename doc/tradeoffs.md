# Architectural and operational trade-offs

This document captures deliberate choices, known limitations, and engineering trade-offs in **Dominus**. It complements **[broker-streaming.md](broker-streaming.md)**, **[concurrency.md](concurrency.md)**, **[grpc-security.md](grpc-security.md)**, and **[configuration.md](configuration.md)**.

---

## Layering and module layout

| Trade-off | Benefit | Cost |
|-----------|---------|------|
| **`internal/` + `tests/...` packages** | Production code cannot be imported by other modules; tests stay as first-class packages with clear `dominus-broker/tests/...` imports. | Extra path depth; no `_test.go` next to handlers (must open two trees when editing). |
| **Thin domain (`entities`, `repositories` interfaces)** | Swappable Redis / gRPC adapters; mocks in `mocks/`. | More boilerplate and `mockgen` churn when interfaces change. |
| **Application use cases as facades** | Single place for orchestration (broker streams, Sqs-like flows). | Complex stream code (channels, defers) is harder to follow than a single linear handler. |

---

## Dual HTTP stacks (gRPC + fasthttp)

| Trade-off | Benefit | Cost |
|-----------|---------|------|
| **gRPC for Broker + Sqs APIs** | Strong contracts (protobuf), streaming, interceptors, OTel/Prometheus hooks. | Operational complexity (ports, TLS, proto version lock to `dominus-proto-definition`). |
| **fasthttp for monitor/health only** | Very small surface; fits Prometheus scrape and liveness. | Two server types to secure, tune, and document; `fasthttp` API differs from `net/http` idioms. |
| **TLS gated by cert files on disk** | Simple local/dev story (plain HTTP/gRPC when files missing). | Behaviour is implicit; misconfigured prod can run insecure without failing fast at startup. |

---

## Redis as queue and idempotency store

| Trade-off | Benefit | Cost |
|-----------|---------|------|
| **Redis Streams (`cmemory`) instead of AWS SQS** | Single deployment artifact; predictable latency in owned infra. | Not Sqs-semantics-complete; scaling and durability are Redis-cluster concerns, not Amazon's. |
| **Separate DB indices (`MemoryDB`, `CheckerDB`)** | Logical separation of queue vs idempotency keys. | Two pools/config knobs; operational mistakes (wrong DB) are harder to spot. |
| **Unary idempotency: `EXISTS` then async `SET NX`** | Fast happy path; handler proceeds immediately after spawn. | **Logical race**: concurrent duplicate keys can both pass before either `SET NX` completes (see **[grpc-security.md](grpc-security.md)**). Stronger deduplication needs an atomic reserve (e.g. synchronous `SET NX` before handler). |

---

## Broker outbound client (`grpc/outbound`)

| Trade-off | Benefit | Cost |
|-----------|---------|------|
| **Per-message `go cls(m)` fan-out** | Parallel sends to all subscriber URLs. | **Unbounded goroutine growth** under bursty ingress; CPU and scheduler pressure. |
| **Shared `[]byte` across fan-out goroutines** | No copy per subscriber; lower allocation. | Payload must stay **read-only** until all `Send` calls for that message complete; buffer reuse upstream causes races. |
| **Mutex only around stream pointer / error, not around `Send`** | Avoids deadlocks with blocking gRPC I/O. | Subtle ordering / TOCTOU if stream is replaced while another path still holds a stale reference briefly. |
| **Bidirectional `done` vs `closed` channels** | Clear separation: worker completion vs outbound drain of `provMsg`. | Easy to misuse (e.g. `defer close(closed)` at function entry races outbound `tx` send — see **[broker-streaming.md](broker-streaming.md)**). |

---

## Concurrency and shutdown

| Trade-off | Benefit | Cost |
|-----------|---------|------|
| **`StreamServerConn`: forward goroutine + `done` wait loop** | Avoids `select` race where buffered `done` signals are consumed before `stream` payloads are forwarded. | Two goroutines per RPC; must keep `cancel()` / error paths consistent. |
| **Async `go` logging in middleware and Redis** | Request path returns quickly; `slog` is concurrent-safe. | Log ordering not aligned with wall-clock request order; harder to debug without correlation IDs. |
| **Integration tests cancel context to tear down bidi streams** | Unblocks outbound workers when TCP peers keep streams open (echo peers). | Tests do not prove production-only `CloseSend`-only shutdown unless peers cooperate. |

---

## Configuration

| Trade-off | Benefit | Cost |
|-----------|---------|------|
| **Prod: `APP_CONFIG` JSON in env** | 12-factor friendly for containers. | Large blobs in env; typos panic at startup (`panic` on unmarshal). |
| **Dev: Viper file `./../env/env.dev.json`** | Easy local iteration. | Path relative to CWD; wrong working directory breaks startup. |
| **`panic` on missing config / bad Redis ping** | Fail fast in bootstrap. | No graceful degradation; harder to run partial stacks for experiments. |

---

## Observability

| Trade-off | Benefit | Cost |
|-----------|---------|------|
| **Prometheus + gRPC middleware metrics** | Standard dashboards and RED-style views. | Cardinality and histogram defaults may need tuning under high method cardinality. |
| **OpenTelemetry stats handlers on gRPC** | Distributed tracing hooks. | Additional deps and CPU; sampling not documented in-repo. |
| **Monitor metrics scrape triggers host metrics collection** | Fresh CPU/memory gauges on each `/metrics` scrape. | Extra work per scrape; noisy under aggressive Prometheus intervals. |

---

## Testing strategy

| Trade-off | Benefit | Cost |
|-----------|---------|------|
| **`Makefile.ps1` / `-coverpkg` over `internal/...` and `tests/...`** | Meaningful % across `internal/...`. | Slower runs; `-race` multiplies time. |
| **bufconn for ingress, TCP for broker peers** | Mix of hermetic and realistic I/O. | Flaky tests if ports/timeouts are tight; different failure modes than full mesh bufconn. |
| **miniredis for Redis tests** | Fast, no Docker required for CI. | Behaviour diverges from real Redis / cluster edge cases (ACL, TLS to Redis, etc.). |

---

## Security and API surface

| Trade-off | Benefit | Cost |
|-----------|---------|------|
| **API key via metadata (hashed compare)** | Simple shared secret. | No per-user identity; rotation is operational (deploy new key everywhere). |
| **Idempotency on unary only** | Clear scope. | Streaming RPCs have no idempotency middleware in the current chain. |
| **Host allowlist via CIDR** | Basic network guard for REST. | gRPC port not covered by the same rule; IP spoofing / proxy headers not addressed here. |

---

## Bootstrap and lifecycle

| Trade-off | Benefit | Cost |
|-----------|---------|------|
| **`RunApp` blocks on signal; `GracefulStop` gRPC** | Clean drain for in-flight RPCs when shutdown is orderly. | `log.Fatal` on signal path; `close(st)` after shutdown is unusual; listener errors use `log.Fatal` in a goroutine (process exit). |
| **Single binary `cmd/api`** | Simple deploy. | No separate worker binary; scaling story is “more replicas”, not role-based split. |

---

## Dependencies and coupling

| Trade-off | Benefit | Cost |
|-----------|---------|------|
| **Pinned protobuf module `dominus-proto-definition@v1.3.7`** | Reproducible builds. | Server and clients must upgrade protos in lockstep. |
| **jsoniter, gopsutil, etc.** | Performance / system metrics. | More supply-chain surface; upgrades need regression runs. |

---

## When to revisit these trade-offs

- **Strict idempotency** under concurrent retries → synchronous atomic key claim in Redis.
- **Sustained high broker fan-out** → limit concurrency, batching, or back-pressure from `provMsg`.
- **Production hardening** → explicit TLS-required mode, structured shutdown without `Fatal`, Redis Sentinel/Cluster runbooks.
- **Broader REST API** → reconsider fasthttp vs `net/http` ecosystem and middleware portability.
