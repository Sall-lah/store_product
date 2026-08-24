# Design Document: Containerized Docker Deployment

## Context

The `store_product` microservice is built using Go (Chi router), Prisma Client Go (`github.com/steebchen/prisma-client-go`), and integrates with an existing PostgreSQL database (e.g. Supabase) and Redis cache.

When running in production or standard container environments, the service requires:
1. Compiling for Linux with Prisma Query Engine binaries generated specifically for Linux/musl.
2. Including dynamic library dependencies (`openssl`, `ca-certificates`) required by Prisma's query engine and TLS database connections.
3. Operating without bundling or spawning database and Redis daemons inside the container, as external managed or self-hosted instances are already provisioned.

## Goals / Non-Goals

**Goals:**
- Provide a clean, reproducible multi-stage `Dockerfile`.
- Keep image size minimal using `alpine:3.21` (~35MB total image size).
- Isolate the build stage to generate Linux Prisma client engines and avoid host OS cross-compilation conflicts.
- Implement security best practices: non-root user execution (`appuser`), stripped binary symbols (`-ldflags="-w -s"`), and health checking.
- Provide a comprehensive `.dockerignore` file.

**Non-Goals:**
- Provisioning PostgreSQL or Redis containers (the user manages them externally).
- Modifying internal Go source code or API handler signatures.

## Decisions

### 1. Multi-Stage Build with Alpine Base
- **Decision**: Use `golang:1.24-alpine` for the builder stage and `alpine:3.21` for the final runtime stage.
- **Rationale**: Minimal attack surface, fast image pull/startup, low disk footprint.
- **Alternatives Considered**:
  - `distroless`: Harder to bundle Prisma query engine runtime dependencies (like OpenSSL / dynamic links) and lack `curl` for basic health checking.
  - `debian:bookworm-slim`: Works well with Prisma but produces larger images (~120MB vs ~35MB).

### 2. Linux Prisma Query Engine Generation in Container
- **Decision**: Execute `go run github.com/steebchen/prisma-client-go generate` inside the Docker builder container.
- **Rationale**: Prisma Client Go packages a platform-specific query engine binary. Generating inside the Linux builder container guarantees the engine binary matches the target architecture and musl C library.
- **Alternatives Considered**:
  - Pre-generating engine binaries on the host: Fails if the host is Windows/macOS.

### 3. OpenSSL & CA Certificates in Runtime Stage
- **Decision**: Install `ca-certificates`, `openssl`, `tzdata`, and `curl` in the Alpine runtime container.
- **Rationale**:
  - `openssl` is required by the Prisma query engine binary.
  - `ca-certificates` is required for TLS database connections (Supabase/Neon/RDS).
  - `curl` provides standard HTTP health check execution.

## Risks / Trade-offs

- **[Risk] Prisma Query Engine Engine Download during Build** → Network access is required during `docker build` to download the Prisma query engine.
  - *Mitigation*: The `RUN go run ... generate` step is cached by Docker layer caching when `schema.prisma` and `go.sum` do not change.
- **[Risk] Connecting to Localhost DB/Redis on Host Machine** → Containers cannot connect to host's `localhost` directly.
  - *Mitigation*: Document use of `host.docker.internal` for Docker Desktop (Windows/Mac) or `--network="host"` on Linux when running with local host databases.
