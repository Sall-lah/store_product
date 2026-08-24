## Context

In the microservices architecture, identity verification and token inspection are offloaded to `store_gateway`. The gateway validates incoming Bearer JWT tokens with `store_auth` (or internal token verification) and injects trusted identity headers (`X-User-Role`, `X-User-Id`) when proxying requests to downstream services such as `store_product`.

Currently:
1. `store_auth` returns uppercase role strings (`"ADMIN"`).
2. `store_product` uses a strict equality check (`role != "admin"`), causing valid admin requests routed via the gateway to fail with `403 Forbidden`.
3. The OpenAPI 3.1.0 specification and Swagger UI documentation currently specify exact lowercase `"admin"` and lack clear architectural notes on Gateway Offloading.
4. Database connection strings in `.env` / `.env.example` need clear differentiation between local containerized PostgreSQL and Supabase poolers.

## Goals / Non-Goals

**Goals:**
- Implement case-insensitive and whitespace-tolerant role validation in `RequireAdmin` middleware.
- Ensure unit tests validate uppercase (`"ADMIN"`), lowercase (`"admin"`), and mixed-case (`"Admin"`) headers.
- Update OpenAPI 3.1 specification (`openapi.yaml` and `openapi.json`) to accurately describe the Gateway Offloading security headers and casing rules.
- Clarify PostgreSQL database connection configuration for both local Docker and cloud environments.

**Non-Goals:**
- Do not perform direct JWT decoding or cryptographic verification inside `store_product` (strictly preserve the Gateway Offloading pattern).
- Do not alter `store_auth` token issuance schemas.

## Decisions

### Decision 1: Case-Insensitive Role Comparison via `strings.EqualFold`
- **Choice**: Use `strings.EqualFold(strings.TrimSpace(role), "admin")`.
- **Rationale**: `EqualFold` provides idiomatic, high-performance, Unicode case-folding comparison in Go without creating auxiliary lowercase string allocations. Trimming leading/trailing whitespace prevents whitespace bugs from upstream proxy header formatting.
- **Alternatives Considered**:
  - `strings.ToLower(role) == "admin"`: Introduces an extra string allocation per request.
  - Modifying `store_auth` exclusively: Fragile; downstream services should follow Postel's Law (be liberal in what you accept).

### Decision 2: OpenAPI 3.1 / Swagger Documentation Synchronization
- **Choice**: Maintain single source of truth across `docs/openapi.yaml` and embedded `docs/openapi.json`.
- **Rationale**: `store_product` embeds both files via `go:embed` in `docs/docs.go`. Updating both ensures consistency across Swagger UI, Scalar UI, and CLI tooling.

## Risks / Trade-offs

- **[Risk]** Downstream services receiving spoofed headers if exposed directly to the public internet.
  - **Mitigation**: In production environments, `store_product` is hosted within a private Docker/Kubernetes network where only `store_gateway` can route traffic. The gateway strips incoming client headers before injecting validated claims.
