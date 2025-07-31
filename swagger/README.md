# API Documentation

This directory contains auto-generated API documentation using Swaggo.

## Installation

Install swag CLI tool:
```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

## Generate Documentation

```bash
make swag
```

This will generate:
- `docs/swagger.json` - OpenAPI specification
- `docs/swagger.yaml` - OpenAPI specification (YAML format)
- `docs/docs.go` - Go code for embedding documentation

## View Documentation

1. Start the server:
   ```bash
   gh-inspector serve
   ```

2. Open Swagger UI:
   http://localhost:8080/swagger/index.html

## Adding Documentation

Add Swagger annotations to your code:

```go
// @Summary Score repositories
// @Description Analyze GitHub repositories
// @Tags analysis
// @Accept json
// @Produce json
// @Param request body ScoreRequest true "Repositories to score"
// @Success 200 {object} ScoreResponse
// @Router /api/v1/score [post]
func handleScore(w http.ResponseWriter, r *http.Request) {
    // ...
}
```

Then regenerate documentation with `make swag`.