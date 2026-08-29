## MODIFIED Requirements

### Requirement: Interactive Swagger and Scalar Documentation Endpoints
The system SHALL expose HTTP endpoints to serve interactive documentation user interfaces and raw OpenAPI specification files embedded into the application binary, utilizing relative paths (`./openapi.json`) and sub-route spec registration for reverse proxy and gateway prefix compatibility.

#### Scenario: Swagger UI accessibility
- **WHEN** a client performs a `GET` request to `/swagger` or `/swagger/`
- **THEN** the server SHALL respond with `200 OK`, `Content-Type: text/html; charset=utf-8`, and render the interactive Swagger UI interface configured with relative specification URL `./openapi.json`

#### Scenario: Scalar UI accessibility
- **WHEN** a client performs a `GET` request to `/docs` or `/docs/`
- **THEN** the server SHALL respond with `200 OK`, `Content-Type: text/html; charset=utf-8`, and render the modern Scalar UI interface configured with relative specification URL `data-url="./openapi.json"`

#### Scenario: Raw OpenAPI specification retrieval
- **WHEN** a client performs a `GET` request to `/openapi.json`, `/docs/openapi.json`, or `/swagger/openapi.json`
- **THEN** the server SHALL respond with `200 OK`, `Content-Type: application/json; charset=utf-8`, and return the complete OpenAPI 3.1.0 document

#### Scenario: Raw OpenAPI YAML retrieval
- **WHEN** a client performs a `GET` request to `/openapi.yaml`, `/docs/openapi.yaml`, or `/swagger/openapi.yaml`
- **THEN** the server SHALL respond with `200 OK`, `Content-Type: application/yaml; charset=utf-8`, and return the complete YAML specification
