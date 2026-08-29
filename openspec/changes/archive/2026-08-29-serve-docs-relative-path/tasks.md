## 1. Documentation HTML & UI Templates

- [x] 1.1 Update `ServeSwaggerUI` in `internal/handler/docs.go` to use relative spec URL `./openapi.json`
- [x] 1.2 Update `ServeScalarUI` in `internal/handler/docs.go` to use relative spec URL `data-url="./openapi.json"`

## 2. Router Sub-route Spec Mounting

- [x] 2.1 Update `internal/handler/router.go` to expose `openapi.json` and `openapi.yaml` under `/docs` and `/swagger` sub-routes
- [x] 2.2 Ensure wildcard doc routes do not shadow spec routes

## 3. OpenAPI Specification Definitions

- [x] 3.1 Add relative server URL `./` to `servers` array in `docs/openapi.json`
- [x] 3.2 Add relative server URL `./` to `servers` array in `docs/openapi.yaml`

## 4. Testing & Verification

- [x] 4.1 Update `internal/handler/docs_test.go` to verify Swagger UI and Scalar UI contain relative `./openapi.json` references
- [x] 4.2 Add test cases in `internal/handler/docs_test.go` verifying `GET /docs/openapi.json` and `GET /swagger/openapi.json` return OpenAPI spec
- [x] 4.3 Run `go test ./...` to verify all tests pass
