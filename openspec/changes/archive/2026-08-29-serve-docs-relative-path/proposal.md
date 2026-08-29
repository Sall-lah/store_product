## Why

When the `store_product` microservice is deployed behind reverse proxies, Kubernetes Ingresses, or API Gateways (such as `store_gateway`) mounted under a subpath prefix (e.g. `/products` or `/catalog`), the hardcoded absolute path `/openapi.json` causes Swagger UI and Scalar UI to request the OpenAPI specification from the root domain, resulting in 404 Not Found errors. Updating documentation serving to use relative paths (`./openapi.json` and relative server URLs `./`) ensures seamless documentation exploration and interactive API testing regardless of where the microservice is mounted.

## What Changes

- **Relative OpenAPI Spec Path in Swagger UI**: Update Swagger UI HTML in `internal/handler/docs.go` to reference `./openapi.json` instead of absolute `/openapi.json`.
- **Relative OpenAPI Spec Path in Scalar UI**: Update Scalar UI HTML in `internal/handler/docs.go` to reference `data-url="./openapi.json"` instead of absolute `/openapi.json`.
- **Sub-Route Spec Registration**: In `internal/handler/router.go`, mount `openapi.json` and `openapi.yaml` within `/docs/` and `/swagger/` route groups so that relative path `./openapi.json` resolves cleanly whether accessed with or without a trailing slash (`/docs` vs `/docs/`).
- **Relative OpenAPI Server URL**: Add relative server URL `./` to `docs/openapi.json` and `docs/openapi.yaml` so interactive API executions ("Try It Out") target the current gateway subpath rather than domain root.
- **Unit and Integration Tests**: Update documentation route tests in `internal/handler/docs_test.go` to assert relative path references and verify sub-route spec resolution.

## Capabilities

### Modified Capabilities
- `swagger-scalar-ui`: Update Swagger and Scalar documentation interfaces and routing to resolve OpenAPI specifications via relative paths (`./openapi.json`) across root and sub-route endpoints.
- `openapi-documentation`: Update OpenAPI 3.1.0 server definitions to include relative base `./` for subpath-aware API Gateway offloading and exploration.

## Impact

- `internal/handler/docs.go`: Swagger UI and Scalar UI HTML generation templates.
- `internal/handler/router.go`: Chi router documentation endpoint definitions.
- `docs/openapi.json` & `docs/openapi.yaml`: OpenAPI 3.1 specification server configuration.
- `internal/handler/docs_test.go`: Endpoint assertions for relative paths.
