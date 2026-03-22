package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	emrtypes "github.com/aws/aws-sdk-go-v2/service/emrserverless/types"
	"github.com/scttfrdmn/ood-emr-adapter/internal/emr"
	internalood "github.com/scttfrdmn/ood-emr-adapter/internal/ood"
	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status <application-id/job-run-id>",
	Short: "Get the status of an EMR Serverless job run",
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

		js := internalood.JobStatus{
			ID:     args[0],
			Status: emrStateToOod(detail.JobRun.State),
		}
		if detail.JobRun.StateDetails != nil {
			js.Message = *detail.JobRun.StateDetails
		}

		return json.NewEncoder(os.Stdout).Encode(js)
	},
}

func parseJobID(id string) (appID, jobRunID string, err error) {
	parts := strings.SplitN(id, "/", 2)
	if len(parts) == 2 {
		return parts[0], parts[1], nil
	}
	if applicationID != "" {
		return applicationID, id, nil
	}
	return "", "", fmt.Errorf("job ID must be in format application-id/job-run-id or --application-id must be set")
}

func emrStateToOod(s emrtypes.JobRunState) string {
	switch s {
	case emrtypes.JobRunStateSubmitted, emrtypes.JobRunStatePending, emrtypes.JobRunStateScheduled, emrtypes.JobRunStateQueued:
		return internalood.StatusQueued
	case emrtypes.JobRunStateRunning:
		return internalood.StatusRunning
	case emrtypes.JobRunStateSuccess:
		return internalood.StatusCompleted
	case emrtypes.JobRunStateFailed:
		return internalood.StatusFailed
	case emrtypes.JobRunStateCancelling, emrtypes.JobRunStateCancelled:
		return internalood.StatusCancelled
	default:
		return internalood.StatusUnknown
	}
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
