# Redis: checker and queue memory

Two adapters use [go-redis](https://github.com/redis/go-redis/v9) against the same Redis cluster but typically **different logical DB indices** from `RedisConfig.CheckerDB` and `RedisConfig.MemoryDB`.

## `cchecker` — idempotency (`internal/infrastructure/redis/cchecker`)

- Implements `repositories.CheckerClient`.
- Keys are built with prefix `enum.ID_POTENCY_TAG` and the caller-provided key.
- `SaveConsumer`: `SET` with **NX** and TTL from `IdPotency` (seconds) — value is a simple string payload suitable for marshalling by go-redis.
- `CheckConsumer`: `EXISTS` on the same key.
- Constructed in bootstrap as `cchecker.NewCheckerClient(...)` with pool/timeouts from config.

**Tests**: `tests/cases/infrastructure/redis/cchecker_test/` using [miniredis](https://github.com/alicebob/miniredis).

## `cmemory` — Sqs-like streams (`internal/infrastructure/redis/cmemory`)

- Implements `repositories.MemoryClient` for producer/consumer/ack.
- Uses **Redis Streams** (`XADD`, `XREADGROUP`, `XACK`, `XGROUP CREATE`): stream name from config `StreamID`, consumer group `GroupID`.
- Payload field name from `enum.PAYLOAD`; messages are JSON (`jsoniter`) of `entities.Message`.
- `GetMessage` must unmarshal into a concrete value (e.g. `var q entities.Message` then `&q`); a nil pointer breaks `Unmarshal`.
- **Ack** at the application layer validates `message_id` as a Redis stream–style ID before `XACK`; see **[sqs-use-cases.md](sqs-use-cases.md)** (validation and gRPC mapping).

**Tests**: `tests/cases/infrastructure/redis/cmemory_test/` (miniredis + `MockEvent` for async `WriteLog`). Higher-level Sqs gRPC + Redis flows: `tests/integration/sqs_flow_test/`.

## Configuration fields (excerpt)

| Field | Typical use |
|-------|-------------|
| `checker_db` | DB index for idempotency keys |
| `memory_db` | DB index for stream-backed queue |
| `stream_id` | Stream key for `XADD` / reads |
| `group_id` | Consumer group name |
| `id_potency` | TTL seconds for checker keys |

Bootstrap calls `qclient.Group(cf.RedisConfig.GroupID)` at startup to ensure the group exists.
