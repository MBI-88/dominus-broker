# Configuration

Runtime settings are loaded through `config.NewConfig(prod bool)` in `config/config.go`, using **Viper** and JSON.

## Modes

| `prod` | Source |
|--------|--------|
| `false` (default CLI) | File `./../env/env.dev.json` relative to the process working directory (typically run from `cmd/api`, so `env/env.dev.json` at repo root). |
| `true` (`-prod` flag) | Environment variable **`APP_CONFIG`**: full JSON string unmarshalled into `Config`. `viper.AutomaticEnv()` is enabled. |

## Root JSON shape (`Config`)

Struct tags define the expected keys:

- `grpc_config` → `GrpcConfig`: `grpc_port`, `connection_key`
- `rest_config` → `RestConfig`: `rest_port`, `api_token`, `allow_origins`
- `cert_config` → `CertConfig`: paths for TLS materials
- `redis_config` → `RedisConfig`: Redis client, DB indices, stream/group names, idempotency TTL (`id_potency`), timeouts, TLS flag
- `infra_config` → `InfraConfig`: `log_mode`, `log_url`

**Important**: `redis_config` must sit at the **root** of the JSON object, matching `json:"redis_config"` on `Config`. Nesting it (e.g. under `provider_config`) will not populate `cf.RedisConfig` unless you change the struct or file layout.

## Bootstrap usage

`internal/boostrap/boostrap.go` reads `cf` to build:

- gRPC server options (TLS vs insecure, interceptors, Prometheus, logging)
- Redis-backed `cchecker` and `cmemory` clients
- REST server and routes

If certificate files are missing (`Stat` on cert paths), gRPC and REST fall back to **insecure** / plain HTTP for local development.

## Related files

- `env/env.dev.json`, `env/env.prod.json` — samples / deployment templates
- **[doc/terraform.md](terraform.md)** — local Docker stack via Terraform (certs volume; use `APP_CONFIG` or extra mounts for in-container JSON)
