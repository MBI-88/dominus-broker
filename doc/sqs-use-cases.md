# SQS-like gRPC use cases

The **`SqsAPI`** service (protobuf package `dominus`) exposes **Producer**, **Consumer**, and **Ack** operations. This document describes layering, validation, Redis behavior, and how tests map to that behavior.

| Layer | Location |
|-------|----------|
| gRPC handlers | `internal/infrastructure/grpc/inbound/sqs_v1.3.7.go` |
| Application | `internal/application/use_cases/sqs/` |
| Domain entity | `internal/domain/entities/sqs_message.go` (`Message`) |
| Redis adapter | `internal/infrastructure/redis/cmemory/outbound.go` (`MemoryClient`) |

---

## Ports and DTOs

- **`repositories.MemoryClient`** — Redis Streams adapter (`cmemory`).
- **DTOs** — `internal/application/dtos/sqs_dto.go` (`ProducerDto`, `ConsumerDto`) for use-case inputs; protobuf requests implement these interfaces at the edge.

Factory: **`sqs.NewSQS(client)`** in bootstrap; the same Redis client is configured with stream and consumer group from config (`Group` is created at startup).

---

## Operations

| RPC | Use case | Redis-backed behavior (via `MemoryClient`) |
|-----|----------|---------------------------------------------|
| **Producer** | Validates payload, builds `entities.Message`, sends | `SendMessage` → `XADD` with JSON payload under `enum.PAYLOAD`; optional explicit stream ID from the message |
| **Consumer** | Reads for a worker + group | `GetMessage` → `XREADGROUP` (blocking read, count 1); stream entry ID is applied to the entity after unmarshal |
| **Ack** | Validates message ID format, then acknowledges | `AckMessage` → `XACK` for `(stream, group, messageId)` |

---

## Validation and error mapping

### gRPC inbound (`sqs_v1.3.7.go`)

| Check | RPCs | gRPC code |
|-------|------|-----------|
| Non-empty payload | Producer | `OutOfRange` |
| Non-empty `group_id` | Consumer, Ack | `NotFound` |
| Non-empty `worker_id` | Consumer, Ack | `NotFound` |
| Non-empty `message_id` | Ack | `NotFound` |
| Use-case / Redis failure | All | `Aborted` (message from `err.Error()`) |

### Application use cases

- **Producer** — Rejects **empty** payload (`empty payload`) before calling Redis.
- **Consumer** — Does **not** duplicate group/worker checks; relies on gRPC for empty IDs. Delegates to `GetMessage(workerId, groupId)`.
- **Ack** — After gRPC ensures a non-empty `message_id`, the use case validates that the ID matches a **Redis stream ID shape**: two non-negative integer components separated by a single `-` (e.g. `1700000000001-0`). Invalid values return **`invalid messageId`** and **never** call `AckMessage`.

Domain rules live on **`entities.Message.SetMessageId`** / `checkValidFormatID`; UUIDs and arbitrary strings are rejected.

**Implication for clients:** `Consumer` returns a `message_id` that is suitable for **`Ack`** in normal operation (Redis-assigned ID). Tests and mocks should use the same `<milliseconds>-<sequence>` style when simulating success paths.

---

## Tests

### Unit — application (`tests/cases/application/use_cases/sqs_test/`)

- **Producer** — Asserts delegation to `SendMessage` with a non-empty `*entities.Message` whose **body** matches the DTO payload; empty payload and Redis errors.
- **Consumer** — Asserts `GetMessage` arguments and, on success, returned **body**, **id**, and **created-at** match the mocked message.
- **Ack** — Valid stream-style IDs with `AckMessage` success/failure; **invalid** IDs (e.g. `message-1`, `123456789`) ensure **`AckMessage` is not invoked** and error text is `invalid messageId`.

### Unit — gRPC inbound (`tests/cases/infrastructure/grpc/inbound_test/sqs_test.go`)

- Handlers with **mocked** `sqs.SQS` for wiring, empty-field validation, and error propagation.
- **`Ack` with invalid `message_id`** uses the **real** `sqs.NewSQS(MockMemoryClient)` so application validation runs end-to-end through the handler (expects `Aborted` / `invalid messageId`).
- Happy-path **Consumer** / **Ack** mocks use **stream-shaped** message IDs so they align with production contracts.

### Redis adapter (`tests/cases/infrastructure/redis/cmemory_test/`)

- `SendMessage`, `GetMessage`, `AckMessage`, and `Group` against **miniredis** with realistic stream IDs.

### Integration (`tests/integration/sqs_flow_test/`)

- **bufconn** gRPC + **miniredis** + real **`cmemory`** and **`sqs.NewSQS`**: full Producer / Consumer / Ack flows, empty payload, Redis failures, empty IDs, multiple messages per group, and post-Ack pending checks (`XPendingExt`).

### Domain (`tests/cases/domain/entities_test/message_test.go`)

- `Message` fields and **stream ID** acceptance/rejection (including UUID rejection).

---

## Related documentation

- Redis streams adapter details: **[redis.md](redis.md)** (`cmemory` section).
- Overall test layout and commands: **[testing.md](testing.md)**.
