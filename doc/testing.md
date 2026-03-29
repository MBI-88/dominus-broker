# Testing layout

Go test packages live under **`tests/...`** so they import **`internal/...`** as normal dependencies of the same module (`dominus-project`), without placing `_test.go` files next to production packages.

---

## `tests/cases/`

| Area | Path (example) | Typical focus |
|------|----------------|----------------|
| Broker use cases | `tests/cases/application/use_cases/broker_test/` | `StreamClientConn`, `StreamServerConn`, `StreamBiConn` with mocks |
| SQS use cases | `tests/cases/application/use_cases/sqs_test/` | Producer / Consumer / Ack with `MockMemoryClient`; payload and message-ID assertions; invalid Ack IDs without calling Redis |
| gRPC inbound | `tests/cases/infrastructure/grpc/inbound_test/` | `BrokerAPI`, `SqsAPI` handlers with bufconn / mocks; SqsAPI Ack invalid-ID path uses real `sqs.NewSQS` + mock client |
| gRPC middleware | `tests/cases/infrastructure/grpc/middleware_test/` | `ApiToken`, interceptors, logging |
| gRPC mappers | `tests/cases/infrastructure/grpc/mappers_test/` | DTO wrappers around streams |
| Redis | `tests/cases/infrastructure/redis/cchecker_test/`, `cmemory_test/` | miniredis, real `SetArgs` / streams |
| fasthttp | `tests/cases/infrastructure/fasthttp/` | monitor API, API-token and host middleware |
| Event | `tests/cases/infrastructure/event_test/` | `event` package behavior |

**Mocks** live in **`mocks/`** (generated with `go.uber.org/mock` / `mockgen`) for `BrokerClient`, `MemoryClient`, `Event`, DTOs, etc.

---

## `tests/integration/`

| Package | Purpose |
|---------|---------|
| `broker_flow_test` | Ingress gRPC over **bufconn**, real `broker` + **`outbound.NewGrpClient`** where applicable; downstream “peers” on **TCP** (`helpers_test.go`) or mocks only for **`Event`** / **`CheckerClient`**; covers ClientStream, ServerStream, BidirectionalStream. **BidirectionalStream:** see **`doc/concurrency.md`** (*Integration tests: CloseSend vs context cancel*) — use **`context.Cancel`** after assertions, not **`CloseSend()`** alone, when peers keep their gRPC stream open (e.g. echo loop). |
| `sqs_flow_test` | End-to-end **SqsAPI** over bufconn with **miniredis** and real **`cmemory`**: Producer (stream payload + errors), Consumer (ordering, empty stream, validation), Ack after consume + pending verification. |

Integration tests are slower and assert **wiring + transport**; keep **unit tests** in `tests/cases/` for branches, errors, and focused contracts (e.g. Ack message-ID format at the use-case layer).

---

## Commands

Run all tests (no full-module coverage attribution):

```bash
go test ./tests/...
```

Verbose / single package:

```bash
go test -v ./tests/integration/broker_flow_test/...
```

---

## Coverage (mandatory: `run_coverage.ps1`)

**Do not** treat `go test -cover ./tests/...` alone as project coverage: it does not set `-coverpkg` across `internal/...`.

Use the repo script from the **repository root**:

```powershell
.\run_coverage.ps1
```

That produces **`coverage.out`** and prints **`go tool cover -func`** including the **total statement %** (baseline **90.0%** when last recorded).

Full script walkthrough, `-coverpkg` package list rules, HTML report command, **per-function table**, and gaps (`IdPotency` 0%, bootstrap untested, etc.): **[coverage.md](coverage.md)**.

Engineering trade-offs (testing vs production realism, coverage scope): **[tradeoffs.md](tradeoffs.md)**.

---

## Generating mocks

When an interface in `internal/...` changes, regenerate with **`mockgen`** (see `//go:generate` comments if present). Output stays in **`mocks/`**.
