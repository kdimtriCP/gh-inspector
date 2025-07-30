package testutil

import (
	"os"
	"path/filepath"
	"testing"
)

func CreateTestConfig(t *testing.T) string {
	t.Helper()
	
	configContent := `github_token: "test-token"
output_format: "table"
cache:
  enabled: true
  ttl: 3600
  directory: ""

scoring:
  weights:
    stars: 0.30
    forks: 0.15
    recent_activity: 0.25
    open_issues: 0.10
    open_prs: 0.05
    has_license: 0.05
    has_cicd: 0.05
    has_contributing: 0.05
`
	
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.yaml")
	
	err := os.WriteFile(configPath, []byte(configContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test config: %v", err)
	}
	
	return configPath
}

func SetEnv(t *testing.T, key, value string) {
	t.Helper()
	
	oldValue := os.Getenv(key)
	err := os.Setenv(key, value)
	if err != nil {
		t.Fatalf("Failed to set env var %s: %v", key, err)
	}
	
	t.Cleanup(func() {
		os.Setenv(key, oldValue)
	})
}