# Broker streaming

gRPC service **`BrokerAPI`** (protobuf package `dominus`) exposes three streaming RPCs. The server implementation lives in `internal/infrastructure/grpc/inbound/broker_v1.3.7.go`; business logic in `internal/application/use_cases/broker/`.

## Layers

| Layer | Path | Role |
|--------|------|------|
| Transport | `inbound/broker_v1.3.7.go` | Implements `BrokerAPIServer`, maps streams to DTOs |
| Mappers | `internal/infrastructure/grpc/mappers/mappers.go` | Wraps gRPC streams as `dtos.Broker*` interfaces |
| Use cases | `internal/application/use_cases/broker/*.go` | `StreamClientConn`, `StreamServerConn`, `StreamBiConn` |
| Outbound client | `internal/infrastructure/grpc/outbound/client_v1.3.7.go` | `repositories.BrokerClient` — dials subscriber URLs |

`broker.NewBroker(client)` receives the outbound `BrokerClient` created in `internal/boostrap/boostrap.go` with the same dial options as production (TLS or insecure, metrics, auth interceptor, compression).

## Client stream (ingress client → this server → outbound to subscribers)

1. Ingress receives a **client-streaming** RPC: many `StreamRequestMessage`, one `StreamResponseMessage` at close.
2. `StreamClientConn` (`stream_client_conn.go`) reads the **first** message and takes `Subscribers` (URLs). Further messages only carry `Payload`.
3. It starts `go client.ClientStream(subscribers, stream, ctx)` and forwards each payload on a `chan []byte`.
4. Outbound opens one **client** stream per subscriber URL and sends the same payload shape (`Subscribers` + `Payload`) to each.

**Inbound success path**: use case returns `io.EOF` when the ingress client ends the stream; handler responds with `codes.OK` via `SendAndClose`.

## Server stream

1. Unary-style request: first `StreamRequestMessage` includes `Subscribers` and initial `Payload`.
2. `StreamServerConn` (`stream_server_conn.go`) calls `ServerStream` on the outbound client for each URL, merges responses into one channel, and sends chunks to the ingress **server-stream** response.
3. When all subscriber streams signal completion, the use case returns `fmt.Errorf("connection closed")`. The gRPC handler maps that to **`codes.Aborted`** with that message (see inbound code).

## Bidirectional stream

1. First message must include non-empty `Subscribers`.
2. `StreamBiConn` (`stream_bi_conn.go`) runs outbound `BidirectionalStream`, pumps provider payloads from the ingress stream into `streamProv`, and forwards subscriber payloads back through `stream.Send`.
3. Completion is driven by per-URL workers in `outbound/client_v1.3.7.go` and the shared `done` channel; the use case returns `"connection closed"` when all subscribers finish.

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

`StreamBiConn` calls `cancel()` when `stream.Send` to the client fails (and `defer cancel()` on exit). Workers observe `ctx` in their send `select` and in retry loops, so cancellation can unblock a blocked send path. Normal subscriber completion is usually driven by `Recv` returning `EOF` rather than by an explicit cancel before `close(streamSub)`.

## Related tests

- Unit-style: `tests/cases/application/use_cases/broker_test/`
- gRPC wiring: `tests/cases/infrastructure/grpc/inbound_test/broker_test.go`
- Integration: `tests/integration/broker_flow_test/` (real broker + outbound + TCP peers or bufconn)
