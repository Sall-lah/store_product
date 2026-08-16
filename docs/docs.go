// Package docs provides embedded OpenAPI 3.1.0 specifications and API documentation assets.
package docs

import (
	_ "embed"
)

// OpenAPIJSON contains the raw embedded OpenAPI 3.1.0 JSON specification.
// Embedding at compile time ensures zero filesystem dependency in containerized environments.
//
//go:embed openapi.json
var OpenAPIJSON []byte

// OpenAPIYAML contains the raw embedded OpenAPI 3.1.0 YAML specification.
//
//go:embed openapi.yaml
var OpenAPIYAML []byte
