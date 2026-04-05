# Terraform (local Docker stack)

This repository ships a **Terraform** configuration under [`terraform/`](../terraform/) that provisions a **development-oriented** stack on **Docker** using the [**kreuzwerker/docker**](https://registry.terraform.io/providers/kreuzwerker/docker/latest/docs) provider. It builds the Dominus image from the repo [`Dockerfile`](../Dockerfile), wires a user-defined bridge network, and starts Redis, Prometheus, Grafana, an Nginx sidecar, and the broker container.

## Prerequisites

- [Terraform](https://www.terraform.io/) (CLI)
- **Docker Engine** running locally (the provider talks to the Docker API)
- On **Windows**, the default provider host in [`terraform/terraform.tf`](../terraform/terraform.tf) is the named pipe `npipe:////.//pipe//docker_engine`. On **Linux** or **macOS**, set `provider "docker" { host = "unix:///var/run/docker.sock" }` (or your socket path) before `terraform apply`.

## Layout

| Path | Role |
|------|------|
| [`terraform/main.tf`](../terraform/main.tf) | Root module: network + stacks for Dominus, Redis, sidecar, Prometheus, Grafana |
| [`terraform/modules/network`](../terraform/modules/network) | Docker network (name/driver from variables) |
| [`terraform/modules/image`](../terraform/modules/image) | `docker_image` — build from context or pull by name |
| [`terraform/modules/container`](../terraform/modules/container) | `docker_container` — ports, health checks, volumes, bind mounts |
| [`terraform/modules/volume`](../terraform/modules/volume) | Named volumes for data/certs |
| [`terraform/dev/dominus`](../terraform/dev/dominus) | Dominus broker image build + container |
| [`terraform/dev/redis`](../terraform/dev/redis) | Redis with config bind mount + data volume |
| [`terraform/dev/sidecar`](../terraform/dev/sidecar) | Nginx (`nginx:latest`) with [`nginx.conf`](../terraform/dev/sidecar/nginx.conf) |
| [`terraform/dev/prometheus`](../terraform/dev/prometheus) | Prometheus with [`prometheus.yml`](../terraform/dev/prometheus/prometheus.yml) |
| [`terraform/dev/grafana`](../terraform/dev/grafana) | Grafana with persistent volume |

State uses the **local** backend ([`terraform/terraform.tf`](../terraform/terraform.tf)): `terraform/terraform.tfstate` (see root `.gitignore` for `*.tfstate`; do not commit secrets).

## Variables and defaults

[`terraform/vars.auto.tfvars`](../terraform/vars.auto.tfvars) sets:

- **Network**: `network_name = "dominus"`, `network_mode = "bridge"`
- **Dominus**: build context path `dominus_server_file` (default points at repo root `Dockerfile`), gRPC/REST port mapping **5000** and **8000**
- **Redis**: **6379**, **Prometheus**: **9090**, **Grafana**: **3000**, **Sidecar (Nginx)**: **80**

Override by editing `vars.auto.tfvars` or passing `-var` / a separate `*.tfvars` file.

## Apply and destroy

From the repository root:

```bash
cd terraform
terraform init
terraform plan
terraform apply
```

Tear down:

```bash
terraform destroy
```

Outputs (see [`terraform/outputs.tf`](../terraform/outputs.tf)) expose container names, image IDs, and port mappings for each service.

## `Makefile.ps1` (Windows)

From the **repository root**, you can run Terraform in [`terraform/`](../terraform/) without `cd` by using **`Makefile.ps1`** targets. The script resolves **`terraform`** on **`PATH`**, changes into **`-TerraformDir`** (default: `terraform`), and forwards **`terraform <subcommand>`** with the right working directory.

| Target | Terraform command | Notes |
|--------|-------------------|--------|
| **`terraform-init`** | `terraform init` | Append **`-TerraformExtraArgs`** (e.g. **`-upgrade`**) |
| **`terraform-fmt`** | `terraform fmt -recursive` | Optional **`-TerraformExtraArgs`** (e.g. **`-check`**) |
| **`terraform-validate`** | `terraform validate` | Run **`terraform-init`** first |
| **`terraform-plan`** | `terraform plan` | Optional **`-TerraformExtraArgs`** |
| **`terraform-apply`** | `terraform apply` | **`-TerraformAutoApprove`** for **`-auto-approve`**; optional **`-TerraformExtraArgs`** |
| **`terraform-destroy`** | `terraform destroy` | **`-TerraformAutoApprove`** for **`-auto-approve`**; optional **`-TerraformExtraArgs`** |
| **`terraform-output`** | `terraform output` | Optional **`-TerraformExtraArgs`** (e.g. **`-json`**, or an output name) |

**`-TerraformDir`**: alternate root module path relative to the repo (default **`terraform`**).

Examples (PowerShell):

```powershell
.\Makefile.ps1 -Target terraform-init
.\Makefile.ps1 -Target terraform-plan
.\Makefile.ps1 -Target terraform-apply -TerraformAutoApprove
.\Makefile.ps1 -Target terraform-output -TerraformExtraArgs '-json'
.\Makefile.ps1 -Target terraform-init -TerraformExtraArgs '-upgrade'
```

## Services (high level)

| Component | Notes |
|-----------|--------|
| **Dominus** | Image built from `dominus_server_file`; named volume for certs under `/etc/dominus/certs`; HTTP health check hits **`/health`** with `x-api-key` (see [`doc/http-monitor.md`](http-monitor.md)) |
| **Redis** | Config from `terraform/dev/redis/redis.conf`; data on a named volume |
| **Sidecar** | Nginx reverse proxy; config bind-mounted from the module directory |
| **Prometheus / Grafana** | Standard images; Prometheus scrapes per `prometheus.yml`; Grafana data on a volume |

Dependency order in root [`main.tf`](../terraform/main.tf) ensures Redis before Dominus, Dominus before sidecar/Prometheus, and those before Grafana.

## Configuration vs `go run`

For **`go run ./cmd/api`** locally, Viper loads `env/env.dev.json` relative to the process (see [`doc/configuration.md`](configuration.md)). The Terraform Dominus module does **not** bind-mount the `env/` folder; use **`APP_CONFIG`** (production-style JSON) inside the container if you need injected config, or extend the module with `container_bind_mounts` / `container_env` in [`terraform/modules/container`](../terraform/modules/container) as needed.

## Related docs

- [`doc/configuration.md`](configuration.md) — JSON config shape and bootstrap
- [`doc/http-monitor.md`](http-monitor.md) — `/health` and `x-api-key` (matches container health check)
