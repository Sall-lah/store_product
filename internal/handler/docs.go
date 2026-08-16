package handler

import (
	"net/http"

	"github.com/Sall-lah/store_product/docs"
)

// DocsHandler handles requests for OpenAPI documentation specifications and interactive UIs.
// Encapsulating documentation delivery within a dedicated handler ensures modularity and testability.
type DocsHandler struct{}

// NewDocsHandler instantiates a new DocsHandler instance.
func NewDocsHandler() *DocsHandler {
	return &DocsHandler{}
}

// ServeOpenAPIJSON outputs the raw OpenAPI 3.1.0 JSON specification.
// Returning the embedded byte slice directly provides sub-millisecond responses without disk I/O.
func (h *DocsHandler) ServeOpenAPIJSON(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(docs.OpenAPIJSON)
}

// ServeOpenAPIYAML outputs the raw OpenAPI 3.1.0 YAML specification.
// Supporting YAML alongside JSON accommodates diverse developer tooling (e.g. Swagger CLI, Postman).
func (h *DocsHandler) ServeOpenAPIYAML(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/yaml; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(docs.OpenAPIYAML)
}

// ServeSwaggerUI renders the interactive Swagger UI standalone web interface.
// Using CDN-hosted assets eliminates the need to bundle static CSS/JS files into the Go binary.
func (h *DocsHandler) ServeSwaggerUI(w http.ResponseWriter, r *http.Request) {
	html := `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <title>Store Product Microservice - Swagger UI</title>
  <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5.18.2/swagger-ui.css" />
  <link rel="icon" type="image/png" href="https://unpkg.com/swagger-ui-dist@5.18.2/favicon-32x32.png" sizes="32x32" />
  <style>
    html { box-sizing: border-box; overflow: -moz-scrollbars-vertical; overflow-y: scroll; }
    *, *:before, *:after { box-sizing: inherit; }
    body { margin:0; background: #fafafa; font-family: sans-serif; }
    .topbar { display: none !important; }
  </style>
</head>
<body>
  <div id="swagger-ui"></div>
  <script src="https://unpkg.com/swagger-ui-dist@5.18.2/swagger-ui-bundle.js" charset="UTF-8"></script>
  <script src="https://unpkg.com/swagger-ui-dist@5.18.2/swagger-ui-standalone-preset.js" charset="UTF-8"></script>
  <script>
    window.onload = function() {
      window.ui = SwaggerUIBundle({
        url: "/openapi.json",
        dom_id: '#swagger-ui',
        deepLinking: true,
        presets: [
          SwaggerUIBundle.presets.apis,
          SwaggerUIStandalonePreset
        ],
        plugins: [
          SwaggerUIBundle.plugins.DownloadUrl
        ],
        layout: "StandaloneLayout",
        persistAuthorization: true,
        displayRequestDuration: true
      });
    };
  </script>
</body>
</html>`

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(html))
}

// ServeScalarUI renders the modern Scalar API reference interface.
// Scalar provides deep native OpenAPI 3.1 support, search, and multi-language code snippets.
func (h *DocsHandler) ServeScalarUI(w http.ResponseWriter, r *http.Request) {
	html := `<!DOCTYPE html>
<html>
<head>
  <title>Store Product Microservice - API Reference</title>
  <meta charset="utf-8" />
  <meta name="viewport" content="width=device-width, initial-scale=1" />
  <style>
    body { margin: 0; }
  </style>
</head>
<body>
  <script
    id="api-reference"
    data-url="/openapi.json"
    data-configuration='{"theme":"purple","layout":"modern"}'
    src="https://cdn.jsdelivr.net/npm/@scalar/api-reference">
  </script>
</body>
</html>`

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(html))
}
