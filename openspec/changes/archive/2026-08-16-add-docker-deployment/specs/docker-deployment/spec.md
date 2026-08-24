## ADDED Requirements

### Requirement: Multi-Stage Container Packaging
The project SHALL provide a multi-stage `Dockerfile` that builds a self-contained, statically compiled Linux binary of the `store_product` service and bundles necessary runtime libraries.

#### Scenario: Build image from clean repository
- **WHEN** a user executes `docker build -t store_product .` on the repository
- **THEN** the builder stage downloads Go dependencies, generates the Linux Prisma Query Engine, compiles the Go application binary, and copies the resulting binary into a minimal Alpine runtime image without build dependencies.

### Requirement: Exclusion of Local Host Artifacts
The project SHALL provide a `.dockerignore` file preventing host-specific binary artifacts, environment secrets, git directories, and caches from being transferred to the Docker build context.

#### Scenario: Ignore local artifacts
- **WHEN** the Docker build context is sent to the Docker daemon
- **THEN** `.env`, `.git`, `.agent`, `openspec`, and Windows query engine files (`internal/db/*_gen.go`) are excluded from the build context.

### Requirement: Non-Root Execution Security
The container runtime SHALL execute under an unprivileged non-root user.

#### Scenario: Verify user in container
- **WHEN** the container is started
- **THEN** the application process runs under UID/GID `10001` (`appuser:appgroup`) instead of root.

### Requirement: Liveness and Health Probe
The container image SHALL define a built-in `HEALTHCHECK` directive probing the service's HTTP health endpoint.

#### Scenario: Container health reporting
- **WHEN** the container is running and healthy
- **THEN** the Docker daemon marks the container status as `healthy` by querying `GET /health` on the configured port.

### Requirement: External Service Configuration via Environment Variables
The containerized service SHALL accept database connection strings, Redis host/port, and server configurations through standard environment variables.

#### Scenario: Connect to pre-existing Database and Redis
- **WHEN** the container is started with `DATABASE_URL` and `REDIS_HOST` pointing to external database and Redis instances
- **THEN** the service successfully establishes connections to PostgreSQL and Redis without requiring embedded database or Redis daemon containers.
