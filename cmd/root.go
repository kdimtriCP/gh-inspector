package cmd

import (
	"fmt"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "gh-inspector",
	Short: "Analyze GitHub repositories",
	Long:  `A CLI tool for scoring GitHub repositories.`,
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().String("config", "", "config file (default is ./configs/config.yaml)")
	if err := viper.BindPFlag("config", rootCmd.PersistentFlags().Lookup("config")); err != nil {
		fmt.Printf("Error binding config flag: %v\n", err)
	}
}

func initConfig() {
	viper.SetConfigType("yaml")
	viper.SetConfigName("config")
	viper.AddConfigPath("./configs")
	_ = godotenv.Load(".env")
	viper.AutomaticEnv()

	_ = viper.ReadInConfig()
}
