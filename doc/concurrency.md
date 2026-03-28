# Concurrency, shutdown, and data races

Notes for contributors maintaining goroutines, channels, and gRPC streams.

## Race detector

CI-style coverage uses `go test -race` (see `doc/coverage.md` and `run_coverage.ps1`). A green run means no **memory** data races were observed under the executed tests; it does not prove absence of **logical** races (e.g. idempotency).

## Logging

`Event.WriteLog` is invoked from many goroutines (often via `go ...` in middleware and Redis adapters). The implementation uses `log/slog`, whose methods are safe for concurrent use.

## Broker bidirectional stream

See **`doc/broker-streaming.md`** (section *Bidirectional shutdown: `done` vs `closed`*):

- **`done`**: per-worker, signals that the worker goroutine has **returned** — this is the barrier that makes `close(streamSub)` safe in `StreamBiConn`.
- **`closed` / `tx`**: signals that the outbound client finished draining `streamProv` and waited for workers — used so the ingress recv side can wait after closing the provider channel.

Do not replace these with a single channel without re-deriving the same ordering guarantees.

## Outbound `BrokerClient` (`client_v1.3.7.go`)

- **ClientStream**: one mutex per subscriber URL serializes access to the stream handle and connection error state; each ingress message fans out with `go cls(m)` (same `[]byte` must remain read-only for all callees).
- **BidirectionalStream**: `Recv` runs in worker goroutines; `Send` runs from closures launched per message. Stream pointers are updated under a mutex; do not hold that mutex across blocking gRPC calls.

## Unary idempotency middleware

See **`doc/grpc-security.md`** (*Concurrency semantics*). The current `EXISTS` then async `SET NX` pattern does not fully deduplicate concurrent duplicate keys; fix at the middleware/Redis boundary if product requirements demand strict single-flight per key.
