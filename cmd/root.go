package cmd

import (
	"github.com/spf13/cobra"
)

var (
	region        string
	applicationID string
)

var rootCmd = &cobra.Command{
	Use:   "ood-emr-adapter",
	Short: "OOD compute adapter for Amazon EMR Serverless",
	Long:  "Translates Open OnDemand job submissions to Amazon EMR Serverless API calls.",
}

// Execute runs the root command.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringVar(&region, "region", "us-east-1", "AWS region")
	rootCmd.PersistentFlags().StringVar(&applicationID, "application-id", "", "EMR Serverless application ID")
}
