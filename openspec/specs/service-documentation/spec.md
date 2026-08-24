# Service Documentation Specification

## Purpose
Provides comprehensive, standardized documentation, onboarding assets, and configuration templates for the `store_product` microservice.

## Requirements

### Requirement: Comprehensive Root README Documentation
The repository SHALL maintain a root `README.md` document that serves as the single source of truth for microservice onboarding, architectural design, environment configuration, endpoint exploration, local development, Docker deployment, and testing.

#### Scenario: Developer reads project overview and architecture
- **WHEN** a developer accesses the root `README.md`
- **THEN** it SHALL provide a clear architecture breakdown detailing Go 1.22+, Chi routing, PostgreSQL via Prisma Client Go, Redis caching, Keyset Cursor pagination, Gateway offloading authentication, and embedded Scalar/Swagger API documentation

#### Scenario: Developer configures environment variables
- **WHEN** a developer inspects the configuration section of `README.md`
- **THEN** it SHALL document all configuration variables (`PORT`, `ENV`, `DATABASE_URL`, `REDIS_HOST`, `REDIS_PORT`, `REDIS_PASSWORD`, `REDIS_DB`, `RATE_LIMIT_PUBLIC_RPM`, `RATE_LIMIT_SEARCH_RPM`, `RATE_LIMIT_ADMIN_RPM`) with types, defaults, and descriptions

#### Scenario: Developer follows local development and Docker deployment instructions
- **WHEN** a developer or DevOps engineer sets up the service
- **THEN** it SHALL provide step-by-step instructions for native Go execution, multi-stage Docker build, and Docker container execution

### Requirement: Environment Variable Configuration Template
The repository SHALL provide a `.env.example` file documenting all supported environment variables with default values and descriptive comments.

#### Scenario: Developer copies configuration template
- **WHEN** a developer inspects `.env.example`
- **THEN** it SHALL contain sample values for server port, environment, database connection, Redis configuration, and rate limiting thresholds without containing real secrets or API keys

### Requirement: API Integration and Gateway Offloading Documentation
The documentation SHALL specify integration protocols including Gateway Offloading headers, rate limiting headers, and Keyset pagination mechanics.

#### Scenario: API consumer integrates with authenticated endpoints
- **WHEN** an API developer reviews documentation for admin endpoints
- **THEN** it SHALL document the required `X-User-Role` (case-insensitive `admin`) and `X-User-Id` headers injected by the API gateway and the rate limit headers (`X-RateLimit-Limit`, `X-RateLimit-Remaining`, `X-RateLimit-Reset`)

#### Scenario: API consumer navigates interactive API documentation
- **WHEN** a developer seeks interactive endpoint documentation
- **THEN** `README.md` SHALL link to and describe the `/docs` (Scalar UI), `/swagger` (Swagger UI), `/openapi.json`, and `/openapi.yaml` endpoints
