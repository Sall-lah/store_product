## ADDED Requirements

### Requirement: Relative Server URL Definition
The OpenAPI 3.1.0 specification SHALL declare a relative server target (`./`) in its `servers` configuration to support reverse proxy and API Gateway subpath prefix execution.

#### Scenario: Specification includes relative server URL
- **WHEN** a client inspects the `servers` array in the OpenAPI specification
- **THEN** it SHALL include a server entry with `url: "./"` and description indicating relative gateway host execution
