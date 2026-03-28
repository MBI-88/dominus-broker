# gRPC middleware and security

Server and client interceptors are wired in `internal/boostrap/boostrap.go`. Middleware types live in `internal/infrastructure/grpc/middlewares/`.

## API token (streams and unary)

- **Server**: `auth.StreamServerInterceptor(midGs.ApiToken)` and `auth.UnaryServerInterceptor(midGs.ApiToken)` with `Middleware.ApiToken` (`middlewares.go`).
- **Client** (outbound): `Interceptor.StreamAuthInterceptor` / `UnaryAuthInterceptor` attach the configured secret as metadata (see `enum` for header name, e.g. `X_API_KEY`).
- `ApiToken` hashes the incoming token with SHA-256 and compares it to the hash of `connection_key` from config (`GrpcConfig.ConnectionKey`).

`Event.CheckID` runs inside `ApiToken` to ensure a correlation id exists in context.

## Idempotency (unary only)

- **Server**: `auth.UnaryServerInterceptor(midGs.IdPotency)` — **not** applied to streaming RPCs in the current bootstrap chain.
- Uses header from `enum` (e.g. idempotency key) and `repositories.CheckerClient` (`cchecker`) to reject duplicates or persist keys.

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
