## 1. Environment Configuration Template

- [x] 1.1 Create `.env.example` with all configuration variables (`PORT`, `ENV`, `DATABASE_URL`, `REDIS_HOST`, `REDIS_PORT`, `REDIS_PASSWORD`, `REDIS_DB`, rate limits), safe defaults, and explanatory comments

## 2. Root README Documentation

- [x] 2.1 Create comprehensive `README.md` with project overview, architecture diagram/breakdown, key features, and tech stack
- [x] 2.2 Add Prerequisites, Environment Configuration guide, and Local Development Quick Start sections to `README.md`
- [x] 2.3 Add Docker deployment instructions, multi-stage build details, and container execution commands to `README.md`
- [x] 2.4 Document Interactive API Documentation (`/docs` Scalar UI, `/swagger` Swagger UI, raw `/openapi.json` and `/openapi.yaml` endpoints) in `README.md`
- [x] 2.5 Document API Gateway Offloading authentication headers (`X-User-Role`, `X-User-Id`), Redis sliding-window rate limiting, and Keyset Cursor Pagination mechanics in `README.md`
- [x] 2.6 Document Testing guidelines (`go test ./...`) and Project Directory Structure in `README.md`

## 3. Verification

- [x] 3.1 Verify markdown links, code blocks, tables, and consistency between `.env.example`, `internal/config/config.go`, and `README.md`
