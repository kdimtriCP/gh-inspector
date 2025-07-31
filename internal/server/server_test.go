package server

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/golang/mock/gomock"
	"github.com/kdimtriCP/gh-inspector/internal/metrics"
	"github.com/kdimtriCP/gh-inspector/internal/mock/mock_github"
)

func TestHealthEndpoint(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAnalyzer := mock_github.NewMockAnalyzer(ctrl)
	srv := New(mockAnalyzer, nil)

	req := httptest.NewRequest("GET", "/health", nil)
	rr := httptest.NewRecorder()

	srv.router.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)

	var response HealthResponse
	err := json.NewDecoder(rr.Body).Decode(&response)
	require.NoError(t, err)
	require.Equal(t, "healthy", response.Status)
	require.False(t, response.Timestamp.IsZero())
}

func TestRootEndpoint(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAnalyzer := mock_github.NewMockAnalyzer(ctrl)
	srv := New(mockAnalyzer, nil)

	req := httptest.NewRequest("GET", "/", nil)
	rr := httptest.NewRecorder()

	srv.router.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)

	var response map[string]interface{}
	err := json.NewDecoder(rr.Body).Decode(&response)
	require.NoError(t, err)
	require.Equal(t, "gh-inspector", response["service"])
	require.Equal(t, "1.0.0", response["version"])
	require.NotNil(t, response["endpoints"])
}

func TestScoreEndpoint(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAnalyzer := mock_github.NewMockAnalyzer(ctrl)
	srv := New(mockAnalyzer, nil)

	t.Run("successful scoring", func(t *testing.T) {
		reqBody := ScoreRequest{
			Repositories: []string{"test/repo1", "test/repo2"},
		}

		mockAnalyzer.EXPECT().
			Analyze(gomock.Any(), "test/repo1").
			Return(&metrics.Repository{
				Owner: "test",
				Name:  "repo1",
				Stars: 100,
				Score: 85.5,
			}, nil)

		mockAnalyzer.EXPECT().
			Analyze(gomock.Any(), "test/repo2").
			Return(&metrics.Repository{
				Owner: "test",
				Name:  "repo2",
				Stars: 200,
				Score: 92.0,
			}, nil)

		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("POST", "/api/v1/score", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		srv.router.ServeHTTP(rr, req)

		require.Equal(t, http.StatusOK, rr.Code)

		var response ScoreResponse
		err := json.NewDecoder(rr.Body).Decode(&response)
		require.NoError(t, err)
		require.Len(t, response.Repositories, 2)
		require.Equal(t, 2, response.TotalCount)
		require.Equal(t, 2, response.SuccessCount)
		require.Equal(t, 0, response.ErrorCount)
	})

	t.Run("empty repositories", func(t *testing.T) {
		reqBody := ScoreRequest{
			Repositories: []string{},
		}

		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("POST", "/api/v1/score", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		srv.router.ServeHTTP(rr, req)

		require.Equal(t, http.StatusBadRequest, rr.Code)

		var response ErrorResponse
		err := json.NewDecoder(rr.Body).Decode(&response)
		require.NoError(t, err)
		require.Equal(t, "No repositories provided", response.Error)
		require.Equal(t, "NO_REPOSITORIES", response.Code)
	})

	t.Run("too many repositories", func(t *testing.T) {
		repos := make([]string, 51)
		for i := 0; i < 51; i++ {
			repos[i] = "test/repo"
		}

		reqBody := ScoreRequest{
			Repositories: repos,
		}

		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("POST", "/api/v1/score", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		srv.router.ServeHTTP(rr, req)

		require.Equal(t, http.StatusBadRequest, rr.Code)

		var response ErrorResponse
		err := json.NewDecoder(rr.Body).Decode(&response)
		require.NoError(t, err)
		require.Equal(t, "Too many repositories (max 50)", response.Error)
		require.Equal(t, "TOO_MANY_REPOSITORIES", response.Code)
	})

	t.Run("invalid request body", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/api/v1/score", bytes.NewReader([]byte("invalid json")))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		srv.router.ServeHTTP(rr, req)

		require.Equal(t, http.StatusBadRequest, rr.Code)

		var response ErrorResponse
		err := json.NewDecoder(rr.Body).Decode(&response)
		require.NoError(t, err)
		require.Equal(t, "Invalid request body", response.Error)
		require.Equal(t, "INVALID_REQUEST", response.Code)
	})
}

func TestCORSMiddleware(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAnalyzer := mock_github.NewMockAnalyzer(ctrl)
	srv := New(mockAnalyzer, nil)

	req := httptest.NewRequest("OPTIONS", "/api/v1/score", nil)
	rr := httptest.NewRecorder()

	srv.router.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
	require.Equal(t, "*", rr.Header().Get("Access-Control-Allow-Origin"))
	require.Contains(t, rr.Header().Get("Access-Control-Allow-Methods"), "POST")
	require.Contains(t, rr.Header().Get("Access-Control-Allow-Headers"), "Content-Type")
}

func TestServerShutdown(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAnalyzer := mock_github.NewMockAnalyzer(ctrl)
	srv := New(mockAnalyzer, &Config{
		Port:         0, // Use random port
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	})

	// Start server in goroutine
	serverStarted := make(chan struct{})
	go func() {
		close(serverStarted)
		srv.Start()
	}()

	// Wait for server to start
	<-serverStarted
	time.Sleep(200 * time.Millisecond)

	// Shutdown server
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := srv.Shutdown(ctx)
	require.NoError(t, err)
}
