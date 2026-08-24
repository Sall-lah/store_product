## ADDED Requirements

### Requirement: Interactive Swagger and Scalar Documentation Endpoints
The system SHALL expose HTTP endpoints to serve interactive documentation user interfaces and raw OpenAPI specification files embedded into the application binary.

#### Scenario: Swagger UI accessibility
- **WHEN** a client performs a `GET` request to `/swagger` or `/swagger/`
- **THEN** the server SHALL respond with `200 OK`, `Content-Type: text/html; charset=utf-8`, and render the interactive Swagger UI interface pointing to `/openapi.json`

#### Scenario: Scalar UI accessibility
- **WHEN** a client performs a `GET` request to `/docs` or `/docs/`
- **THEN** the server SHALL respond with `200 OK`, `Content-Type: text/html; charset=utf-8`, and render the modern Scalar UI interface with theme toggle and multi-language code snippets

#### Scenario: Raw OpenAPI specification retrieval
- **WHEN** a client performs a `GET` request to `/openapi.json`
- **THEN** the server SHALL respond with `200 OK`, `Content-Type: application/json; charset=utf-8`, and return the complete OpenAPI 3.1.0 document

#### Scenario: Raw OpenAPI YAML retrieval
- **WHEN** a client performs a `GET` request to `/openapi.yaml`
- **THEN** the server SHALL respond with `200 OK`, `Content-Type: application/yaml; charset=utf-8`, and return the complete YAML specification
