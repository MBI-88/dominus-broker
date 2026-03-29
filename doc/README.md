# Dominus — technical documentation

In-depth notes for contributors. The [main README](../README.md) stays high-level; these pages focus on specific subsystems.

| Document | Topics |
|----------|--------|
| [broker-streaming.md](broker-streaming.md) | `BrokerAPI` flows, outbound fan-out, bidi shutdown (`done` / `closed`), channel lifecycle, shared payload bytes |
| [concurrency.md](concurrency.md) | Race detector, logging, bidi integration teardown (`CloseSend` vs cancel), idempotency |
| [configuration.md](configuration.md) | Viper, JSON shape, dev vs prod |
| [grpc-security.md](grpc-security.md) | Middleware, API key, idempotency (incl. concurrent duplicate keys), TLS |
| [redis.md](redis.md) | `cchecker`, `cmemory`, streams |
| [sqs-use-cases.md](sqs-use-cases.md) | Producer, Consumer, Ack: layers, validation, message IDs, test map |
| [http-monitor.md](http-monitor.md) | fasthttp monitor API and middleware |
| [testing.md](testing.md) | `tests/cases` vs `tests/integration`, commands |
| [tradeoffs.md](tradeoffs.md) | Architectural choices, limits, Redis/gRPC/broker trade-offs (English) |
| [coverage.md](coverage.md) | **`run_coverage.ps1`** (canonical), `-coverpkg`, baseline **90.0%**, per-function table |
