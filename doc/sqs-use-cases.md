# SQS-like gRPC use cases

The **`SqsAPI`** service (protobuf `dominus`) exposes producer, consumer, and ack operations. Implementation: `internal/infrastructure/grpc/inbound/sqs_v1.3.7.go`. Application layer: `internal/application/use_cases/sqs/`.

## Ports

- `repositories.MemoryClient` — Redis stream adapter (`cmemory`).
- DTOs under `internal/application/dtos/` for request/response mapping.

## Operations

| RPC | Use case | Redis-backed behavior (via `MemoryClient`) |
|-----|----------|---------------------------------------------|
| Producer | Sends a message into the stream | `SendMessage` → `XADD` with payload |
| Consumer | Pulls for a worker + group | `GetMessage` → `XREADGROUP` |
| Ack | Acknowledges processing | `AckMessage` → `XACK` |

Factory: `sqs.NewSQS(client)` in bootstrap; same Redis client is configured with stream and group from config.

## Tests

- `tests/cases/application/use_cases/sqs_test/` — mocks for DTOs and `MemoryClient`
- `tests/cases/infrastructure/grpc/inbound_test/sqs_test.go` — inbound wiring
- `tests/integration/sqs_flow_test/` — reserved for higher-level flows (extend as needed)
