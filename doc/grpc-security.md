# gRPC middleware and security

Server and client interceptors are wired in `internal/boostrap/boostrap.go`. Middleware types live in `internal/infrastructure/grpc/middlewares/`.

## API token (streams and unary)

- **Server**: `auth.StreamServerInterceptor(midGs.ApiToken)` and `auth.UnaryServerInterceptor(midGs.ApiToken)` with `Middleware.ApiToken` (`middlewares.go`).
- **Client** (outbound): `Interceptor.StreamAuthInterceptor` / `UnaryAuthInterceptor` attach the configured secret as metadata (see `enum` for header name, e.g. `X_API_KEY`).
- `ApiToken` hashes the incoming token with SHA-256 and compares it to the hash of `connection_key` from config (`GrpcConfig.ConnectionKey`).

`Event.CheckID` runs inside `ApiToken` to ensure a correlation id exists in context.

## Idempotency (unary only)

- **Server**: `auth.UnaryServerInterceptor(midGs.IdemPotency)` — **not** applied to streaming RPCs in the current bootstrap chain.
- Uses header from `enum` (e.g. idempotency key) and `repositories.CheckerClient` (`cchecker`) to reject duplicates or persist keys.

### Concurrency semantics (important for operators and future changes)

The middleware flow is roughly: `CheckConsumer` (Redis `EXISTS` on the idempotency key) → if missing, spawn a goroutine that calls `SaveConsumer` (`SET` with `NX` and TTL) → continue to the handler.

**Limitation:** two concurrent unary requests with the **same** idempotency key can both observe “key absent” before either `SET NX` completes, so **both** may reach the handler. This is a classic check-then-act window at the application level; the race detector will not flag it because Redis handles its own concurrency.

**Stronger guarantee (if required):** reserve the key **synchronously** before invoking the handler (for example: single `SET key NX` — proceed only if the set succeeds; if the key already exists, return duplicate/aborted). That collapses check and insert into one atomic step. See also `doc/concurrency.md`.

## TLS

When cert, key, and CA files exist:

- gRPC server uses server credentials from cert/key.
- gRPC client dial uses CA for server verification.
- REST can use `ListenAndServeTLS` (see `restServer` in bootstrap).

Otherwise both stacks use `insecure` credentials / plain HTTP.

## Observability

- **OpenTelemetry**: `otelgrpc` stats handlers on server and client dial options.
- **Prometheus**: `grpc-ecosystem/go-grpc-middleware/providers/prometheus` unary and stream interceptors.

## Related tests

- `tests/cases/infrastructure/grpc/middleware_test/`
- `tests/cases/infrastructure/grpc/middleware_test/interceptors_test.go`
