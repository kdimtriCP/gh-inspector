package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/kdimtriCP/gh-inspector/internal/cache"
	"github.com/kdimtriCP/gh-inspector/internal/github"
	"github.com/kdimtriCP/gh-inspector/internal/scoring"
	"github.com/kdimtriCP/gh-inspector/internal/server"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the gh-inspector web service",
	Long: `Start the gh-inspector as a web service with REST API endpoints.
	
The service provides:
- POST /api/v1/score - Score GitHub repositories
- GET /health - Health check endpoint
- GET /metrics - Prometheus metrics endpoint
- GET /swagger/* - Swagger UI documentation`,
	RunE: runServe,
}

var (
	port         int
	readTimeout  int
	writeTimeout int
)

func init() {
	rootCmd.AddCommand(serveCmd)

	serveCmd.Flags().IntVarP(&port, "port", "p", 8080, "Server port")
	serveCmd.Flags().IntVar(&readTimeout, "read-timeout", 15, "Read timeout in seconds")
	serveCmd.Flags().IntVar(&writeTimeout, "write-timeout", 15, "Write timeout in seconds")
}

func runServe(_ *cobra.Command, _ []string) error {
	token := viper.GetString("github_token")
	if token == "" || token == "ghp_yourtokenhere" {
		return fmt.Errorf("GitHub token not configured. Please set github_token in config file or GITHUB_TOKEN environment variable")
	}

	cacheEnabled := viper.GetBool("cache.enabled")
	cacheDir := viper.GetString("cache.directory")
	cacheTTL := time.Duration(viper.GetInt("cache.ttl")) * time.Second

	var cacheInstance cache.Cache
	if cacheEnabled {
		var err error
		cacheInstance, err = cache.New(cacheDir)
		if err != nil {
			fmt.Printf("Warning: failed to initialize cache: %v\n", err)
		}
	}

	scoringConfig := &scoring.Config{
		Weights: scoring.Weights{
			Stars:            viper.GetFloat64("scoring.weights.stars"),
			Forks:            viper.GetFloat64("scoring.weights.forks"),
			RecentActivity:   viper.GetFloat64("scoring.weights.recent_activity"),
			OpenIssues:       viper.GetFloat64("scoring.weights.open_issues"),
			OpenPRs:          viper.GetFloat64("scoring.weights.open_prs"),
			HasLicense:       viper.GetFloat64("scoring.weights.has_license"),
			HasCICD:          viper.GetFloat64("scoring.weights.has_cicd"),
			HasContributing:  viper.GetFloat64("scoring.weights.has_contributing"),
			ReleaseFrequency: viper.GetFloat64("scoring.weights.release_frequency"),
			HasReadme:        viper.GetFloat64("scoring.weights.has_readme"),
			HasCodeOfConduct: viper.GetFloat64("scoring.weights.has_code_of_conduct"),
			HasSecurity:      viper.GetFloat64("scoring.weights.has_security"),
			Watchers:         viper.GetFloat64("scoring.weights.watchers"),
		},
	}

	analyzer := github.NewRepoAnalyzer(token, scoringConfig)
	if cacheInstance != nil {
		analyzer.SetCache(cacheInstance)
		analyzer.SetCacheTTL(cacheTTL)
	}

	serverConfig := &server.Config{
		Port:         port,
		ReadTimeout:  time.Duration(readTimeout) * time.Second,
		WriteTimeout: time.Duration(writeTimeout) * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	srv := server.New(analyzer, serverConfig)

	// Inject metrics recorder into analyzer
	analyzer.SetMetricsRecorder(srv.MetricsRecorder())

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		if err := srv.Start(); err != nil {
			fmt.Printf("Server error: %v\n", err)
			os.Exit(1)
		}
	}()

	<-sigChan
	fmt.Println("\nShutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		return fmt.Errorf("server shutdown failed: %w", err)
	}

	fmt.Println("Server stopped gracefully")
	return nil
}
