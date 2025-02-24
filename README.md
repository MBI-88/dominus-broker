# Dominus Project

``
==========================================================	
██████    ██████  ███    ██  ██ ███    ██ ██    ██ ███████
██   ██  ██    ██ ████  ████ ██ ████   ██ ██    ██ ██
██   ██  ██    ██ ██ ████ ██ ██ ██ ██  ██ ██    ██ ███████
██   ██  ██    ██ ██  ██  ██ ██ ██  ██ ██ ██    ██      ██
██████    ██████  ██      ██ ██ ██   ████  ██████  ███████
==========================================================
👉 Github: <https://github.com/MBI-88>

```

## Overview

The Dominus Project is a containerized application that uses Docker Compose to manage its services. It consists of two main services: `dominus` and `sidecar`. The `dominus` service is the core application, while the `sidecar` service uses Nginx as a reverse proxy.

## Services

### Dominus

The `dominus` service is built from a Dockerfile located at `./Dockerfile.dominus`. It is configured to restart on failure with a delay of 2 seconds and a maximum of 3 attempts. The service exposes ports 8000 and 5000 and includes a health check that verifies the service's health by making an HTTP request to the `/health` endpoint.

#### Configuration

- **Build Context**: `./`
- **Dockerfile**: `./Dockerfile.dominus`
- **Container Name**: `dominus`
- **Replicas**: 1
- **Ports**:
  - `8000:8000`
  - `5000:5000`
- **Volumes**:
  - `dominus_certs:/etc/dominus/certs:ro`
  - `dominus_logs:/var/dominus/logs:rw`
- **Health Check**:
  - **Test**: `wget -qO- --header "API_TOKEN: dominus_example_@10102024KeyServerToken" http://localhost:8000/health`
  - **Interval**: 1 minute
  - **Timeout**: 2 seconds
  - **Retries**: 3

### Sidecar

The `sidecar` service uses the latest Nginx image and acts as a reverse proxy for the `dominus` service. It is configured to restart on failure with a delay of 1 second and a maximum of 3 attempts. The service depends on the `dominus` service and uses a custom Nginx configuration file located at `./nginx.conf`.

#### Infrastructure

- **Image**: `nginx:latest`
- **Container Name**: `sidecar`
- **Replicas**: 1
- **Volumes**:
  - `./nginx.conf:/etc/nginx/nginx.conf:ro`
- **Environment Variables**:
  - `NGINX_HOST=0.0.0.0`
  - `NGINX_PORT=80`
  - `NGINX_ENVSUBST_TEMPLATE_SUFFIX=.conf`
- **Depends On**: `dominus`

## Networks

Both services are connected to the `dominus` network.

## Volumes

- `dominus_certs`: Stores the certificates for the `dominus` service.
- `dominus_logs`: Stores the logs for the `dominus` service.

## Usage

To start the Dominus Project, use the following command: go run . start  

## Test Coverage

The following is the test coverage for the Dominus project:

- **dominus/tests/init.go**: 100.0%
- **dominus/tests/mock_conn.go**:
  - `helperErr`: 100.0%
  - `Descriptor`: 0.0%
  - `GetPayload`: 100.0%
  - `GetSubscribers`: 100.0%
  - `Reset`: 0.0%
  - `String`: 0.0%
  - `Validate`: 0.0%
  - `ValidateAll`: 0.0%
  - `NewSimpleConnMock`: 100.0%
  - `Recv`: 100.0%
  - `NewClienConnMock`: 100.0%
  - `Send`: 100.0%
  - `Context`: 0.0%
  - `NewServerConnMock`: 100.0%
  - `Recv`: 100.0%
  - `Send`: 100.0%
  - `NewBiConnMock`: 100.0%
- **dominus/tests/mock_event.go**:
  - `WriteLog`: 0.0%
  - `Printf`: 0.0%
  - `NewEventMock`: 100.0%
- **dominus/tests/mock_grpc_client.go**:
  - `Simple`: 100.0%
  - `ClientStream`: 100.0%
  - `ServerStream`: 84.2%
  - `BidirectionalStream`: 88.0%
  - `NewGrpcClientMock`: 100.0%
- **dominus/tests/mock_repo.go**:
  - `Migrations`: 0.0%
  - `InsertObject`: 100.0%
  - `DeleteObjects`: 0.0%
  - `CountPages`: 0.0%
  - `Filter`: 0.0%
  - `Backup`: 0.0%
  - `Stats`: 0.0%
  - `NewRepoMock`: 100.0%

**Total Test Coverage**: 81.8%  
