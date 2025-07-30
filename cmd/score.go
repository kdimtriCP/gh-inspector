package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var repos []string

var scoreCmd = &cobra.Command{
	Use:   "score",
	Short: "Score GitHub repositories",
	Run: func(cmd *cobra.Command, args []string) {
		for _, repo := range repos {
			fmt.Printf("Analyzing repo: %s ", repo)
			// Здесь будет вызов логики анализа
		}
	},
}

func init() {
	rootCmd.AddCommand(scoreCmd)
	scoreCmd.Flags().StringSliceVarP(&repos, "repos", "r", []string{}, "List of GitHub repositories")
}
