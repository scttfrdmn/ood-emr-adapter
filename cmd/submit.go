package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/scttfrdmn/ood-emr-adapter/internal/emr"
	"github.com/spf13/cobra"
)

// JobSpec is the EMR Serverless-specific job submission payload.
type JobSpec struct {
	ApplicationID         string            `json:"application_id,omitempty"`
	ExecutionRoleArn      string            `json:"execution_role_arn,omitempty"`
	EntryPoint            string            `json:"entry_point"`
	EntryPointArgs        []string          `json:"entry_point_args,omitempty"`
	SparkSubmitParameters string            `json:"spark_submit_parameters,omitempty"`
	Env                   map[string]string `json:"env,omitempty"`
	JobName               string            `json:"job_name,omitempty"`
}

var submitCmd = &cobra.Command{
	Use:   "submit",
	Short: "Submit an OOD job to Amazon EMR Serverless",
	Long:  "Reads a JSON job spec from stdin and submits it as an EMR Serverless job run.",
	RunE: func(cmd *cobra.Command, args []string) error {
		var spec JobSpec
		if err := json.NewDecoder(os.Stdin).Decode(&spec); err != nil {
			return fmt.Errorf("decode job spec: %w", err)
		}

		appID := applicationID
		if appID == "" {
			appID = spec.ApplicationID
		}
		if appID == "" {
			return fmt.Errorf("--application-id is required (or set application_id in job spec)")
		}
		if spec.EntryPoint == "" {
			return fmt.Errorf("job spec must include entry_point")
		}
		if spec.ExecutionRoleArn == "" {
			return fmt.Errorf("job spec must include execution_role_arn")
		}

		ctx := context.Background()
		client, err := emr.New(ctx, region)
		if err != nil {
			return err
		}

		jobRunID, err := client.StartJobRun(ctx, emr.JobRunSpec{
			ApplicationID:         appID,
			ExecutionRoleArn:      spec.ExecutionRoleArn,
			EntryPoint:            spec.EntryPoint,
			EntryPointArgs:        spec.EntryPointArgs,
			SparkSubmitParameters: spec.SparkSubmitParameters,
			Env:                   spec.Env,
			JobName:               spec.JobName,
		})
		if err != nil {
			return err
		}

		// Output format: "applicationID/jobRunID" so status/delete can parse both
		fmt.Printf("%s/%s\n", appID, jobRunID)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(submitCmd)
}
