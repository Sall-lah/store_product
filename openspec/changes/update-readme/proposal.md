## Why

The current `store_product` README provides useful foundational documentation, but lacks consistent structure, modern badges, detailed Kafka inventory event flow diagrams, and a comprehensive API catalog matching the enterprise standard established across the ecosystem (e.g., `store_order`). Standardizing the documentation improves developer onboarding, operational clarity, and architectural visibility.

## What Changes

- Update `README.md` to follow the standardized enterprise microservice README template.
- Add modern shield badges for Go version, Chi v5 router, PostgreSQL, Prisma Go Client, Kafka event streaming, Redis cache/rate limiting, and Cloudflare R2 storage.
- Include a structured Table of Contents with emoji headers linking to all major sections.
- Enhance the Architecture Overview with Mermaid diagrams for request flows and the event-driven inventory synchronization lifecycle.
- Provide a structured Key Features summary detailing atomic stock management, Cloudflare R2 presigned uploads, multi-tiered Redis caching & sliding-window rate limiting, and embedded OpenAPI 3.1 docs (Scalar/Swagger).
- Provide an updated Technology Stack list, ASCII repository tree diagram, and complete environment configuration table (`.env`).
- Standardize Database Setup & Prisma ORM generation workflows, Local Development Quickstart, and Docker Deployment instructions.
- Provide a complete Endpoint Catalog table covering Public, Admin, and Media Management routes.
- Document Kafka Event Processing details (`order.events` consumption, idempotency tracking via `processed_events`, cache invalidation triggers) and Redis rate limiting / degraded fail-open policies.

## Capabilities

### New Capabilities
None.

### Modified Capabilities
- `service-documentation`: Update requirements for the root `README.md` to include standardized badge styling, expanded Mermaid architecture/lifecycle diagrams, Cloudflare R2 storage documentation, Kafka inventory synchronization mechanics, and an exhaustive API catalog.

## Impact

- **Documentation**: Updates `README.md` with complete, up-to-date documentation.
- **Codebase**: No code changes required; all documented behaviors reflect the existing implementation in `store_product`.
