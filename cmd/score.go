package cmd

import (
	"context"
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var repos []string

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

		analyzer := NewRepoAnalyzer(token)
		ctx := context.Background()

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
		_, _ = fmt.Fprintln(w, "REPOSITORY\tSTARS\tFORKS\tOPEN ISSUES\tOPEN PRS\tLAST COMMIT\tLANGUAGE\tCI/CD\tLICENSE")
		_, _ = fmt.Fprintln(w, "----------\t-----\t-----\t-----------\t--------\t-----------\t--------\t-----\t-------")

		for _, repo := range repos {
			metrics, err := analyzer.Analyze(ctx, repo)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error analyzing %s: %v\n", repo, err)
				continue
			}

			lastCommit := "N/A"
			if !metrics.LastCommitDate.IsZero() {
				daysAgo := int(time.Since(metrics.LastCommitDate).Hours() / 24)
				lastCommit = fmt.Sprintf("%d days ago", daysAgo)
			}

			cicd := "No"
			if metrics.HasCICD {
				cicd = "Yes"
			}

			license := "No"
			if metrics.HasLicense {
				license = "Yes"
			}

			lang := metrics.PrimaryLanguage
			if lang == "" {
				lang = "N/A"
			}

			fmt.Fprintf(w, "%s/%s\t%d\t%d\t%d\t%d\t%s\t%s\t%s\t%s\n",
				metrics.Owner, metrics.Name,
				metrics.Stars,
				metrics.Forks,
				metrics.OpenIssues,
				metrics.OpenPRs,
				lastCommit,
				lang,
				cicd,
				license,
			)
		}

		_ = w.Flush()
		return nil
	},
}

func init() {
	rootCmd.AddCommand(scoreCmd)
	scoreCmd.Flags().StringSliceVarP(&repos, "repos", "r", []string{}, "List of GitHub repositories")
}
