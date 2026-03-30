# Concurrency, shutdown, and data races

Notes for contributors maintaining goroutines, channels, and gRPC streams.

## Race detector

CI-style coverage uses `go test -race` (see `doc/coverage.md` and **`Makefile.ps1 -Target test-cover`**). A green run means no **memory** data races were observed under the executed tests; it does not prove absence of **logical** races (e.g. idempotency).

## Logging

`Event.WriteLog` is invoked from many goroutines (often via `go ...` in middleware and Redis adapters). The implementation uses `log/slog`, whose methods are safe for concurrent use.

## Broker bidirectional stream

See **`doc/broker-streaming.md`** (section *Bidirectional shutdown: `done` vs `closed`*):

- **`done`**: per-worker, signals that the worker goroutine has **returned** — this is the barrier that makes `close(streamSub)` safe in `StreamBiConn`.
- **`closed` / `tx`**: signals that the outbound client finished draining `streamProv` and waited for workers — used so the ingress recv side can wait after closing the provider channel.

Do not replace these with a single channel without re-deriving the same ordering guarantees.

### Integration tests: `CloseSend` vs context cancel (`broker_flow_test`)

`tests/integration/broker_flow_test/bidirectional_stream_flow_test.go` dials **real TCP** subscriber peers (`bidirEchoPeer` in `helpers_test.go`) that loop forever on `Recv` until the connection breaks. That mirrors long-lived third-party brokers.

**Important:** if the test only calls **`CloseSend()`** on the ingress client and expects a clean shutdown:

1. The server’s provider goroutine gets **EOF**, **`close(streamProv)`**, and the outbound exits `for msg := range provMsg`.
2. Outbound then **`workerWG.Wait()`** for each per-URL worker.
3. Those workers are often still blocked on **`Recv()`** toward the echo peer, which does **not** close the stream just because the ingress half-closed.

So **`Wait()` never completes**, the outbound never sends on **`tx` (`closed`)**, and the provider stays blocked on **`<-closed`** → **deadlock**. Closing channels in the use case does not fix that; the workers must be unblocked.

**Recommended teardown in tests:** after asserting application-level responses, **`cancel()`** the RPC context (the same tree passed into `BidirectionalStream`). Workers observe `ctx.Done()`, exit, emit **`done`**, and the outbound can finish and signal **`closed`**. Then drain the final client **`Recv()`** and accept terminal codes such as **`Canceled`**, **`Aborted`**, **`Unavailable`**, or **`io.EOF`** depending on timing and gRPC version.

Documented channel rules and the `closed` / `done` split: **`doc/broker-streaming.md`** (*Channel lifecycle in `StreamBiConn`*).

## Outbound `BrokerClient` (`client_v1.3.7.go`)

- **ClientStream**: one mutex per subscriber URL serializes access to the stream handle and connection error state; each ingress message fans out with `go cls(m)` (same `[]byte` must remain read-only for all callees).
- **BidirectionalStream**: `Recv` runs in worker goroutines; `Send` runs from closures launched per message. Stream pointers are updated under a mutex; do not hold that mutex across blocking gRPC calls.

## Unary idempotency middleware

See **`doc/grpc-security.md`** (*Concurrency semantics*). The current `EXISTS` then async `SET NX` pattern does not fully deduplicate concurrent duplicate keys; fix at the middleware/Redis boundary if product requirements demand strict single-flight per key.
