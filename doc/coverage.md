# Test coverage (`run_coverage.ps1`)

This project treats **`run_coverage.ps1`** (repository root) as the **canonical** way to measure coverage. Do not rely on a plain `go test -cover` without `-coverpkg`, or you will only see coverage inside test packages—not `internal/...`.

---

## Prerequisites

| Requirement | Notes |
|-------------|--------|
| **Windows PowerShell** or **PowerShell 7+** | Script is `.ps1`; run from a shell that can execute it. |
| **Working directory** | Repository **root** (`dominus-project/`), same folder as `go.mod` and `run_coverage.ps1`. |
| **Go toolchain** | Matches `go.mod` (e.g. 1.26.x). `go` on `PATH`. |

```powershell
cd <path-to>\dominus-project
.\run_coverage.ps1
```

---

## What the script does (step by step)

1. **`$excluded`** — Package path substrings to **omit** from `-coverpkg`:
   - `docs`
   - `mocks`

   Any `go list ./...` package whose path matches one of these patterns is dropped.  
   **Note:** the documentation folder is `doc/`, not `docs/`; `doc` is **not** excluded and would be included if it ever contained Go packages.

2. **`go list ./...`** — Enumerates all module packages.

3. **Filter** — Keeps packages that do **not** match any excluded pattern.

4. **Print** — Lists every package that will be instrumented (helpful for debugging “why is X not covered?”).

5. **`$coverpkg`** — Joins the kept packages with commas for a single `-coverpkg=...` argument. That forces the compiler to record statement hits in **production code** when tests under `./tests/...` run.

6. **`go test`** — Exact invocation:

   ```text
   go test -timeout 120s -race -coverpkg="<comma-separated packages>" -covermode=atomic -coverprofile=coverage.out ./tests/...
   ```

   | Flag | Purpose |
   |------|---------|
   | `-timeout 120s` | Avoid hanging builds on stuck tests. |
   | `-race` | Race detector; failures indicate real data races. |
   | `-coverpkg=...` | Attribute coverage to `internal/...` (and other listed packages). |
   | `-covermode=atomic` | Safe counts under concurrency. |
   | `-coverprofile=coverage.out` | Raw profile for `go tool cover`. |
   | `./tests/...` | All test packages only (no tests live next to `internal`). |

7. **Colored log** — Prints test output; lines containing `FAIL` (red), `PASS` (green), `coverage:` (yellow).

8. **`go tool cover -func=coverage.out`** — Prints **per-function** coverage (file, function, %).

---

## Artifacts

| File | Git | Description |
|------|-----|-------------|
| `coverage.out` | **Ignored** (`*.out` in `.gitignore`) | Binary/text profile; regenerate locally or in CI. |

### Optional HTML report

```powershell
go tool cover -html=coverage.out -o coverage.html
```

Open `coverage.html` in a browser for line-by-line highlighting.

---

## Baseline report (instrumented with `run_coverage.ps1`)

> **Regenerate:** after meaningful test or production code changes, run `.\run_coverage.ps1` and update this section (especially **total** and any rows you care about tracking).

| Field | Value |
|--------|--------|
| **Command** | `.\run_coverage.ps1` from repo root |
| **Captured** | 2026-03-27 (re-run after code changes and refresh this doc) |
| **Total statements** | **80.9%** |
| **Output of** | `go tool cover -func=coverage.out` (last line: `total: (statements) 80.9%`) |

### By area (summary)

| Area | Coverage notes |
|------|----------------|
| **Broker use cases** | `StreamClientConn`, `StreamServerConn`, `StreamBiConn`, `NewBroker` — **100%** |
| **SQS use cases** | `Producer`, `Consumer`, `Ack`, `NewSQS` — **100%** |
| **Domain `Message`** | All listed methods — **100%** |
| **gRPC mappers** | All listed methods — **100%** |
| **gRPC interceptors** | `NewInterceptor`, `UnaryAuthInterceptor`, `StreamAuthInterceptor` — **100%** |
| **gRPC middleware** | `ApiToken`, `UnaryLog`, `StreamLog`, `LogErrors`, `NewMiddleware` — **100%**; **`IdPotency` — 0.0%** (not exercised by current tests) |
| **gRPC inbound broker** | `ServerStream`, `BidirectionalStream`, `NewBrokerAPI` — **100%**; **`ClientStream` — 71.4%** |
| **gRPC inbound SQS** | **100%** on listed handlers |
| **gRPC outbound** | **`BidirectionalStream` — 73.5%**, **`ClientStream` — 66.7%**, **`ServerStream` — 55.6%** (retry / error branches partially hit) |
| **Redis `cchecker`** | **`SaveConsumer` / `CheckConsumer` — 100%**; **`NewCheckerClient` — 71.4%** (e.g. TLS branch not hit) |
| **Redis `cmemory`** | **`AckMessage` / `Group` — 100%**; **`GetMessage` — 81.8%**; **`SendMessage` — 77.8%**; **`NewMemoryClient` — 71.4%** |
| **Event** | **`CheckID` / `NewEvent` — 100%**; **`WriteLog` — 77.8%**; **`typeLog` — 66.7%**; **`cmdLog` — 75.0%**; **`clientLog` — 0.0%** (`LOGCLIENT` path) |
| **fasthttp monitor** | **`getHealthCheck` / `collectMetrics` / `NewMonitorAPI` — 100%**; **`getMetrics` / `convertToHTTP` — 75%**; **`WriteHeader` — 0%** |
| **fasthttp middleware** | **API token** middleware — **100%**; **host** `CheckMiddleware` — **60%**; **chain** (`middlewares.go` `NewMiddleware` / `AddMiddleware` / `Middlewares`) — **0%** (composition not run in tests) |

### Per-function detail (`go tool cover -func=coverage.out`)

```
dominus-project/internal/application/use_cases/broker/factory.go:18:                    NewBroker               100.0%
dominus-project/internal/application/use_cases/broker/stream_bi_conn.go:10:         StreamBiConn            100.0%
dominus-project/internal/application/use_cases/broker/stream_client_conn.go:10:       StreamClientConn        100.0%
dominus-project/internal/application/use_cases/broker/stream_server_conn.go:10:       StreamServerConn        100.0%
dominus-project/internal/application/use_cases/sqs/ack.go:8:                          Ack                     100.0%
dominus-project/internal/application/use_cases/sqs/consumer.go:9:                     Consumer                100.0%
dominus-project/internal/application/use_cases/sqs/factory.go:20:                     NewSQS                  100.0%
dominus-project/internal/application/use_cases/sqs/producer.go:10:                    Producer                100.0%
dominus-project/internal/domain/entities/sqs_message.go:15:                            NewMessage              100.0%
dominus-project/internal/domain/entities/sqs_message.go:23:                            SetMessage              100.0%
dominus-project/internal/domain/entities/sqs_message.go:26:                            GetMessage              100.0%
dominus-project/internal/domain/entities/sqs_message.go:30:                            SetMessageId            100.0%
dominus-project/internal/domain/entities/sqs_message.go:33:                            GetMessageId            100.0%
dominus-project/internal/domain/entities/sqs_message.go:37:                            GetCreatedAt            100.0%
dominus-project/internal/domain/entities/sqs_message.go:40:                            SetCreateAt             100.0%
dominus-project/internal/infrastructure/event/event.go:26:                            NewEvent                100.0%
dominus-project/internal/infrastructure/event/event.go:39:                            WriteLog                77.8%
dominus-project/internal/infrastructure/event/event.go:60:                            CheckID                 100.0%
dominus-project/internal/infrastructure/event/event.go:69:                            typeLog                 66.7%
dominus-project/internal/infrastructure/event/event.go:78:                            clientLog               0.0%
dominus-project/internal/infrastructure/event/event.go:81:                            cmdLog                  75.0%
dominus-project/internal/infrastructure/fasthttp/inbound/monitor_api.go:31:           Header                  100.0%
dominus-project/internal/infrastructure/fasthttp/inbound/monitor_api.go:35:           Write                   100.0%
dominus-project/internal/infrastructure/fasthttp/inbound/monitor_api.go:39:           WriteHeader             0.0%
dominus-project/internal/infrastructure/fasthttp/inbound/monitor_api.go:50:           NewMonitorAPI           100.0%
dominus-project/internal/infrastructure/fasthttp/inbound/monitor_api.go:61:           convertToHTTP           75.0%
dominus-project/internal/infrastructure/fasthttp/inbound/monitor_api.go:72:           collectMetrics          100.0%
dominus-project/internal/infrastructure/fasthttp/inbound/monitor_api.go:83:           getMetrics              75.0%
dominus-project/internal/infrastructure/fasthttp/inbound/monitor_api.go:95:           getHealthCheck          100.0%
dominus-project/internal/infrastructure/fasthttp/inbound/monitor_api.go:101:          path                    100.0%
dominus-project/internal/infrastructure/fasthttp/middlewares/api_middleware.go:20:    NewMiddlewareApiToken   100.0%
dominus-project/internal/infrastructure/fasthttp/middlewares/api_middleware.go:27:    CheckMiddleware         100.0%
dominus-project/internal/infrastructure/fasthttp/middlewares/host_allowed.go:17:      NewMiddlewareHost       100.0%
dominus-project/internal/infrastructure/fasthttp/middlewares/host_allowed.go:24:    CheckMiddleware         60.0%
dominus-project/internal/infrastructure/fasthttp/middlewares/middlewares.go:22:     NewMiddleware           0.0%
dominus-project/internal/infrastructure/fasthttp/middlewares/middlewares.go:28:     AddMiddleware           0.0%
dominus-project/internal/infrastructure/fasthttp/middlewares/middlewares.go:32:     Middlewares             0.0%
dominus-project/internal/infrastructure/grpc/inbound/broker_v1.3.7.go:24:           NewBrokerAPI            100.0%
dominus-project/internal/infrastructure/grpc/inbound/broker_v1.3.7.go:32:           ClientStream            71.4%
dominus-project/internal/infrastructure/grpc/inbound/broker_v1.3.7.go:46:           ServerStream            100.0%
dominus-project/internal/infrastructure/grpc/inbound/broker_v1.3.7.go:53:           BidirectionalStream     100.0%
dominus-project/internal/infrastructure/grpc/inbound/sqs_v1.3.7.go:22:              NewSqsAPI               100.0%
dominus-project/internal/infrastructure/grpc/inbound/sqs_v1.3.7.go:28:              Producer                100.0%
dominus-project/internal/infrastructure/grpc/inbound/sqs_v1.3.7.go:43:              Consumer                100.0%
dominus-project/internal/infrastructure/grpc/inbound/sqs_v1.3.7.go:68:              Ack                     100.0%
dominus-project/internal/infrastructure/grpc/mappers/mappers.go:14:                  NewBiStreamConn         100.0%
dominus-project/internal/infrastructure/grpc/mappers/mappers.go:20:                  Recv                    100.0%
dominus-project/internal/infrastructure/grpc/mappers/mappers.go:24:                  Send                    100.0%
dominus-project/internal/infrastructure/grpc/mappers/mappers.go:31:                  Context                 100.0%
dominus-project/internal/infrastructure/grpc/mappers/mappers.go:39:                  NewClientStreamContext  100.0%
dominus-project/internal/infrastructure/grpc/mappers/mappers.go:45:                  Recv                    100.0%
dominus-project/internal/infrastructure/grpc/mappers/mappers.go:49:                  Context                 100.0%
dominus-project/internal/infrastructure/grpc/mappers/mappers.go:57:                  NewServerStreamContext  100.0%
dominus-project/internal/infrastructure/grpc/mappers/mappers.go:63:                  Send                    100.0%
dominus-project/internal/infrastructure/grpc/mappers/mappers.go:70:                  Context                 100.0%
dominus-project/internal/infrastructure/grpc/middlewares/interceptors.go:24:       NewInterceptor          100.0%
dominus-project/internal/infrastructure/grpc/middlewares/interceptors.go:31:       UnaryAuthInterceptor    100.0%
dominus-project/internal/infrastructure/grpc/middlewares/interceptors.go:38:       StreamAuthInterceptor   100.0%
dominus-project/internal/infrastructure/grpc/middlewares/middlewares.go:32:        NewMiddleware           100.0%
dominus-project/internal/infrastructure/grpc/middlewares/middlewares.go:40:        ApiToken                100.0%
dominus-project/internal/infrastructure/grpc/middlewares/middlewares.go:57:        UnaryLog                100.0%
dominus-project/internal/infrastructure/grpc/middlewares/middlewares.go:62:        StreamLog               100.0%
dominus-project/internal/infrastructure/grpc/middlewares/middlewares.go:67:        LogErrors               100.0%
dominus-project/internal/infrastructure/grpc/middlewares/middlewares.go:85:        IdPotency               0.0%
dominus-project/internal/infrastructure/grpc/outbound/client_v1.3.7.go:22:          NewGrpClient            100.0%
dominus-project/internal/infrastructure/grpc/outbound/client_v1.3.7.go:29:          ClientStream            66.7%
dominus-project/internal/infrastructure/grpc/outbound/client_v1.3.7.go:79:          ServerStream            55.6%
dominus-project/internal/infrastructure/grpc/outbound/client_v1.3.7.go:124:         BidirectionalStream     73.5%
dominus-project/internal/infrastructure/redis/cchecker/outbound.go:19:             NewCheckerClient        71.4%
dominus-project/internal/infrastructure/redis/cchecker/outbound.go:64:             SaveConsumer            100.0%
dominus-project/internal/infrastructure/redis/cchecker/outbound.go:74:             CheckConsumer           100.0%
dominus-project/internal/infrastructure/redis/cmemory/outbound.go:23:              NewMemoryClient         71.4%
dominus-project/internal/infrastructure/redis/cmemory/outbound.go:70:              SendMessage             77.8%
dominus-project/internal/infrastructure/redis/cmemory/outbound.go:92:              AckMessage              100.0%
dominus-project/internal/infrastructure/redis/cmemory/outbound.go:101:             GetMessage              81.8%
dominus-project/internal/infrastructure/redis/cmemory/outbound.go:128:             Group                   100.0%
total:                                                                             (statements)            80.9%
```

### Packages not listed above

`go tool cover -func` only prints functions that appear in the profile. Code that is **never reached** by any test may be absent or folded into totals. In particular, **`cmd/`**, **`config/`**, and **`internal/boostrap/`** are not in the table above—those paths are not currently exercised by `./tests/...` with this profile. Raising **total** coverage usually requires integration tests that start the app or targeted tests for config parsing and bootstrap.

---

## CI suggestion

In a pipeline, run the same logic as the script (PowerShell on Windows, or translate the `go list` + filter + `go test` steps to bash). Archive `coverage.out` and fail the job if `total` drops below a threshold you set.
