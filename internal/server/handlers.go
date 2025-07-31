package server

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/kdimtriCP/gh-inspector/internal/formatter"
)

// ScoreRequest represents the request body for scoring repositories
// @Description Request body for scoring GitHub repositories
type ScoreRequest struct {
	// List of repositories in owner/name format
	// @example ["kubernetes/kubernetes", "golang/go"]
	Repositories []string `json:"repositories" example:"kubernetes/kubernetes,golang/go"`
	// Output format (optional)
	// @example json
	OutputFormat string `json:"output_format,omitempty" example:"json"`
}

// ScoreResponse represents the response from the score endpoint
// @Description Response containing scored repositories
type ScoreResponse struct {
	// List of scored repositories
	Repositories []*formatter.Record `json:"repositories"`
	// Timestamp of the response
	Timestamp time.Time `json:"timestamp" example:"2025-01-31T10:30:00Z"`
	// Total number of repositories requested
	TotalCount int `json:"total_count" example:"2"`
	// Number of successfully scored repositories
	SuccessCount int `json:"success_count" example:"2"`
	// Number of failed scorings
	ErrorCount int `json:"error_count" example:"0"`
}

// HealthResponse represents the health check response
// @Description Health status of the service
type HealthResponse struct {
	// Health status
	Status string `json:"status" example:"healthy"`
	// Current timestamp
	Timestamp time.Time `json:"timestamp" example:"2025-01-31T10:30:00Z"`
}

// ErrorResponse represents an error response
// @Description Error response from the API
type ErrorResponse struct {
	// Error message
	Error string `json:"error" example:"No repositories provided"`
	// Error code
	Code string `json:"code,omitempty" example:"NO_REPOSITORIES"`
	// Timestamp of the error
	Timestamp time.Time `json:"timestamp" example:"2025-01-31T10:30:00Z"`
}

// handleScore godoc
// @Summary Score GitHub repositories
// @Description Analyzes and scores a list of GitHub repositories based on various metrics
// @Tags analysis
// @Accept json
// @Produce json
// @Param request body ScoreRequest true "List of repositories to analyze"
// @Success 200 {object} ScoreResponse "Successfully scored repositories"
// @Failure 400 {object} ErrorResponse "Bad request"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /api/v1/score [post]
func (s *Server) handleScore(w http.ResponseWriter, r *http.Request) {
	var req ScoreRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body", "INVALID_REQUEST")
		return
	}

	if len(req.Repositories) == 0 {
		writeError(w, http.StatusBadRequest, "No repositories provided", "NO_REPOSITORIES")
		return
	}

	if len(req.Repositories) > 50 {
		writeError(w, http.StatusBadRequest, "Too many repositories (max 50)", "TOO_MANY_REPOSITORIES")
		return
	}

	response := &ScoreResponse{
		Repositories: make([]*formatter.Record, 0),
		Timestamp:    time.Now(),
		TotalCount:   len(req.Repositories),
	}

	for _, repoName := range req.Repositories {
		start := time.Now()
		metricsData, err := s.analyzer.Analyze(r.Context(), repoName)
		duration := time.Since(start)
		if err != nil {
			response.ErrorCount++
			s.metricsRecorder.RecordRepositoryAnalysis("error", duration)
			continue
		}

		record := formatter.MetricsToRecord(metricsData)
		response.Repositories = append(response.Repositories, record)
		response.SuccessCount++
		s.metricsRecorder.RecordRepositoryAnalysis("success", duration)
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to encode response", "ENCODING_ERROR")
	}
}

// handleHealth godoc
// @Summary Health check
// @Description Returns the health status of the service
// @Tags monitoring
// @Produce json
// @Success 200 {object} HealthResponse "Service is healthy"
// @Router /health [get]
func (s *Server) handleHealth(w http.ResponseWriter, _ *http.Request) {
	response := HealthResponse{
		Status:    "healthy",
		Timestamp: time.Now(),
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(response)
}

// handleRoot godoc
// @Summary Get API information
// @Description Returns information about the API and available endpoints
// @Tags info
// @Produce json
// @Success 200 {object} map[string]interface{} "API information"
// @Router / [get]
func (s *Server) handleRoot(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"service": "gh-inspector",
		"version": "1.0.0",
		"endpoints": map[string]string{
			"score":   "/api/v1/score",
			"health":  "/health",
			"metrics": "/metrics",
		},
	})
}

func writeError(w http.ResponseWriter, statusCode int, message string, code string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(ErrorResponse{
		Error:     message,
		Code:      code,
		Timestamp: time.Now(),
	})
}
