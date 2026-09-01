# Service Documentation Specification

## Purpose
Provides comprehensive, standardized documentation, onboarding assets, and configuration templates for the `store_product` microservice.

## Requirements

### Requirement: Comprehensive Root README Documentation
The repository SHALL maintain a root `README.md` document structured according to the enterprise microservice documentation template, serving as the single source of truth for onboarding, architectural design, environment configuration, endpoint exploration, Kafka event integration, Redis caching and rate limiting rules, Cloudflare R2 media management, local development, Docker deployment, and testing.

#### Scenario: Developer reads project overview and architecture
- **WHEN** a developer accesses the root `README.md`
- **THEN** it SHALL display modern status badges (Go 1.26+, Chi v5, PostgreSQL, Prisma Go Client, Apache Kafka, Redis, Cloudflare R2), a structured Table of Contents, an ASCII repository tree, and Mermaid diagrams illustrating HTTP request flows and Kafka event-driven inventory synchronization

#### Scenario: Developer inspects key features and technical capabilities
- **WHEN** a developer reviews the key features section
- **THEN** it SHALL outline atomic stock deduction and restocking, Cloudflare R2 presigned image uploads and gallery management, multi-tiered Redis caching and sliding-window rate limiting with fail-open circuit breakers, API Gateway offloading authentication, and embedded OpenAPI 3.1 Scalar/Swagger documentation

#### Scenario: Developer configures environment variables
- **WHEN** a developer inspects the configuration section of `README.md`
- **THEN** it SHALL document all configuration variables (`PORT`, `ENV`, `DATABASE_URL`, `REDIS_HOST`, `REDIS_PORT`, `REDIS_PASSWORD`, `REDIS_DB`, `RATE_LIMIT_PUBLIC_RPM`, `RATE_LIMIT_SEARCH_RPM`, `RATE_LIMIT_ADMIN_RPM`, `KAFKA_BROKERS`, `KAFKA_TOPIC_ORDER_EVENTS`, `KAFKA_CONSUMER_GROUP`, `R2_ACCOUNT_ID`, `R2_ACCESS_KEY_ID`, `R2_SECRET_ACCESS_KEY`, `R2_BUCKET_NAME`, `R2_PUBLIC_BASE_URL`) with types, defaults, and descriptions

#### Scenario: Developer navigates endpoint catalog and Kafka pipeline
- **WHEN** an API consumer or backend engineer explores integration points
- **THEN** it SHALL provide a full endpoint catalog table (Public, Admin, Media) with HTTP methods, paths, authentication requirements, and rate limit rules, along with Kafka event topic schemas and idempotency ledger specifications

#### Scenario: Developer follows local development and Docker deployment instructions
- **WHEN** a developer or DevOps engineer sets up or packages the service
- **THEN** it SHALL provide step-by-step instructions for Prisma schema push and Go client generation, native Go execution, testing suites, and multi-stage Docker container build and runtime commands

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
