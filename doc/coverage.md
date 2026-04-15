# Test coverage (`Makefile.ps1`)

This project measures coverage with **`Makefile.ps1`** at the repository root (**`-Target test-cover`**). Do not rely on `go test -cover ./...` without **`-coverpkg`**, or totals will mostly reflect test packages—not **`internal/...`**.

---

## Prerequisites

| Requirement | Notes |
|-------------|--------|
| **Windows PowerShell** 5.1+ or **PowerShell 7+** | Script is `Makefile.ps1`; run from a shell that can execute it. |
| **Working directory** | Implicit: script **`Set-Location`** to its directory (repo root next to `go.mod`). |
| **Go toolchain** | Matches `go.mod` (e.g. 1.26.x). **`go`** on **`PATH`**. |

```powershell
cd <path-to>\dominus-broker
.\Makefile.ps1 -Target test-cover
```

---

## What `test-cover` does

1. **`-coverpkg`** is built from **`go list ./internal/...`** and **`go list ./tests/...`**, joined as comma-separated import paths, with **`dominus-broker/internal/boostrap`** (package **`internal/boostrap`**) **removed** — bootstrap wiring is not part of the coverage budget here.

2. **`go test ./...`** with that **`-coverpkg`**, **`-race`**, **`-count=1`**, **`-coverprofile=coverage.out`**.

### Packages intentionally **not** in `-coverpkg`

| Area | Reason |
|------|--------|
| **`cmd/`** (e.g. `cmd/api`) | Entrypoint / `main`; not aggregated into the service coverage total. |
| **`config/`** | Root config package; not listed (same as not being under `internal/` + `tests/` lists). |
| **`mocks/`** | Generated stubs; excluded (not under `./internal/...` or `./tests/...`). |
| **`internal/boostrap/`** | Composition / wiring; excluded explicitly after `go list`. |

So the profile reflects **application, domain, infrastructure** under **`internal/`** (except bootstrap) plus **`tests/...`**, not mains nor generated mocks.

3. **`go tool cover -func coverage.out`** — per-function lines and **total** statement %.

There is **no** `timeout` flag in the Makefile target; add one in CI if jobs need a hard cap.

### Compared to the former `run_coverage.ps1`

An older script (**removed**) ran **`go test ./tests/...`** only and built **`-coverpkg`** from **`go list ./...`** with exclusions (**`mocks`**, **`docs`**). **`Makefile.ps1`** runs **`go test ./...`** and builds **`-coverpkg`** from **`go list ./internal/...`** (minus **`boostrap`**) + **`go list ./tests/...`**. Totals can differ from older baselines; refresh the table below after workflow changes.

---

## Artifacts

| File | Git | Description |
|------|-----|-------------|
| `coverage.out` | **Ignored** (`*.out` in `.gitignore`) | Profile for `go tool cover`. |

### Optional HTML report

```powershell
go tool cover -html=coverage.out -o coverage.html
```

Open `coverage.html` in a browser for line-by-line highlighting.

---

## Baseline report (instrumented with `Makefile.ps1 -Target test-cover`)

> **Regenerate:** after meaningful test or production code changes, run **`.\Makefile.ps1 -Target test-cover`** and update this section (especially **total** and any rows you care about tracking).

| Field | Value |
|--------|--------|
| **Command** | `.\Makefile.ps1 -Target test-cover` from repo root |
| **Captured** | 2026-04-01 (re-run after code changes and refresh this doc) |
| **Total statements** | **91.8%** (then `go tool cover -func coverage.out`) |
| **Output of** | `go tool cover -func coverage.out` (last line: `total: (statements) 91.8%`) |

### By area (summary)

| Area | Coverage notes |
|------|----------------|
| **Broker use cases** | `StreamClientConn`, `StreamServerConn`, `StreamBiConn`, `NewBroker` — **100%** |
| **Sqs use cases** | `Producer`, `Consumer`, `Ack`, `NewSqs` — **100%** |
| **Domain `Message`** | `NewMessageWithID`, `NewMessage`, getters/setters, `checkValidFormatID` — **100%** |
| **gRPC mappers** | All listed methods — **100%** |
| **gRPC interceptors** | `NewInterceptor`, `UnaryAuthInterceptor`, `StreamAuthInterceptor` — **100%** |
| **gRPC middleware** | `ApiToken`, `UnaryLog`, `StreamLog`, `LogErrors`, `NewMiddleware`, **`IdPotency`** — **100%** (last full run) |
| **gRPC inbound broker** | `ClientStream`, `ServerStream`, `BidirectionalStream`, `NewBrokerAPI` — **100%** |
| **gRPC inbound Sqs** | **100%** on listed handlers |
| **gRPC outbound** | **`ClientStream` — 78.8%**, **`ServerStream` — 90.6%**, **`BidirectionalStream` — 74.4%** (retry / error branches partially hit) |
| **Redis `cchecker`** | **`SaveConsumer` / `CheckConsumer` — 100%**; **`NewCheckerClient` — 71.4%** (e.g. TLS branch not hit) |
| **Redis `cmemory`** | **`AckMessage` / `Group` — 100%**; **`GetMessage` — 90.9%**; **`SendMessage` — 83.3%**; **`NewMemoryClient` — 71.4%** |
| **Event** | **`CheckID` / `NewEvent` / `typeLog` / `cmdLog` — 100%**; **`WriteLog` — 88.9%**; **`clientLog` — 0.0%** (`LOGCLIENT` path) |
| **fasthttp monitor** | **`getHealthCheck` / `collectMetrics` / `NewMonitorAPI` / `convertToHTTP` — 100%**; **`getMetrics` — 75%**; **`WriteHeader` — 0%** |
| **fasthttp middleware** | **API token** and **host** `CheckMiddleware` — **100%**; **chain** (`middlewares.go`) — **100%** (last full run) |

### Per-function detail (`go tool cover -func coverage.out`)

Snapshot from a recent **`test-cover`** run; re-run locally after changes—line-level % can drift while **total** above stays the tracked baseline until you refresh it.

```
dominus-broker/internal/application/usecases/broker/factory.go:18:                    NewBroker               100.0%
dominus-broker/internal/application/usecases/broker/stream_bi_conn_service.go:10:         StreamBiConn            100.0%
dominus-broker/internal/application/usecases/broker/stream_client_conn_service.go:10:       StreamClientConn        100.0%
dominus-broker/internal/application/usecases/broker/stream_server_conn_service.go:10:       StreamServerConn        100.0%
dominus-broker/internal/application/usecases/sqs/ack_service.go:10:                          Ack                     100.0%
dominus-broker/internal/application/usecases/sqs/consumer_service.go:9:                     Consumer                100.0%
dominus-broker/internal/application/usecases/sqs/factory.go:20:                     NewSQS                  100.0%
dominus-broker/internal/application/usecases/sqs/producer_service.go:10:                    Producer                100.0%
dominus-broker/internal/domain/entities/sqs_message.go:16:                            NewMessageWithID        100.0%
dominus-broker/internal/domain/entities/sqs_message.go:24:                            NewMessage              100.0%
dominus-broker/internal/domain/entities/sqs_message.go:28:                            SetMessage              100.0%
dominus-broker/internal/domain/entities/sqs_message.go:31:                            GetMessage              100.0%
dominus-broker/internal/domain/entities/sqs_message.go:35:                            SetMessageId            100.0%
dominus-broker/internal/domain/entities/sqs_message.go:43:                            GetMessageId            100.0%
dominus-broker/internal/domain/entities/sqs_message.go:47:                            GetCreatedAt            100.0%
dominus-broker/internal/domain/entities/sqs_message.go:50:                            SetCreateAt             100.0%
dominus-broker/internal/domain/entities/sqs_message.go:54:                            checkValidFormatID      100.0%
dominus-broker/internal/infrastructure/event/event.go:26:                            NewEvent                100.0%
dominus-broker/internal/infrastructure/event/event.go:39:                            WriteLog                88.9%
dominus-broker/internal/infrastructure/event/event.go:60:                            CheckID                 100.0%
dominus-broker/internal/infrastructure/event/event.go:69:                            typeLog                 100.0%
dominus-broker/internal/infrastructure/event/event.go:78:                            clientLog               0.0%
dominus-broker/internal/infrastructure/event/event.go:81:                            cmdLog                  100.0%
dominus-broker/internal/infrastructure/fasthttp/inbound/monitor_api.go:31:           Header                  100.0%
dominus-broker/internal/infrastructure/fasthttp/inbound/monitor_api.go:35:           Write                   100.0%
dominus-broker/internal/infrastructure/fasthttp/inbound/monitor_api.go:39:           WriteHeader             0.0%
dominus-broker/internal/infrastructure/fasthttp/inbound/monitor_api.go:50:           NewMonitorAPI           100.0%
dominus-broker/internal/infrastructure/fasthttp/inbound/monitor_api.go:61:           convertToHTTP           100.0%
dominus-broker/internal/infrastructure/fasthttp/inbound/monitor_api.go:72:           collectMetrics          100.0%
dominus-broker/internal/infrastructure/fasthttp/inbound/monitor_api.go:83:           getMetrics              75.0%
dominus-broker/internal/infrastructure/fasthttp/inbound/monitor_api.go:95:           getHealthCheck          100.0%
dominus-broker/internal/infrastructure/fasthttp/inbound/monitor_api.go:101:          path                    100.0%
dominus-broker/internal/infrastructure/fasthttp/middlewares/api_middleware.go:20:    NewMiddlewareApiToken   100.0%
dominus-broker/internal/infrastructure/fasthttp/middlewares/api_middleware.go:27:    CheckMiddleware         100.0%
dominus-broker/internal/infrastructure/fasthttp/middlewares/host_allowed.go:17:      NewMiddlewareHost       100.0%
dominus-broker/internal/infrastructure/fasthttp/middlewares/host_allowed.go:24:    CheckMiddleware         100.0%
dominus-broker/internal/infrastructure/fasthttp/middlewares/middlewares.go:22:     NewMiddleware           100.0%
dominus-broker/internal/infrastructure/fasthttp/middlewares/middlewares.go:28:     AddMiddleware           100.0%
dominus-broker/internal/infrastructure/fasthttp/middlewares/middlewares.go:32:     Middlewares             100.0%
dominus-broker/internal/infrastructure/grpc/inbound/broker_v1.3.7.go:24:           NewBrokerAPI            100.0%
dominus-broker/internal/infrastructure/grpc/inbound/broker_v1.3.7.go:30:           ClientStream            100.0%
dominus-broker/internal/infrastructure/grpc/inbound/broker_v1.3.7.go:50:           ServerStream            100.0%
dominus-broker/internal/infrastructure/grpc/inbound/broker_v1.3.7.go:57:           BidirectionalStream     100.0%
dominus-broker/internal/infrastructure/grpc/inbound/sqs_v1.3.7.go:22:              NewSqsAPI               100.0%
dominus-broker/internal/infrastructure/grpc/inbound/sqs_v1.3.7.go:28:              Producer                100.0%
dominus-broker/internal/infrastructure/grpc/inbound/sqs_v1.3.7.go:43:              Consumer                100.0%
dominus-broker/internal/infrastructure/grpc/inbound/sqs_v1.3.7.go:68:              Ack                     100.0%
dominus-broker/internal/infrastructure/grpc/mappers/mappers.go:14:                  NewBiStreamConn         100.0%
dominus-broker/internal/infrastructure/grpc/mappers/mappers.go:20:                  Recv                    100.0%
dominus-broker/internal/infrastructure/grpc/mappers/mappers.go:24:                  Send                    100.0%
dominus-broker/internal/infrastructure/grpc/mappers/mappers.go:31:                  Context                 100.0%
dominus-broker/internal/infrastructure/grpc/mappers/mappers.go:39:                  NewClientStreamContext  100.0%
dominus-broker/internal/infrastructure/grpc/mappers/mappers.go:45:                  Recv                    100.0%
dominus-broker/internal/infrastructure/grpc/mappers/mappers.go:49:                  Context                 100.0%
dominus-broker/internal/infrastructure/grpc/mappers/mappers.go:57:                  NewServerStreamContext  100.0%
dominus-broker/internal/infrastructure/grpc/mappers/mappers.go:63:                  Send                    100.0%
dominus-broker/internal/infrastructure/grpc/mappers/mappers.go:70:                  Context                 100.0%
dominus-broker/internal/infrastructure/grpc/middlewares/interceptors.go:24:       NewInterceptor          100.0%
dominus-broker/internal/infrastructure/grpc/middlewares/interceptors.go:31:       UnaryAuthInterceptor    100.0%
dominus-broker/internal/infrastructure/grpc/middlewares/interceptors.go:38:       StreamAuthInterceptor   100.0%
dominus-broker/internal/infrastructure/grpc/middlewares/middlewares.go:32:        NewMiddleware           100.0%
dominus-broker/internal/infrastructure/grpc/middlewares/middlewares.go:40:        ApiToken                100.0%
dominus-broker/internal/infrastructure/grpc/middlewares/middlewares.go:57:        UnaryLog                100.0%
dominus-broker/internal/infrastructure/grpc/middlewares/middlewares.go:62:        StreamLog               100.0%
dominus-broker/internal/infrastructure/grpc/middlewares/middlewares.go:67:        LogErrors               100.0%
dominus-broker/internal/infrastructure/grpc/middlewares/middlewares.go:85:        IdPotency               100.0%
dominus-broker/internal/infrastructure/grpc/outbound/client_v1.3.7.go:22:          NewGrpClient            100.0%
dominus-broker/internal/infrastructure/grpc/outbound/client_v1.3.7.go:29:          ClientStream            78.8%
dominus-broker/internal/infrastructure/grpc/outbound/client_v1.3.7.go:79:          ServerStream            90.6%
dominus-broker/internal/infrastructure/grpc/outbound/client_v1.3.7.go:132:         BidirectionalStream     74.4%
dominus-broker/internal/infrastructure/redis/cchecker/outbound.go:19:             NewCheckerClient        71.4%
dominus-broker/internal/infrastructure/redis/cchecker/outbound.go:64:             SaveConsumer            100.0%
dominus-broker/internal/infrastructure/redis/cchecker/outbound.go:74:             CheckConsumer           100.0%
dominus-broker/internal/infrastructure/redis/cmemory/outbound.go:21:              NewMemoryClient         71.4%
dominus-broker/internal/infrastructure/redis/cmemory/outbound.go:66:              SendMessage             83.3%
dominus-broker/internal/infrastructure/redis/cmemory/outbound.go:83:              AckMessage              100.0%
dominus-broker/internal/infrastructure/redis/cmemory/outbound.go:90:              GetMessage              90.9%
dominus-broker/internal/infrastructure/redis/cmemory/outbound.go:118:             Group                   100.0%
total:                                                                             (statements)            91.8%
```

### Packages not listed above

`go tool cover -func` only prints functions that appear in the profile. Code that is **never reached** by any test may be absent or folded into totals. **`cmd/`**, **`config/`**, and **`internal/boostrap/`** are **out of `-coverpkg`** (see table above), so they do not appear in this report. Raising **total** for the instrumented set usually means more tests against **`internal/...`**; covering bootstrap or `main` would require a deliberate change to **`-coverpkg`** or separate integration runs.

---

## CI suggestion

Mirror **`Makefile.ps1 -Target test-cover`** (or **`check-coverage`**) with the same **`go test`** flags and **`-coverpkg`**, or run **`.\Makefile.ps1`** on Windows agents. Archive **`coverage.out`** and optionally fail if **total** drops below a threshold (**`check-coverage`** uses **`-CoverageMinPct`**; default **85** in the script—align with team policy).
