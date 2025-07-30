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
```bash
docker build -t gh-inspector .
docker run --rm gh-inspector score --repos=your/repo
```
