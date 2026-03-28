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

## Related tests

- Unit-style: `tests/cases/application/use_cases/broker_test/`
- gRPC wiring: `tests/cases/infrastructure/grpc/inbound_test/broker_test.go`
- Integration: `tests/integration/broker_flow_test/` (real broker + outbound + TCP peers or bufconn)
