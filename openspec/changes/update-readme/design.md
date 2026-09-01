## Context

The `store_product` microservice is the product catalog, variant, stock, and media asset authority within the e-commerce platform. Other microservices in the platform (such as `store_order`) have standardized on a modern, high-clarity README structure featuring visual badges, navigable Table of Contents, Mermaid flowcharts for service interactions and event lifecycles, structured API catalog tables, Kafka topic and idempotency documentation, Redis cache/rate limit specifications, and deployment runbooks.

Updating `store_product` to match this template ensures developer consistency, architectural clarity, and smooth onboarding across all platform microservices.

## Goals / Non-Goals

**Goals:**
- Update `README.md` to follow the exact visual styling, structure, and depth demonstrated in the template.
- Document all core capabilities of `store_product`:
  - Architecture diagram with Gateway offloading, Redis caching/rate limiting, and Kafka inventory synchronization.
  - Event-Driven Inventory Lifecycle State Machine (`order.created`/`order.placed` decrement, `order.cancelled`/`order.expired` restock) with database idempotency (`processed_events`).
  - Key Features list highlighting atomic inventory mutations, Cloudflare R2 presigned media uploads, Keyset Cursor pagination, Redis multi-level caching, and embedded OpenAPI 3.1 (Scalar UI / Swagger UI).
  - Accurate Go 1.26+, Chi v5, Prisma Client Go, PostgreSQL, Kafka, Redis, and Cloudflare R2 badges.
  - Comprehensive configuration table covering all `.env` parameters.
  - Prisma ORM database migration and client generation instructions.
  - Full API Endpoint Catalog table categorized by Public, Admin, and Image routes.
  - Object storage workflow explanation for direct-to-R2 presigned URL uploads.
  - Testing commands and multi-stage Docker deployment instructions.

**Non-Goals:**
- Modifying Go code or service runtime behavior (this is purely a documentation alignment change).
- Changing OpenAPI specification files or `.env.example` templates.

## Decisions

1. **Badge Alignment**:
   - Use shields.io badges matching the `store_order` aesthetic for Go version (`Go-1.26+`), HTTP Router (`Router-Chi%20v5`), Database (`Database-PostgreSQL`), ORM (`ORM-Prisma%20Go%20Client`), Event Streaming (`Streaming-Apache%20Kafka`), Caching & Rate Limiting (`Cache%20&%20Rate%20Limit-Redis`), and Media Storage (`Storage-Cloudflare%20R2`).
   - *Rationale*: Establishes visual uniformity across all repositories.

2. **Mermaid Diagrams**:
   - Create a service architecture flowchart mapping Client -> API Gateway -> Chi Router -> Middleware -> Handlers -> Services -> Repositories -> PostgreSQL/Redis/R2/Kafka.
   - Create an Event-Driven Stock Synchronization flow / state diagram detailing how inbound `order.events` trigger stock reservation, release, and cache invalidation.
   - *Rationale*: Visual diagrams drastically reduce onboarding cognitive load.

3. **Complete Endpoint Catalog**:
   - Organize routes clearly into Public Catalog, Admin Product & Variant Management, and Admin Image & Media Management.
   - Include Method, Path, Auth/Headers (`X-User-Role: admin`), Rate Limit, and Description.
   - *Rationale*: Provides quick reference without needing to immediately launch Swagger/Scalar.

4. **Dedicated Sections for Domain-Specific Mechanics**:
   - Detail Kafka event handling (`order.events`), idempotency checks via `processed_events`, and cache purging (`product:detail:<id>`, `product:slug:<slug>`, `product:list:*`).
   - Detail Cloudflare R2 presigned URL upload flow (Request presigned PUT URL -> Client direct upload to R2 -> Register metadata in Postgres).
   - *Rationale*: Captures operational details necessary for platform integration.

## Risks / Trade-offs

- **[Risk]** Out-of-sync documentation if environment variables or endpoints change in the future.
  - **Mitigation**: All documented endpoints and environment variables are verified directly against `internal/config/config.go`, `internal/handler/router.go`, and `cmd/server/main.go`.
