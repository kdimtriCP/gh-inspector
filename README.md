# gh-inspector

[![CI](https://github.com/kdimtricp/gh-inspector/actions/workflows/ci.yml/badge.svg)](https://github.com/kdimtricp/gh-inspector/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/kdimtricp/gh-inspector)](https://goreportcard.com/report/github.com/kdimtricp/gh-inspector)

`gh-inspector` is a CLI tool for analyzing and scoring GitHub repositories.  
It helps you quickly assess the quality and activity of multiple projects before choosing a library, dependency, or contributing to open source.

---

## 🧠 What it does

- Accepts a list of public GitHub repositories
- Collects key metrics:
    - stars, forks
    - last activity
    - open issues and pull requests
    - presence of CI/CD, license, contributing guide
- Builds a comparison table
- Provides a final score for each repository

---

## 🚀 Usage Example

```bash
gh-inspector score --repos=go-chi/chi,labstack/echo,gin-gonic/gin
Analyzing repo: go-chi/chi
Analyzing repo: labstack/echo
Analyzing repo: gin-gonic/gin

+------------------+--------+--------------+------+--------+--------+
| Repository       | Stars  | Last Commit  | CI   | Issues | Score  |
+------------------+--------+--------------+------+--------+--------+
| go-chi/chi       | 14.9k  | 2024-12-01   | ❌   | 21     | 68     |
| labstack/echo    | 27.2k  | 2025-05-10   | ✅   | 84     | 78     |
| gin-gonic/gin    | 73.4k  | 2025-07-15   | ✅   | 200+   | 91     |
+------------------+--------+--------------+------+--------+--------+
```

## Installation
```bash
go install github.com/kdimtriCP/gh-inspector@latest
```

## Docker:

### Production Build
```bash
# Build with auto-generated Swagger documentation
docker build -t gh-inspector .

# Run CLI commands
docker run --rm gh-inspector score --repos=your/repo

# Run web service
docker run -p 8080:8080 -e GITHUB_TOKEN=$GITHUB_TOKEN gh-inspector serve
```

### Docker Compose
```bash
# Start all services (gh-inspector, prometheus, grafana)
docker-compose up -d

# View logs
docker-compose logs -f gh-inspector

# Stop all services
docker-compose down
```

## Web Service Mode

Run gh-inspector as a web service:

```bash
# Start the server
gh-inspector serve --port 8080

# Or using Docker Compose
docker-compose up -d
```

### API Endpoints

- `POST /api/v1/score` - Score GitHub repositories
- `GET /health` - Health check endpoint
- `GET /metrics` - Prometheus metrics
- `GET /` - API information
- `GET /swagger/index.html` - Swagger UI documentation

### Example API Usage

```bash
# Score repositories via API
curl -X POST http://localhost:8080/api/v1/score \
  -H "Content-Type: application/json" \
  -d '{
    "repositories": [
      "kubernetes/kubernetes",
      "golang/go"
    ]
  }'

# Check health
curl http://localhost:8080/health
```

### API Documentation

The API documentation is auto-generated using Swagger/OpenAPI:

1. **Generate documentation**:
   ```bash
   make swag
   ```

2. **Access Swagger UI**:
   - Start the server: `gh-inspector serve`
   - Open browser: http://localhost:8080/swagger/index.html

3. **View OpenAPI spec**:
   - JSON: http://localhost:8080/swagger/doc.json
