# Dominus — technical documentation

In-depth notes for contributors. The [main README](../README.md) stays high-level; these pages focus on specific subsystems.

| Document | Topics |
|----------|--------|
| [broker-streaming.md](broker-streaming.md) | `BrokerAPI` flows, use cases, outbound fan-out |
| [configuration.md](configuration.md) | Viper, JSON shape, dev vs prod |
| [grpc-security.md](grpc-security.md) | Middleware, API key, idempotency, TLS |
| [redis.md](redis.md) | `cchecker`, `cmemory`, streams |
| [sqs-use-cases.md](sqs-use-cases.md) | Producer, Consumer, Ack: layers, validation, message IDs, test map |
| [http-monitor.md](http-monitor.md) | fasthttp monitor API and middleware |
| [testing.md](testing.md) | `tests/cases` vs `tests/integration`, commands |
| [coverage.md](coverage.md) | **`run_coverage.ps1`** (canonical), `-coverpkg`, baseline **80.9%**, per-function table |
