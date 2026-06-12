# Broker streaming

gRPC service **`BrokerStream`** (protobuf package `dominus`) exposes three streaming RPCs. The server implementation lives in `internal/infrastructure/grpc/inbound/` (gRPC handlers); business logic lives under `internal/application/usecases/` in the `stream_client`, `stream_server`, and `stream_bidirectional` packages.

## Layers

| Layer | Path | Role |
|--------|------|------|
| Transport | `inbound/broker_stream.go` | Implements `BrokerAPIServer`, maps streams to DTOs |
| Mappers | `internal/infrastructure/grpc/mappers/mappers.go` | Wraps gRPC streams as `Broker*Dto` interfaces |
| Use cases | `internal/application/usecases/stream_client`, `.../stream_server`, `.../stream_bidirectional` | `StreamClient`, `StreamServer`, `StreamBidirectional` |
| Outbound client | `internal/infrastructure/grpc/outbound/client_stream.go` | `repositories.BrokerClient` — dials subscriber URLs |

## Client stream (ingress client → this server → outbound to subscribers)

1. Ingress receives a **client-streaming** RPC: many `StreamRequestMessage`, one `StreamResponseMessage` at close.
2. The client-stream use case reads the **first** message and takes `Subscribers` (URLs). Further messages only carry `Payload`.
3. It starts `go client.ClientStream(subscribers, stream, ctx)` and forwards each payload on a `chan []byte`.
4. Outbound opens one **client** stream per subscriber URL and sends the same payload shape (`Subscribers` + `Payload`) to each.

**Inbound success path**: use case returns `io.EOF` when the ingress client ends the stream; handler responds with `codes.OK` via `SendAndClose`.

## Server stream

1. Unary-style request: first `StreamRequestMessage` includes `Subscribers` and initial `Payload`.
2. The server-stream use case calls the outbound client's `ServerStream` per URL, merges responses into a channel (current implementation uses a buffered channel sized `total + int(total*2/3)`), and sends chunks to the ingress **server-stream** response.
3. When the outbound indicates completion the use case returns an error (e.g. a formatted `connection closed` message). The gRPC handler maps that to **`codes.Aborted`** with that message (see inbound code).

## Bidirectional stream

1. First message must include non-empty `Subscribers`.
2. The bidirectional use case runs the outbound `BidirectionalStream`, pumps provider payloads from the ingress stream into an unbuffered `streamProv`, and forwards subscriber payloads back through `stream.Send` using a buffered `streamSub` sized `total + int(total*2/3)`.
3. Completion is driven by per-URL workers and a `done` signal; the use case returns an error when all subscribers finish and the outbound signals shutdown.

**Concurrency note**: the outbound implementation must **not** hold the process-wide mutex across blocking `Recv`/`Send` calls, or fan-out deadlocks against echoing peers. State per URL is tracked in a small `endpoint` struct guarded by short critical sections only.

### Bidirectional shutdown: `done` vs `closed` (`tx`)

Two signals are easy to confuse; they answer different questions.

| Signal | Where | Meaning |
|--------|--------|---------|
| **`done`** | Each per-URL worker in `BidirectionalStream` sends **once** in a `defer` when that worker goroutine **returns** | All workers have exited; none can still be in the `select` that sends to `streamSub` (`subMsg`). |
| **`closed`** (argument `tx` in `BidirectionalStream`) | Main goroutine of the outbound client, **after** `for msg := range provMsg` ends and `workerWG.Wait()` | Inbound provider side has closed `streamProv`; all workers have stopped; safe for the recv loop to finish waiting. |

**Why `close(streamSub)` is safe after counting `done` to zero**

Workers only forward payloads in a `select` between `ctx.Done()` and `subMsg <- payload`. While a worker is blocked on `subMsg <-`, it has **not** yet returned, so its `defer` has **not** run and it has **not** incremented the use-case’s completion count via `done`. The use case (`StreamBiConn`) decrements `total` only on `<-done` and calls `close(streamSub)` when `total == 0`, i.e. after **every** worker goroutine has finished. Therefore no worker should still send on `streamSub` after that close.

**Provider-side wait on `closed`**

When the ingress `Recv` loop hits an error, it `close(streamProv)` then blocks on `<-closed`. That unblocks only after the outbound client drains `provMsg`, waits for all workers (`workerWG.Wait()`), and writes to `tx` (`closed`). That keeps the provider half and the outbound fan-out torn down in order.

**Context cancellation**

`StreamBidirectional` calls `cancel()` when `stream.Send` to the client fails (and `defer cancel()` on exit). Workers observe `ctx` in their send `select` and in retry loops, so cancellation can unblock a blocked send path. Normal subscriber completion is usually driven by `Recv` returning `EOF` rather than by an explicit cancel before `close(streamSub)`.

### Channel lifecycle in `StreamBidirectional` (`stream_bidirectional_service.go`)

These channels coordinate the provider goroutine, the main loop, the error logger, and `BidirectionalStream` in `outbound/client_stream.go`.

| Channel | Capacity | Who closes | Role |
|---------|----------|------------|------|
| **`streamProv`** | unbuffered | Provider goroutine on first ingress `Recv` error | Feeds payloads into outbound `for msg := range provMsg`; closing it ends that loop so the outbound can `workerWG.Wait()` and signal `tx` (`closed`). |
| **`streamSub`** | buffered | Main loop when `total == 0` after all `<-done` | Merges subscriber responses toward `stream.Send`; safe to close only after every worker has returned (see `done` above). |
| **`errMsg`** | buffered | Main loop together with `streamSub` | Consumed by a small goroutine `for range errMsg`; must be closed so that goroutine exits. |
| **`done`** | `len(subscribers)` | **`defer close(done)`** after allocation | One send per worker when the worker goroutine exits; main counts down to zero, then closes `streamSub` / `errMsg`. |
| **`closed`** | unbuffered | **Provider goroutine only**, after `<-closed` | Handed to outbound as `tx`. Outbound sends **once** after `provMsg` is drained and workers have joined. The provider waits on `<-closed` so it does not exit before that handshake; then it **`close(closed)`**. |

**Do not `defer close(closed)` at the start of `StreamBidirectional`.** The main loop can return (`"connection closed"`) while the outbound goroutine has not yet executed `tx <- struct{}{}`. A top-level `defer close(closed)` would run on return and race the send → **panic: send on closed channel**. The `closed` channel is therefore closed **only** in the provider path, **after** the outbound’s single send has been received.

**`defer close(done)`** is safe here: the main loop only returns after receiving exactly `len(subscribers)` values from `done`, so no worker should still send on `done` when the defer runs. Early returns before `BidirectionalStream` is started never write to `done`; closing an empty buffered channel is fine.

**Implementation notes:** current use-case implementations in `internal/application/usecases/stream_*` follow these concrete patterns observed in code: explicit `len(subscribers)` checks returning errors when empty; `context.WithCancel` used to cancel provider-side loops when `stream.Send` fails; `stream`/`streamSub` buffers sized relative to `total` to reduce blocking; and a `closed` handshake channel used so the provider waits for the outbound to drain and acknowledge shutdown before exiting.

### Payload fan-out and shared `[]byte`

For each `msg` from `provMsg`, outbound `BidirectionalStream` launches `go cls(msg)` per subscriber URL. All goroutines for that iteration observe the **same** slice header (same backing array). Treat payloads as **read-only** until all sends for that message complete; if the producer reuses or mutates the buffer before gRPC finishes reading, you get subtle races or corrupted payloads.

## Related tests

- Unit-style: `tests/cases/application/use_cases/broker_test/`
- gRPC wiring: `tests/cases/infrastructure/grpc/inbound_test/broker_test.go` (handlers, `Broker` interface mocking, context-driven flows)
- Integration: `tests/integration/broker_flow_test/` (real broker + outbound + TCP peers or bufconn). For **BidirectionalStream** shutdown expectations when subscribers keep streams open, see **`doc/concurrency.md`** (*Integration tests: CloseSend vs context cancel*).
