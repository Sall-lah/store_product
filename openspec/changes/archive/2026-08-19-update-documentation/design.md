## Context

The `store_product` microservice handles product catalog querying and variant management in an e-commerce microservices ecosystem. It utilizes Go (1.22+), Chi router, Prisma Client Go with PostgreSQL, Redis for caching and sliding-window rate limiting, and gateway offloading for admin authentication. Additionally, it embeds OpenAPI 3.1 specifications and serves interactive Scalar UI (`/docs`) and Swagger UI (`/swagger`).

Prior changes implemented and verified these features, but the repository lacks a primary entry point document (`README.md`) and an environment variable template (`.env.example`). This change creates authoritative, structured documentation and configuration templates to facilitate rapid onboarding, local development, and operational deployments.

## Goals / Non-Goals

**Goals:**
- Provide a comprehensive, developer-friendly root `README.md` covering architecture, tech stack, configuration, quick start, Docker workflow, API endpoints, gateway auth, keyset pagination, and testing.
- Provide a clean `.env.example` template with documented defaults matching `internal/config/config.go`.
- Ensure documentation accurately reflects the current state of code (e.g. case-insensitive `X-User-Role: admin` gateway offloading, keyset pagination with base64 cursor, multi-layer cache invalidation, embedded Scalar/Swagger endpoints).

**Non-Goals:**
- Modifying internal Go application code, routing, or database schema.
- Creating runtime `.env` files with actual production secrets.

## Decisions

### Decision 1: Structured Root README.md Hierarchy
- **Rationale**: A standardized structure (Overview, Architecture, Tech Stack, Prerequisites, Configuration, Getting Started, Docker Deployment, API Documentation, Gateway Offloading & Auth, Keyset Pagination, Testing, Project Layout) enables developers to find required information immediately.
- **Alternatives Considered**: Creating multiple markdown docs in a `docs/` subfolder. Rejected in favor of a central `README.md` because `README.md` is the primary entry point rendered on GitHub and repository homepages.

### Decision 2: Environment Variable Template (.env.example)
- **Rationale**: An explicit `.env.example` file ensures zero-guesswork setup for local environments while strictly preventing the exposure of credentials.
- **Alternatives Considered**: Inlining environment variables only in `README.md`. Having a dedicated `.env.example` file allows developers to execute `cp .env.example .env` and immediately configure their local environment.

### Decision 3: Document API Gateway Offloading Architecture
- **Rationale**: The microservice relies on the API Gateway (`store_gateway`) to validate JWTs and forward `X-User-Role` and `X-User-Id` headers. Documenting this pattern clearly avoids developer confusion about missing JWT secret keys or auth middleware in the service.

## Risks / Trade-offs

- **[Risk] Accidental inclusion of sensitive credentials in example files** → *Mitigation*: Ensure all `.env.example` and `README.md` configurations use localhost and dummy placeholders.
- **[Risk] Documentation desynchronization as endpoints evolve** → *Mitigation*: Emphasize embedded live UI documentation at `/docs` (Scalar) and `/swagger` (Swagger) as the interactive source of truth.
