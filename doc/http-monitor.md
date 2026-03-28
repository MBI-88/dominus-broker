# HTTP monitor API (fasthttp)

REST surface is built with **valyala/fasthttp** and **fasthttp/router** in `internal/boostrap/boostrap.go` (`restServer`).

## Entry

- `internal/infrastructure/fasthttp/inbound/monitor_api.go` registers routes (e.g. health / metrics exposure depending on setup).
- Middleware stack: `internal/infrastructure/fasthttp/middlewares/` — API token validation and allowed origins, composed in `restServer`.

## Configuration

- `RestConfig.RestPort`, `ApiToken`, `AllowOrigins` from JSON.
- TLS: if certificates are available, REST listens with TLS; otherwise HTTP.

## Health checks

Docker healthcheck (`docker/dominus/docker-compose.yml`) calls HTTP **`/health`** with header **`x-api-key`** set to the same token as in config (see `rest_config.api_token` in `env` samples).

## Tests

- `tests/cases/infrastructure/fasthttp/monitor_test/`
- `tests/cases/infrastructure/fasthttp/middleware_test/`
