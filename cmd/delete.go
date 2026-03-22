package cmd

import (
	"context"
	"fmt"

	"github.com/scttfrdmn/ood-emr-adapter/internal/emr"
	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete <application-id/job-run-id>",
	Short: "Cancel an EMR Serverless job run",
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
		if err := client.CancelJobRun(ctx, appID, jobRunID); err != nil {
			return err
		}
		fmt.Printf("Job run %s cancelled\n", args[0])
		return nil
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}
