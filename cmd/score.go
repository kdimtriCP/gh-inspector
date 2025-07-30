package cmd

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/kdimtriCP/gh-inspector/internal/cache"
	"github.com/kdimtriCP/gh-inspector/internal/formatter"
	"github.com/kdimtriCP/gh-inspector/internal/github"
	"github.com/kdimtriCP/gh-inspector/internal/metrics"
	"github.com/kdimtriCP/gh-inspector/internal/scoring"
)

var (
	repos        []string
	outputFormat string
	noCache      bool
)

var scoreCmd = &cobra.Command{
	Use:   "score",
	Short: "Score GitHub repositories",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(repos) == 0 {
			return fmt.Errorf("no repositories specified")
		}

		token := viper.GetString("github_token")
		if token == "" {
			return fmt.Errorf("GitHub token not configured")
		}

		scoringConfig := &scoring.Config{}
		_ = viper.UnmarshalKey("scoring", scoringConfig)

		analyzer := github.NewRepoAnalyzer(token, scoringConfig)
		ctx := context.Background()

		if viper.GetBool("cache.enabled") && !noCache {
			cacheDir := viper.GetString("cache.directory")
			c, err := cache.New(cacheDir)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Warning: failed to initialize cache: %v\n", err)
			} else {
				analyzer.SetCache(c)
				defer c.Close()

				cacheTTL := viper.GetInt("cache.ttl")
				if cacheTTL > 0 {
					ttlDuration := time.Duration(cacheTTL) * time.Second
					analyzer.SetCacheTTL(ttlDuration)
				}
			}
		}

		var allMetrics []*metrics.Repository

		for _, repo := range repos {
			metrics, err := analyzer.Analyze(ctx, repo)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error analyzing %s: %v\n", repo, err)
				continue
			}
			allMetrics = append(allMetrics, metrics)
		}

		if len(allMetrics) == 0 {
			return fmt.Errorf("no repositories could be analyzed")
		}

		format := outputFormat
		if format == "" {
			format = viper.GetString("output_format")
		}
		if format == "" {
			format = formatter.FormatTable
		}

		formatter, err := formatter.New(format)
		if err != nil {
			return err
		}

		return formatter.Format(os.Stdout, allMetrics)
	},
}

func init() {
	rootCmd.AddCommand(scoreCmd)
	scoreCmd.Flags().StringSliceVarP(&repos, "repos", "r", []string{}, "List of GitHub repositories")
	scoreCmd.Flags().StringVarP(&outputFormat, "output", "o", "", "Output format (table, json, json-compact, csv)")
	scoreCmd.Flags().BoolVar(&noCache, "no-cache", false, "Disable caching")
}
