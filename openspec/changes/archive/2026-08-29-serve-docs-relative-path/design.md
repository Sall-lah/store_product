## Context

`store_product` provides embedded API documentation using Swagger UI (`/swagger`) and Scalar UI (`/docs`). Currently, both UIs load `/openapi.json` using an absolute root path. When running behind reverse proxies or API gateways with path prefixes (e.g. `https://gateway.domain/products/docs`), the browser requests `/openapi.json` from the domain root rather than the proxy prefix, causing 404 Not Found errors.

## Goals / Non-Goals

**Goals:**
- Enable Swagger UI and Scalar UI to successfully fetch OpenAPI schemas when mounted at root or behind reverse proxy subpath prefixes.
- Handle both trailing slash (`/docs/`, `/swagger/`) and non-trailing slash (`/docs`, `/swagger`) URL requests without breaking spec resolution.
- Configure OpenAPI 3.1 server definitions with relative path (`./`) so interactive endpoint calls in Swagger/Scalar target the gateway subpath.
- Maintain backwards compatibility for direct local development access at `http://localhost:8080/docs` and `http://localhost:8080/swagger`.

**Non-Goals:**
- Self-hosting CDN JS/CSS assets (Swagger UI and Scalar bundle assets remain loaded from unpkg/jsdelivr CDNs).
- Altering API route structure under `/api/v1`.

## Decisions

### Decision 1: Relative spec references in UI HTML
- **Choice**: Update `internal/handler/docs.go` to use `url: "./openapi.json"` in `ServeSwaggerUI` and `data-url="./openapi.json"` in `ServeScalarUI`.
- **Rationale**: Browsers resolve `./openapi.json` relative to the current URL path rather than resetting to the root domain.
- **Alternatives Considered**: Injecting dynamic JavaScript `window.location.origin` - rejected because it still fails when a path prefix is present without parsing the path segments.

### Decision 2: Dual-mount spec endpoints in Router
- **Choice**: Expose `openapi.json` and `openapi.yaml` at root (`/openapi.json`) and within `/docs/` (`/docs/openapi.json`) and `/swagger/` (`/swagger/openapi.json`).
- **Rationale**: According to RFC 3986, `./openapi.json` from `/docs` resolves to `/openapi.json`, whereas from `/docs/` it resolves to `/docs/openapi.json`. Mounting the handlers at both routes ensures 100% resolution reliability regardless of client trailing slash habits.
- **Alternatives Considered**: Enforcing 301/302 redirects between `/docs` and `/docs/` - rejected because reverse proxies often mishandle or rewrite redirect `Location` headers.

### Decision 3: Relative server entry in OpenAPI 3.1 definitions
- **Choice**: Add an entry with `url: "./"` and description `"API Gateway / Relative Subpath"` in `docs/openapi.json` and `docs/openapi.yaml`.
- **Rationale**: In OpenAPI 3.1, relative server URLs are resolved relative to the document location, ensuring interactive "Try It Out" requests preserve the gateway subpath prefix.

## Risks / Trade-offs

- **[Risk] Route matching precedence in Chi**: If `/docs/*` or `/swagger/*` catch-all wildcards precede `/docs/openapi.json`, the HTML UI would be returned instead of JSON.
  → **Mitigation**: Register explicit sub-route `/docs/openapi.json` and `/docs/openapi.yaml` before wildcard or as dedicated sub-routes.
- **[Risk] Relative server resolution in older tools**: Some legacy OpenAPI 2.0 tools do not support relative `servers.url`.
  → **Mitigation**: We retain `http://localhost:8080` alongside `./` in `servers`.
