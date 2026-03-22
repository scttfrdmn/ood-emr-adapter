package cmd

import (
	"context"
	"encoding/json"
	"os"

	"github.com/scttfrdmn/ood-emr-adapter/internal/emr"
	"github.com/spf13/cobra"
)

var infoCmd = &cobra.Command{
	Use:   "info <application-id/job-run-id>",
	Short: "Print full EMR Serverless job run details as JSON",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		appID, jobRunID, err := parseJobID(args[0])
		if err != nil {
			return err
		}

		ctx := context.Background()
		client, err := emr.New(ctx, region)
		if err != nil {
			return err
		}
		detail, err := client.GetJobRun(ctx, appID, jobRunID)
		if err != nil {
			return err
		}
		return json.NewEncoder(os.Stdout).Encode(detail)
	},
}

func init() {
	rootCmd.AddCommand(infoCmd)
}
