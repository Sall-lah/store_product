## Why

The `store_product` microservice implements high-performance product catalog and variant management with PostgreSQL, Redis caching, keyset pagination, gateway offloading authentication, and embedded OpenAPI/Scalar/Swagger documentation. However, the repository lacks a comprehensive root `README.md` and `.env.example` file, making it difficult for new developers and operators to set up, configure, run, and integrate with the service. Adding complete project documentation and configuration templates now standardizes developer onboarding and production deployment workflows.

## What Changes

- Add comprehensive root `README.md` documenting architecture, features, prerequisites, configuration, interactive API documentation (Scalar UI & Swagger UI), gateway offloading authentication headers, keyset pagination, local development, Docker deployment, and testing instructions.
- Add `.env.example` template with documented configuration variables for database, Redis, port, environment, and rate limiting limits.
- Establish the `service-documentation` specification to formally govern service documentation standards and onboarding assets.

## Capabilities

### New Capabilities
- `service-documentation`: Comprehensive project onboarding, architecture overview, configuration guide, API documentation references, and deployment guides.

### Modified Capabilities
*(None)*

## Impact

- **Documentation**: Adds root `README.md` and `.env.example`.
- **APIs & Dependencies**: No breaking API changes or new Go dependencies required.
- **Developer Experience**: Significantly streamlines local environment setup, configuration discovery, API exploration via `/docs` and `/swagger`, and containerized deployment.
